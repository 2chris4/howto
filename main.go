// file: list_posts.go
package main

import (
	"fmt"
	"os"
	"strings"
	"sync"
)

func main() {

	var query string
	if len(os.Args) > 1 {
		query = "How to " + strings.Join(os.Args[1:], " ")
		query = strings.Replace(query, "\n", "", -1)
		search(query)
	} else {
		fmt.Println("Usage: howto <Your-Query> : mode <google/stackoverflow>")
	}
}
func search(query string) {
	mode := strings.Split(query, ":")
	var items []post
	if len(mode) == 1 {
		fmt.Println("Searching query ", query, " using google")
		items = chromedpHandler(query)
	} else {
		if mode[1] == "google" {
			fmt.Println("Searching query ", mode[0], " using google")
			items = chromedpHandler(mode[0])

		} else {
			fmt.Println("Searching query ", mode[0], " using Stack Overflow")
			items = searchPost(mode[0])

		}
	}

	var wg sync.WaitGroup
	wg.Add(len(items))
	messages := make(chan string, len(items))
	title, description, _ := getTitleDescription(items)
	acceptedAnswers := make(map[int]solution)
	allAnswers := make(map[int][]solution)
	gettingPost := func(link string, index int, item post) {
		defer wg.Done()
		acceptedTemp, temp := getPost(item)
		allAnswers[index] = temp
		acceptedAnswers[index] = acceptedTemp
		messages <- "done"
	}

	for index, item := range items {
		go gettingPost(item.link, index, item)
	}

	wg.Wait()
	createUI(title, description, allAnswers, acceptedAnswers)
}
