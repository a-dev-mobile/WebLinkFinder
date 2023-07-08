package main

import (
	"fmt"

	weblinkfinder "github.com/a-dev-mobile/weblinkfinder/pkg"
)

func main() {


	// Стартовый URL, максимальное количество горутин, режим отладки
	config := &weblinkfinder.WebCrawlerConfig{
		MaxGoroutine:       3,
		MaxRequests:        10,
		IsDebug:            true,
		QueuesSkipStrRegex: []string{},
		QueuesAllowStrRegex: []string{
			`.*packages\?page.*`,
			`.*/packages$`,
		},
		LinksSkipStrRegex: []string{},
		LinksAllowStrRegex: []string{
			`.*/packages\/[\w]+$`},
	}

	startURL := "https://pub.dev/packages"

	queue := weblinkfinder.NewQueue(startURL, config)

	queueLinks := queue.GetQueue()

	// Вывод ссылок из очереди и их количество
	for _, link := range queueLinks {
		fmt.Println(link)
	}
	fmt.Println("")
	fmt.Printf("Number of links in the queue: %d\n", len(queueLinks))
	fmt.Println("")


	links := weblinkfinder.NewLinks(queueLinks, config)

	findLinks := links.GetLinks()

	// Вывод найденных ссылок и их количество
	for _, link := range findLinks {
		fmt.Println(link)
	}
	fmt.Println("")
	fmt.Printf("Number of found links: %d\n", len(findLinks))
	fmt.Println("")
}
