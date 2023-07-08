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
		IsSaveLocalFile:    true,
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

	findQueue := queue.GetQueue()

	// Вывод ссылок из очереди и их количество
	for _, link := range findQueue {
		fmt.Println(link)
	}
	fmt.Println("")
	fmt.Printf("Number of links in the queue: %d\n", len(findQueue))
	fmt.Println("")

	links := weblinkfinder.NewLinks(findQueue, config)

	findLinks := links.GetLinks()

	// Вывод найденных ссылок и их количество
	for _, link := range findLinks {
		fmt.Println(link)
	}
	fmt.Println("")
	fmt.Printf("Number of found links: %d\n", len(findLinks))
	fmt.Println("")

}
