// Golang RSS Reader Package (rssreader.go)

package rssreader

import (
	"encoding/xml"
    "errors"
    "net/http"
    "sync"
    "time"
)

type RssItem struct {
	Title       string    `json:"title"`
	Source      string    `json:"source"`
	SourceURL   string    `json:"source_url"`
	Link        string    `json:"link"`
	PublishDate time.Time `json:"publish_date"`
	Description string    `json:"description"`
}

type rssFeed struct {
	Channel struct {
		Title string `xml:"title"`
		Items []struct {
			Title       string `xml:"title"`
			Link        string `xml:"link"`
			Description string `xml:"description"`
			PubDate     string `xml:"pubDate"`
		} `xml:"item"`
	} `xml:"channel"`
}

func Parse(urls []string) ([]RssItem, error) {
	var wg sync.WaitGroup
	results := make([]RssItem, 0)
	var mu sync.Mutex

	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			resp, err := http.Get(url)
			if err != nil {
				return
			}
			defer resp.Body.Close()

			var feed rssFeed
			if err := xml.NewDecoder(resp.Body).Decode(&feed); err != nil {
				return
			}

			for _, item := range feed.Channel.Items {
				publishDate, _ := time.Parse(time.RFC1123, item.PubDate)
				rssItem := RssItem{
					Title:       item.Title,
					Source:      feed.Channel.Title,
					SourceURL:   url,
					Link:        item.Link,
					PublishDate: publishDate,
					Description: item.Description,
				}
				mu.Lock()
				results = append(results, rssItem)
				mu.Unlock()
			}
		}(url)
	}

	wg.Wait()
	if len(results) == 0 {
		return nil, errors.New("no items parsed")
	}
	return results, nil
}
