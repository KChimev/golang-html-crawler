package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"sync"
	"time"

	"golang.org/x/net/html"
)

func main() {
	start := flag.String("start", "https://stackoverflow.com/", "The starting point of the crawler")
	maxWorkers := flag.Int("workers", 10, "Maximum ammount of workers")
	timeLimit := flag.Int("time", 60, "Application lifetime in seconds")
	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(*timeLimit))
	defer cancel()
	jobs := make(chan string, 100)
	var wg sync.WaitGroup
	var mu sync.Mutex
	var visited sync.Map
	fmt.Printf("%s\n", *start)
	for i := 0; i < *maxWorkers; i++ {
		wg.Add(1)
		go func() {
			fmt.Printf("Worker %d is working...\n", i)
			defer wg.Done()
			for job := range jobs {
				select {
				case <-ctx.Done():
					fmt.Printf("Done..\n")
					return
				default:
					mu.Lock()
					_, loaded := visited.LoadOrStore(job, nil)
					if loaded {
						fmt.Printf("Already checked %s\n", job)
						continue
					}
					mu.Unlock()
					fmt.Printf("Processing %s... \n", job)
					res, err := processDocument(job)

					mu.Lock()
					if err != nil {
						fmt.Printf("Error processing job: %v\n", err)
						visited.Store(job, false) // Store as false to be reiterated later if needed
					}
					for _, url := range res {
						fmt.Printf("FInished working: %s \n", url)
						jobs <- url
					}
					visited.Store(job, true)
					mu.Unlock()
				}
			}
		}()
	}

	jobs <- *start
	<-ctx.Done()
	close(jobs)
	wg.Wait()
}

func processDocument(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	tokenizer := html.NewTokenizer(resp.Body)
	var links []string

	for {
		tokenType := tokenizer.Next()
		switch tokenType {
		case html.ErrorToken:
			return links, nil
		case html.StartTagToken:
			token := tokenizer.Token()
			if token.Data == "a" {
				for _, attr := range token.Attr {
					if attr.Key == "href" {
						links = append(links, attr.Val)
					}
				}
			}
		}
	}
}
