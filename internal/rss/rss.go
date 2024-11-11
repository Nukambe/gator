package rss

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"net/http"
)

type Feed struct {
	Channel struct {
		Title       string `xml:"title"`
		Link        string `xml:"link"`
		Description string `xml:"description"`
		Item        []Item `xml:"item"`
	} `xml:"channel"`
}

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func (feed Feed) unescape() {
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	for _, i := range feed.Channel.Item {
		i.Title = html.UnescapeString(i.Title)
		i.Description = html.UnescapeString(i.Description)
	}
}

func FetchFeed(ctx context.Context, feedURL string) (*Feed, error) {
	// create request
	req, errReq := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if errReq != nil {
		return nil, fmt.Errorf("unable to create rss request: %w", errReq)
	}

	// set headers
	req.Header.Add("User-Agent", "gator")

	// make request
	client := http.Client{}
	res, errRes := client.Do(req)
	if errRes != nil {
		return nil, fmt.Errorf("network error fetching feed: %w", errRes)
	}
	defer res.Body.Close()

	// read response
	var data Feed
	decoder := xml.NewDecoder(res.Body)
	if err := decoder.Decode(&data); err != nil {
		return nil, fmt.Errorf("unable to decode rss: %w", err)
	}

	data.unescape()
	return &data, nil
}
