// Golang RSS Reader Package (rssreader.go)

package rssreader

import (
	"encoding/xml"
	// "errors"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

// type RssItem struct {
// 	Title       string    `json:"title"`
// 	Source      string    `json:"source"`
// 	SourceURL   string    `json:"source_url"`
// 	Link        string    `json:"link"`
// 	PublishDate time.Time `json:"publish_date"`
// 	Description string    `json:"description"`
// }

type RssItem struct {
	Title       string    `xml:"title"`
	Source      string    `xml:"-"` // Not directly from RSS
	SourceURL   string    `xml:"-"`
	Link        string    `xml:"link"`
	PublishDate time.Time `xml:"pubDate"` // Adjust based on feed format
	Description string    `xml:"description"`
}

// type rssFeed struct {
// 	Channel struct {
// 		Title string `xml:"title"`
// 		Items []struct {
// 			Title       string `xml:"title"`
// 			Link        string `xml:"link"`
// 			Description string `xml:"description"`
// 			PubDate     string `xml:"pubDate"`
// 		} `xml:"item"`
// 	} `xml:"channel"`
// }

type rss struct {
	Channel struct {
		Items []RssItem `xml:"item"`
	} `xml:"channel"`
}

// func Parse(urls []string) ([]RssItem, error) {
// 	var wg sync.WaitGroup
// 	results := make([]RssItem, 0)
// 	var mu sync.Mutex

// 	for _, url := range urls {
// 		wg.Add(1)
// 		go func(url string) {
// 			defer wg.Done()
// 			resp, err := http.Get(url)
// 			if err != nil {
// 				return
// 			}
// 			defer resp.Body.Close()

// 			var feed rssFeed
// 			if err := xml.NewDecoder(resp.Body).Decode(&feed); err != nil {
// 				return
// 			}

// 			for _, item := range feed.Channel.Items {
// 				publishDate, _ := time.Parse(time.RFC1123, item.PubDate)
// 				rssItem := RssItem{
// 					Title:       item.Title,
// 					Source:      feed.Channel.Title,
// 					SourceURL:   url,
// 					Link:        item.Link,
// 					PublishDate: publishDate,
// 					Description: item.Description,
// 				}
// 				mu.Lock()
// 				results = append(results, rssItem)
// 				mu.Unlock()
// 			}
// 		}(url)
// 	}

// 	wg.Wait()
// 	if len(results) == 0 {
// 		return nil, errors.New("no items parsed")
// 	}
// 	return results, nil
// }

func Parse(urls []string) ([]RssItem, error) {
	var allItems []RssItem
	var wg sync.WaitGroup
	var mu sync.Mutex // Protect shared data
	errChan := make(chan error)

	for _, url := range urls {
			wg.Add(1)
			go func(url string) {
					defer wg.Done()

					resp, err := http.Get(url)
					if err != nil {
							errChan <- fmt.Errorf("error fetching %s: %w", url, err)
							return
					}
					defer resp.Body.Close()

					data, err := io.ReadAll(resp.Body)
					if err != nil {
							errChan <- fmt.Errorf("error reading body from %s: %w", url, err)
							return
					}

					var rssData rss
					err = xml.Unmarshal(data, &rssData)
					if err != nil {
							errChan <- fmt.Errorf("error unmarshaling XML from %s: %w", url, err)
							return
					}

					mu.Lock()
					for i := range rssData.Channel.Items {
							rssData.Channel.Items[i].SourceURL = url
							// Extract the source name from the URL (basic example)
							rssData.Channel.Items[i].Source = url // Or more sophisticated parsing
					}
					allItems = append(allItems, rssData.Channel.Items...)
					mu.Unlock()
			}(url)
	}

	go func() {
			wg.Wait()
			close(errChan)
	}()

	for err := range errChan {
	if err != nil {
		return nil, err
	}
}
	return allItems, nil
}
