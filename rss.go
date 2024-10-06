package main

import (
	"encoding/xml"
	"io"
	"net/http"
	"time"
)

// RedditFeed represents a feed from Reddit
type RedditFeed struct {
	XMLName  string    `xml:"feed"`
	Title    string    `xml:"title"`
	Subtitle string    `xml:"subtitle"`
	Updated  time.Time `xml:"updated,attr"`
	Link     []Link    `xml:"link"`
	Entries  []Entry   `xml:"entry"`
}

// Link represents a link in the feed
type Link struct {
	Rel  string `xml:"rel,attr"`
	Href string `xml:"href,attr"`
}

// Entry represents an entry (post) in the feed
type Entry struct {
	ID        string     `xml:"id"`
	Published string     `xml:"published,attr"`
	Updated   time.Time  `xml:"updated,attr"`
	Title     string     `xml:"title"`
	Author    Author     `xml:"author"`
	Content   string     `xml:"content"`
	Link      Link       `xml:"link"`
	Category  []Category `xml:"category"`
	Thumbnail Thumbnail  `xml:"media:thumbnail,attr"`
}

// Author represents the author of an entry
type Author struct {
	Name string `xml:"name"`
	URI  string `xml:"uri"`
}

// Category represents a category of an entry
type Category struct {
	Term  string `xml:"term,attr"`
	Label string `xml:"label,attr"`
}

// Thumbnail represents the thumbnail of an entry
type Thumbnail struct {
	URL string `xml:"url,attr"`
}

func urlToFeed(url string) (RedditFeed, error) {
	httpClient := http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := httpClient.Get(url)
	if err != nil {
		return RedditFeed{}, err
	}
	defer resp.Body.Close()

	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return RedditFeed{}, err
	}

	rFeed := RedditFeed{}
	err = xml.Unmarshal(dat, &rFeed)
	if err != nil {
		return RedditFeed{}, err
	}

	return rFeed, nil

}
