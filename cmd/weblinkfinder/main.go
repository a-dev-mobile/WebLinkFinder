package main

import (
	"WebLinkFinder/pkg/WebLinkFinder"
	"fmt"
	"regexp"
)

func main() {

	queue := WebLinkFinder.NewQueue(
		"https://9animetv.to/tv",
		100,
		false,
		[]*regexp.Regexp{regexp.MustCompile(`.*top-.*`)},
		[]*regexp.Regexp{
			regexp.MustCompile(`.*watch/.*`),
		})

	queueLinks := queue.GetQueue()

	for _, link := range queueLinks {
		fmt.Println(link)
	}
	fmt.Println("")
	fmt.Println("number of links in the queue ", len(queueLinks))
	fmt.Println("")

	links := WebLinkFinder.NewLinks(
		queueLinks,
		100,
		false,
		[]*regexp.Regexp{
			regexp.MustCompile(`.*top-.*`),
		},
		[]*regexp.Regexp{
			regexp.MustCompile(`.*ura.*`),
		})
	findLinks := links.GetLinks()

	for _, link := range findLinks {
		fmt.Println(link)
	}
	fmt.Println("")
	fmt.Println("number of links ", len(findLinks))
	fmt.Println("")

}
