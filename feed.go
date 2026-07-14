package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/ashy558/bootdev-gator/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
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

func (i *RSSItem) String() string {
	fields := []string{}
	fields = append(fields, fmt.Sprintf(" * Title: %s", i.Title))
	fields = append(fields, fmt.Sprintf(" * Description: %s", i.Description))
	fields = append(fields, fmt.Sprintf(" * Link: %s", i.Link))
	fields = append(fields, fmt.Sprintf(" * PubDate: %s", i.PubDate))
	fields = append(fields, "=====================================")
	return strings.Join(fields, "\n")
}

func (i *RSSItem) Unescape() {
	i.Title = html.UnescapeString(i.Title)
	i.Description = html.UnescapeString(i.Description)
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	client := http.DefaultClient
	var feed RSSFeed
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)
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

func scrapeFeeds(s *state) error {
	const feedTimestampFormat = time.RFC1123Z
	ctx := context.Background()
	feed, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		return fmt.Errorf("could not fetch next feed: %s", err)
	}
	_, err = s.db.MarkFeedFetched(ctx, feed.ID)
	if err != nil {
		return fmt.Errorf("could not mark feed fetched: %s", err)
	}
	rssFeed, err := fetchFeed(ctx, feed.Url)
	if err != nil {
		return fmt.Errorf("could not fetch feed from url: %s", err)
	}
	rawPosts := rssFeed.Channel.Item
	if len(rawPosts) == 0 {
		fmt.Println("The RSS Feed is empty.")
	}
	fmt.Printf("Fetched %d Items from the %s RSS Feed:\n", len(rawPosts), rssFeed.Channel.Title)
	for i, item := range rawPosts {
		parsedTimestamp, err := time.Parse(feedTimestampFormat, item.PubDate)
		if err != nil {
			fmt.Printf("Post %d has invalid timestamp: %v: %s\n", i, item.PubDate, err)
		}
		args := database.CreatePostParams{
			ID:          uuid.New(),
			Title:       item.Title,
			Url:         item.Link,
			Description: item.Description,
			PublishedAt: parsedTimestamp,
			FeedID:      feed.ID,
		}
		post, err := s.db.CreatePost(ctx, args)
		if err != nil {
			if pqErr, ok := err.(*pq.Error); ok {
				if pqErr.Code.Name() == "unique_violation" {
					fmt.Println()
					fmt.Printf("%d.\n", i+1)
					fmt.Println("post already present in database")
					fmt.Println("=====================================")
					continue
				}
			}
			return fmt.Errorf("could not create new post: %s", err)
		}
		fmt.Println()
		fmt.Printf("%d.\n", i+1)
		fmt.Println(stringifyPost(post))
		fmt.Println()
	}
	return nil
}

func stringifyPost(post database.Post) string {
	fields := []string{}
	fields = append(fields, fmt.Sprintf(" * ID: %s", post.ID))
	fields = append(fields, fmt.Sprintf(" * Title: %s", post.Title))
	fields = append(fields, fmt.Sprintf(" * URL: %s", post.Url))
	fields = append(fields, fmt.Sprintf(" * Description: %s", post.Description))
	fields = append(fields, fmt.Sprintf(" * PublishedAt: %v", post.PublishedAt))
	fields = append(fields, fmt.Sprintf(" * FeedID: %s", post.FeedID))
	fields = append(fields, "=====================================")
	return strings.Join(fields, "\n")
}
