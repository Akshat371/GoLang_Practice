package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Result struct {
	URL        string
	StatusCode int
	Err        error
}

func checkURL(url string, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	response, err := client.Get(url)

	if err != nil {
		results <- Result{
			URL: url,
			Err: err,
		}
		return
	}

	defer response.Body.Close()

	results <- Result{
		URL:        url,
		StatusCode: response.StatusCode,
	}
}

func main() {
	urls := []string{
		"https://google.com",
		"https://github.com",
		"https://golang.org",
		"https://amazon.com",
		"https://example.com",
	}

	results := make(chan Result, len(urls))

	var wg sync.WaitGroup

	start := time.Now()

	for _, url := range urls {
		wg.Add(1)

		go checkURL(url, results, &wg)
	}

	wg.Wait()
	close(results)

	for result := range results {
		if result.Err != nil {
			fmt.Println(result.URL, "Error:", result.Err)
		} else {
			fmt.Println(result.URL, "Status:", result.StatusCode)
		}
	}

	fmt.Println("Total time:", time.Since(start))
}
