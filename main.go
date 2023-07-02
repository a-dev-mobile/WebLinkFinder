package main

import (
	"WebLinkFinder/utils/arrutils"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

type privateConstants struct {
	startUrl string
}

var skipQueueRegexes = []*regexp.Regexp{
	regexp.MustCompile(`.*blog.*`),
	// regexp.MustCompile(`.*people.*`),

}

// добавлять в очередь ссылки:
var allowQueueRegexes = []*regexp.Regexp{
	// regexp.MustCompile(`.*/diary.*`),
	regexp.MustCompile(`.*calorie.*`),
}

var constants = privateConstants{
	startUrl: "https://health-diet.ru/",
}

func main() {

	// website := "https://health-diet.ru/"
	visited := make(map[string]bool)
	toVisit := []string{constants.startUrl}
	count := 0
	for len(toVisit) > 0 {
		visitURL := toVisit[0]
		toVisit = toVisit[1:]
		
		if visited[visitURL] {
			continue
		}
		count++

		urls, err := fetchLinks(visitURL, skipQueueRegexes, allowQueueRegexes)

		if err != nil {
			fmt.Println("ERROR: fetchLinks \"" + visitURL + "\"")

		}

		toVisit = append(toVisit, urls...)
		// toVisit = arrutils.UniqueStr(toVisit)

		fmt.Println(count, " Visited " + visitURL)

		visited[visitURL] = true
	}
	

}

func fetchLinks(website string, skipRegexes []*regexp.Regexp, allowRegexes []*regexp.Regexp) ([]string, error) {
	resp, err := http.Get(website)
	if err != nil {
		fmt.Println("ERROR: Failed to crawl \"" + website + "\"")
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("ERROR: Could not read \"" + website + "\"")
		return nil, err
	}

	u, _ := url.Parse(website)


	// Формируем шаблон регулярного выражения
	reBaseUrl, err := regexp.Compile(fmt.Sprintf(`^.*%s.*`, u.Host))


	if err != nil {
		fmt.Printf("Invalid pattern: %v\n", err)
		return nil, err
	}

	doc, err := html.Parse(strings.NewReader(string(body)))
	if err != nil {
		fmt.Println("ERROR: Could not parse \"" + website + "\"")
		return nil, err
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

	return arrutils.UniqueStr(results), err

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
