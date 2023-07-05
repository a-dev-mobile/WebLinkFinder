package WebLinkFinder

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"

	"WebLinkFinder/utils/arrutils"
	"WebLinkFinder/utils/dicutils"

	"golang.org/x/net/html"
)

type Queue struct {
	startUrl     string
	maxGoroutine int
	isDebug      bool
	skipRegexes  []*regexp.Regexp
	allowRegexes []*regexp.Regexp
}
type Links struct {
	queues       []string
	maxGoroutine int
	isDebug      bool
	skipRegexes  []*regexp.Regexp
	allowRegexes []*regexp.Regexp
}

func NewQueue(startUrl string, maxGoroutine int, isDebug bool, skipRegexes []*regexp.Regexp, allowRegexes []*regexp.Regexp) *Queue {
	return &Queue{
		startUrl:     startUrl,
		maxGoroutine: maxGoroutine,
		isDebug:      isDebug,
		skipRegexes:  skipRegexes,
		allowRegexes: allowRegexes,
	}
}
func NewLinks(queues []string, maxGoroutine int, isDebug bool, skipRegexes []*regexp.Regexp, allowRegexes []*regexp.Regexp) *Links {
	return &Links{
		queues:       queues,
		maxGoroutine: maxGoroutine,
		isDebug:      isDebug,
		skipRegexes:  skipRegexes,
		allowRegexes: allowRegexes,
	}
}
func (f *Queue) GetQueue() []string {
	queueCh := make(chan string)

	visitMap := make(map[string]bool)
	var mu sync.Mutex
	visitMap[f.startUrl] = false

	var wg sync.WaitGroup

	go func() {

		for s := range queueCh {
			mu.Lock()
			dicutils.AddToMapIfNotExist(visitMap, s, false)
			mu.Unlock()
		}
	}()
	count := 0

	for {
		count++
		mu.Lock()
		visitArr := dicutils.GetKeysWithFalse(visitMap)
		mu.Unlock()
		
		fmt.Printf("It remains to check %d\n", len(visitArr))
		if len(visitArr) == 0 {
			break
		}

		var lenght int

		if f.maxGoroutine < len(visitArr) {
			lenght = f.maxGoroutine
		} else {
			lenght = len(visitArr)

		}

		if f.isDebug {
			fmt.Printf("Debug: Start of loop iteration %d, Goroutine start %d\n", count, lenght)
		}
		for i := 0; i < lenght; i++ {
			wg.Add(1)

			mu.Lock()
			visitMap[visitArr[i]] = true
			mu.Unlock()

			go func(id int, cycle int) {
				defer wg.Done()
				if f.isDebug {
					fmt.Printf("Debug: Inside goroutine for cycle %d, Visiting %s\n", cycle, visitArr[id])
				}
				queues := fetchLinks(visitArr[id], f.skipRegexes, f.allowRegexes)
				if f.isDebug && len(queues) > 0 {
					fmt.Printf("Debug: Found %d links in %s\n", len(queues), visitArr[id])
				}
				for _, s := range queues {
					queueCh <- s

				}

			}(i, count)

		}

		wg.Wait() // wait for all goroutines to finish
	}

	close(queueCh)

	// close(linksCh)

	var keys []string

	for k := range visitMap {
		keys = append(keys, k)
	}

	return keys
}

func (f *Links) GetLinks() []string {
	linksCh := make(chan string)

	findLinks := []string{}
	links := f.queues
	var mu sync.Mutex
	// visitMap[f.startUrl] = false

	var wg sync.WaitGroup

	go func() {

		for s := range linksCh {
			mu.Lock()

			findLinks = arrutils.AddIfUnique(findLinks, s)
			mu.Unlock()
		}
	}()
	count := 0

	for {
		count++
		// mu.Lock()
		// visitArr := dicutils.GetKeysWithFalse(visitMap)
		// mu.Unlock()

		fmt.Printf("It remains to check %d\n", len(links))
		if len(links) == 0 {
			break
		}

		var lenght int

		if f.maxGoroutine < len(links) {
			lenght = f.maxGoroutine
		} else {
			lenght = len(links)

		}

		if f.isDebug {
			fmt.Printf("Debug: Start of loop iteration %d, Goroutine start %d\n", count, lenght)
		}

		for i := 0; i < lenght; i++ {
			wg.Add(1)

			// mu.Lock()
			// if i >= lenght {

			// }
			// mu.Unlock()
			url := links[i]

			go func(url string, cycle int) {
				defer wg.Done()
				if f.isDebug {
					fmt.Printf("Debug: Inside goroutine for cycle %d, Visiting %s\n", cycle, url)
				}

				queues := fetchLinks(url, f.skipRegexes, f.allowRegexes)

				mu.Lock()
				links = arrutils.DeleteElement(links, url)
				mu.Unlock()

				if f.isDebug && len(queues) > 0 {
					fmt.Printf("Debug: Found %d links in %s\n", len(queues), url)
				}
				for _, s := range queues {
					linksCh <- s

				}

			}(url, count)

		}

		wg.Wait() // wait for all goroutines to finish
	}

	close(linksCh)

	return findLinks
}

func fetchLinks(website string, skipRegexes []*regexp.Regexp, allowRegexes []*regexp.Regexp) []string {
	resp, err := http.Get(website)
	if err != nil {
		fmt.Println("ERROR: Failed to crawl \"" + website + "\"")
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("ERROR: Could not read \"" + website + "\"")
		return nil
	}

	u, err := url.Parse(website)
	if err != nil {
		fmt.Println("ERROR: Parse \"" + website + "\"")
	}
	// Формируем шаблон регулярного выражения
	reBaseUrl, err := regexp.Compile(fmt.Sprintf(`^.*%s.*`, u.Host))

	if err != nil {
		fmt.Printf("Invalid pattern: %v\n", err)
		return nil
	}

	doc, err := html.Parse(strings.NewReader(string(body)))
	if err != nil {
		fmt.Println("ERROR: Could not parse \"" + website + "\"")
		return nil

	}

	var f func(*html.Node)

	var results []string
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {

					allow := false

					absolute := fixUrl(a.Val, website)
					if website == absolute {
						continue
					}

					if reBaseUrl.MatchString(absolute) {
						allow = true
					} else {
						continue
					}

					for _, regex := range allowRegexes {
						if regex.MatchString(absolute) {
							allow = true
							break
						} else {
							allow = false

						}
					}

					for _, regex := range skipRegexes {
						if regex.MatchString(absolute) {
							allow = false
							break
						}
					}
					if allow {
						results = append(results, absolute)
					}
				}

			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return arrutils.UniqueStr(results)

}

func fixUrl(href, base string) (absolute string) {
	uri, err := url.Parse(href)
	if err != nil {
		return
	}
	baseUrl, err := url.Parse(base)
	if err != nil {
		return
	}
	absolute = baseUrl.ResolveReference(uri).String()
	return
}
