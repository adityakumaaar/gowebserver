package main

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/adityakumaaar/gowebserver/internal/database"
	"github.com/google/uuid"
)

func startScraping(
	db *database.Queries,
	concurreny int,
	timeBetweenRequest time.Duration,
) {
	log.Printf("Scraping on %v goroutines every %s duration", concurreny, timeBetweenRequest)
	ticker := time.NewTicker(timeBetweenRequest)

	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(context.Background(), int32(concurreny))
		if err != nil {
			log.Println("error fetching feeds: ", err)
			continue
		}

		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)
			go scrapeFeed(wg, feed, db)
		}
		wg.Wait()
	}

}

func scrapeFeed(wg *sync.WaitGroup, feed database.Feed, db *database.Queries) {
	defer wg.Done()
	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Println("error marking feed as fetched: ", err)
		return
	}
	redditFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Println("error fetching feed: ", err)
		return
	}

	for _, item := range redditFeed.Entries {
		desc := sql.NullString{}
		if item.Content != "" {
			desc.String = item.Content
			desc.Valid = true
		}
		_, err := db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       item.Title,
			Description: desc,
			Url:         item.Link.Href,
			FeedID:      feed.ID,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			}
			log.Println("error saving post in DB: ", err)
		}

	}

	log.Printf("Feed %v collected, found %v posts", feed.Name, len(redditFeed.Entries))

}
