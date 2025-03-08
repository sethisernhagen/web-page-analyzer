package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type Stats struct {
	URL          string
	WordCount    int
	ImageCount   int
	LinkCount    int
	StatusCode   int
	ResponseTime time.Duration
	Error        error
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a URL")
		os.Exit(1)
	}

	urls := os.Args[1:]
	results := analyzeURLs(urls)

	fmt.Println("Results:")
	fmt.Println("--------------------")
	for _, result := range results {
		fmt.Printf("URL: %s\n", result.URL)
		if result.Error != nil {
			fmt.Printf("Error: %s\n", result.Error)
			continue
		}
		fmt.Printf("Word count: %d\n", result.WordCount)
		fmt.Printf("Image count: %d\n", result.ImageCount)
		fmt.Printf("Link count: %d\n", result.LinkCount)
		fmt.Printf("Status code: %d\n", result.StatusCode)
		fmt.Printf("Response time: %s\n", result.ResponseTime)
		fmt.Println("--------------------")
	}
}

func analyzeURLs(urls []string) []Stats {
	wg := sync.WaitGroup{}
	results := make([]Stats, len(urls))

	for i, url := range urls {
		wg.Add(1)
		go func(i int, url string) {
			defer wg.Done()
			results[i] = analyzeURL(url)
		}(i, url)
	}
	wg.Wait()
	return results
}

func analyzeURL(url string) Stats {
	url = getUrl(url)
	stats := Stats{URL: url}
	stats.URL = url

	start := time.Now()
	resp, err := http.Get(url)
	stats.StatusCode = resp.StatusCode
	stats.ResponseTime = time.Since(start)
	stats.Error = err
	if err != nil {
		return stats
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Stats{URL: url, Error: err}
	}
	bodyStr := string(body)
	stats.WordCount = countWords(bodyStr)
	stats.ImageCount = countImages(bodyStr)
	stats.LinkCount = countLinks(bodyStr)

	return stats
}

func getUrl(url string) string {
	if !strings.HasPrefix(url, "http") {
		url = "http://" + url
	}

	return url
}

func countWords(body string) int {
	return len(strings.Fields(string(body)))
}

func countImages(body string) int {
	return strings.Count(string(body), "<img")
}

func countLinks(body string) int {
	return strings.Count(string(body), "<a")
}
