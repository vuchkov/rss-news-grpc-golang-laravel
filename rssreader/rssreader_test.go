// Test file (rssreader_test.go)
package rssreader

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	urls := []string{"https://example.com/rss"}
	items, err := Parse(urls)
	if err != nil {
		t.Fatalf("Error occurred: %v", err)
	}
	if len(items) == 0 {
		t.Fatal("Expected items, got none")
	}
	if !reflect.DeepEqual(items[0].SourceURL, urls[0]) {
		t.Errorf("Expected source URL %v, got %v", urls[0], items[0].SourceURL)
	}
}
