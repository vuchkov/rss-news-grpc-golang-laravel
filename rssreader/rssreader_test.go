// Test file (rssreader_test.go)
package rssreader

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	urls := []string{
		"https://rss.cnn.com/rss/cnn_topstories.rss",
		// "https://rss.nytimes.com/services/xml/rss/nyt/HomePage.xml",
	}
	items, err := Parse(urls)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(items) == 0 {
		t.Fatal("Expected items, got none")
	}
	if !reflect.DeepEqual(items[0].SourceURL, urls[0]) {
		t.Errorf("Expected source URL %v, got %v", urls[0], items[0].SourceURL)
	}
}
