package main

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"golang.org/x/net/html"
)
var baseUrl string = "https://scrape-me.dreamsofcode.io"
type VisitedSet struct {
    mu sync.Mutex
    m  map[string]struct{}
}
var visited = VisitedSet{
    m: make(map[string]struct{}),
}
var dead []string
func scrape(url string)  {
	visited.mu.Lock()
	if _, ok := visited.m[url]; ok {
		visited.mu.Unlock()
		fmt.Println("Already visited!")
		return
	}
	visited.m[url] = struct{}{}
	visited.mu.Unlock()
	if !strings.HasPrefix(url, baseUrl) {
		fmt.Println(url)
		fmt.Println("Doesn't have prefix")
		return
	}
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("error fetching url.")
		return
	}
	fmt.Println(resp.StatusCode) 
	if resp.StatusCode != 200 {
		dead = append(dead, url)
		fmt.Println("Dead link reached.")
		return
	}
	if err != nil {
		fmt.Println(err)
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		fmt.Println("Error parsing body..")
	}
	defer resp.Body.Close()
	for n := range doc.Descendants() {
		if n.Type == html.ElementNode && n.Data == "a" {
			for index := range n.Attr {
				if n.Attr[index].Key == "href" {
					link := n.Attr[index].Val
					if strings.HasPrefix(link, "/") {
						scrape(baseUrl+link)
					} else {
						scrape(link)
					}
				}
			}
		}
	}
}

func main() {
	scrape(baseUrl)
	fmt.Println(dead)

}