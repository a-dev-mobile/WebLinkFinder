package main

import (
	"fmt"

	weblinkfinder "weblinkfinder/pkg"
)

func main() {
	// URL-ы для пропуска и допуска в зависимости от совпадения с регулярными выражениями

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

	// Стартовый URL, максимальное количество горутин, режим отладки
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

	// Смена режима отладки и обновление паттернов для сканирования

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
