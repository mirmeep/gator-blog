package main

import (
	"fmt"
	"net/http"
	"io"
	"encoding/xml"
	"context"
	"html"
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

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "", feedURL, nil)
	if err != nil {
		return &RSSFeed{}, err
	}

	req.Header.Set("User-Agent", "gator")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return &RSSFeed{}, err
	}

	defer res.Body.Close()
	if res.StatusCode > 299 {
		return &RSSFeed{}, fmt.Errorf("RSS feed not found")
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return &RSSFeed{}, err
	}

	var rss RSSFeed
	if err := xml.Unmarshal(data, &rss); err != nil {
		return &RSSFeed{}, err
	}

	rss.Channel.Title = html.UnescapeString(rss.Channel.Title)
	rss.Channel.Description = html.UnescapeString(rss.Channel.Description)
	
	for index := range rss.Channel.Item {
		rss.Channel.Item[index].Title = html.UnescapeString(rss.Channel.Item[index].Title)
		rss.Channel.Item[index].Description = html.UnescapeString(rss.Channel.Item[index].Description)
	}

	return &rss, nil
}

func handlerAgg(s *state, cmd command) error {
	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		fmt.Errorf("cannot fetch feed: %w", err)
	}

	fmt.Println(feed)
	return nil
}