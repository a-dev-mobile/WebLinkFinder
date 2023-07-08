// / package weblinkfinder содержит логику для поиска веб-ссылок на сайтах.
package weblinkfinder

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/a-dev-mobile/weblinkfinder/utils/arrutils"
	"github.com/a-dev-mobile/weblinkfinder/utils/dicutils"
	"github.com/a-dev-mobile/weblinkfinder/utils/regexutils"
	"golang.org/x/net/html"
)

/// Queue represents a queue data structure for searching web links.

type Queue struct {
	startUrl string
	config   *WebCrawlerConfig
}

// / Links represents a data structure for storing web links.
type Links struct {
	queues []string
	config *WebCrawlerConfig
}

type WebCrawlerConfig struct {
	MaxGoroutine        int
	MaxRequests         int
	IsDebug             bool
	IsSaveLocalFile     bool
	QueuesSkipStrRegex  []string
	QueuesAllowStrRegex []string
	LinksSkipStrRegex   []string
	LinksAllowStrRegex  []string
}

// / NewQueue creates a new Queue instance.
func NewQueue(startUrl string, config *WebCrawlerConfig) *Queue {
	return &Queue{
		startUrl: startUrl,
		config:   config,
	}
}

// / NewLinks creates a new Links instance.
func NewLinks(queues []string, config *WebCrawlerConfig) *Links {
	return &Links{
		queues: queues,
		config: config,
	}
}

// / GetQueue searches for web links starting at startUrl.
// / Returns a list of web links found on sites.
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

		if len(visitArr) == 0 {
			break
		}
		// do not wait for all links to be processed
		lenghtVisited := len(dicutils.GetKeysWithTrue(visitMap))
		if lenghtVisited >= f.config.MaxRequests {
			break
		}

		fmt.Printf("\nVerified links %d\n\n", lenghtVisited)
		var lenght int

		if f.config.MaxGoroutine < len(visitArr) {
			lenght = f.config.MaxGoroutine
		} else {
			lenght = len(visitArr)

		}

		if f.config.IsDebug {
			fmt.Printf("Debug: Start of loop iteration %d, Goroutine start %d\n", count, lenght)
		}
		for i := 0; i < lenght; i++ {
			wg.Add(1)

			mu.Lock()
			visitMap[visitArr[i]] = true
			mu.Unlock()

			if f.config.IsDebug {
				fmt.Printf("Debug: Visiting %s\n", visitArr[i])
			}

			go func(id int, cycle int) {
				defer wg.Done()
				if f.config.IsDebug {
					fmt.Printf("Debug: Inside goroutine for cycle %d, Visiting %s\n", cycle, visitArr[id])
				}
				queues := fetchLinks(visitArr[id], f.config.QueuesSkipStrRegex, f.config.QueuesAllowStrRegex)
				if f.config.IsDebug && len(queues) > 0 {
					fmt.Printf("Debug: Found %d links in %s\n", len(queues), visitArr[id])
				}
				for _, s := range queues {
					queueCh <- s

				}

			}(i, count)

		}

		wg.Wait()
	}

	close(queueCh)

	var findQueue []string

	for k := range visitMap {
		findQueue = append(findQueue, k)
	}
	if f.config.IsSaveLocalFile {

		saveToFile(findQueue, "queue.txt")

	}

	return findQueue
}

// GetLinks searches for web links for all URLs in queues.
// Returns a list of web links found on sites.
func (f *Links) GetLinks() []string {
	linksCh := make(chan string)

	findLinks := []string{}
	links := f.queues
	var mu sync.Mutex

	var wg sync.WaitGroup

	sum := len(f.queues)

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

		fmt.Printf("\n%d checked - %d found\n\n", sum-len(links), len(findLinks))

		if len(links) == 0 {
			break
		}

		var lenght int

		if f.config.MaxGoroutine < len(links) {
			lenght = f.config.MaxGoroutine
		} else {
			lenght = len(links)

		}

		if f.config.IsDebug {
			fmt.Printf("Debug: Start of loop iteration %d, Goroutine start %d\n", count, lenght)
		}

		for i := 0; i < lenght; i++ {
			wg.Add(1)

			url := links[i]
			if f.config.IsDebug {
				fmt.Printf("Debug: Visiting %s\n", url)
			}
			go func(url string, cycle int) {
				defer wg.Done()
				if f.config.IsDebug {
					fmt.Printf("Debug: Inside goroutine for cycle %d, Visiting %s\n", cycle, url)
				}

				urls := fetchLinks(url, f.config.LinksSkipStrRegex, f.config.LinksAllowStrRegex)

				mu.Lock()
				links = arrutils.DeleteElement(links, url)
				mu.Unlock()

				if f.config.IsDebug && len(urls) > 0 {
					fmt.Printf("Debug: Found %d links in %s\n", len(urls), url)
				}
				for _, s := range urls {
					linksCh <- s

				}

			}(url, count)

		}

		wg.Wait()
	}

	close(linksCh)

	if f.config.IsSaveLocalFile {
		saveToFile(findLinks, "links.txt")

	}

	return findLinks
}

// / fetchLinks retrieves web links from a given site, taking into account the rules
// / skipRegexes и allowRegexes.
func fetchLinks(website string, skipStrRegex []string, allowStrRegex []string) []string {
	skipRegex, err := regexutils.CompileRegexes(skipStrRegex)
	if err != nil {
		log.Fatal(err)
	}
	allowRegex, err := regexutils.CompileRegexes(allowStrRegex)
	if err != nil {
		log.Fatal(err)
	}

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
	/// Forming a Regular Expression Pattern
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

					for _, regex := range allowRegex {
						if regex.MatchString(absolute) {
							allow = true
							break
						} else {
							allow = false

						}
					}

					for _, regex := range skipRegex {
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

// fixUrl converts a relative URL to an absolute URL based on the base URL.
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
func saveToFile(links []string, nameFile string) {
	file, err := os.OpenFile(nameFile, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	datawriter := bufio.NewWriter(file)

	for _, link := range links {
		_, _ = datawriter.WriteString(link + "\n")
	}
	datawriter.Flush()
	file.Close()
}
