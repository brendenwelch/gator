package rss

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
)

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("failed to create feed request: %w", err)
	}
	req.Header.Set("User-Agent", "gator")

	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("feed request failed: %w", err)
	}
	defer res.Body.Close()

	var feed RSSFeed
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("failed to read response: %w", err)
	}
	if err := xml.Unmarshal(data, &feed); err != nil {
		return &RSSFeed{}, fmt.Errorf("failed to unmarshal response data: %w", err)
	}

	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	for i := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(feed.Channel.Item[i].Title)
		feed.Channel.Item[i].Description = html.UnescapeString(feed.Channel.Item[i].Description)
	}

	return &feed, nil
}
