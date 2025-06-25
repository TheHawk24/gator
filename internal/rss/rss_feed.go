package rss

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	//	"log"
	"net/http"
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

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {

	//Make a new request
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		err := fmt.Errorf("Failed to create a request")
		return nil, err
	}
	req.Header.Set("User-Agent", "gator")

	var client http.Client
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		err := fmt.Errorf("Failed to send the request")
		return nil, err
	}

	defer resp.Body.Close()

	// Read the response from server
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		err := fmt.Errorf("Failed to read the body of the response")
		return nil, err
	}

	var feed RSSFeed
	err = xml.Unmarshal(data, &feed)
	if err != nil {
		fmt.Println(err)
		err := fmt.Errorf("Failed to decode data into feed")
		return nil, err
	}

	//Unescape entities
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

	for i, v := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(v.Title)
		feed.Channel.Item[i].Description = html.UnescapeString(v.Description)
	}

	return &feed, nil
}
