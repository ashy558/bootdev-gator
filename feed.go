package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
)

const (
	userAgentHeader = "User-Agent"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func (f *RSSFeed) Unescape() {
	f.Channel.Title = html.UnescapeString(f.Channel.Title)
	f.Channel.Description = html.UnescapeString(f.Channel.Title)
	for i := range f.Channel.Item {
		f.Channel.Item[i].Unescape()
	}
}

func (i *RSSItem) Unescape() {
	i.Title = html.UnescapeString(i.Title)
	i.Description = html.UnescapeString(i.Description)
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	client := http.DefaultClient
	var feed RSSFeed
	request, err := http.NewRequestWithContext(context.Background(), http.MethodGet, feedURL, nil)
	if err != nil {
		return &feed, fmt.Errorf("could not create request: %s", err)
	}
	request.Header.Set(userAgentHeader, "gator")
	res, err := client.Do(request)
	if err != nil {
		return &feed, fmt.Errorf("could not perform request: %s", err)
	}
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return &feed, fmt.Errorf("could not read response: %s", err)
	}
	if err := xml.Unmarshal(data, &feed); err != nil {
		return &feed, fmt.Errorf("could not unmarshal response: %s", err)
	}
	feed.Unescape()
	return &feed, nil
}
