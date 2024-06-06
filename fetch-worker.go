package main

import (
	"context"
	"database/sql"
	"log"
	"sync"
	"time"

	"github.com/KMalkowski/rss-feed/internal/database"
	feedfetcher "github.com/KMalkowski/rss-feed/internal/feed_fetcher"
	"github.com/google/uuid"
)

func (a *apiConfig) FetchWorker(n int) error {
	for {
		feeds, err := a.DB.GetNextToFetch(context.TODO(), int32(n))
		if err != nil {
			return err
		}

		var wg sync.WaitGroup
		for _, feed := range feeds {
			wg.Add(1)

			go func(feed database.Feed) {
				defer wg.Done()

				_, err := a.DB.MarkFetched(context.TODO(), database.MarkFetchedParams{
					LastFetchedAt: sql.NullTime{Time: time.Now(), Valid: true},
					ID:            feed.ID,
				})

				if err != nil {
					log.Println("could not mark feed as fetched" + feed.Url)
					return
				}

				xml, err := feedfetcher.FetchFeed(feed.Url)
				if err != nil {
					log.Println("invalid feed url" + feed.Url)
					return
				}

				for _, item := range xml.Channel.Item {
					log.Println(item.Title)
					pubTime, err := time.Parse(time.DateOnly, item.PubDate)
					if err != nil {
						pubTime = time.Now()
					}

					_, err = a.DB.CreatePost(context.TODO(), database.CreatePostParams{
						ID:          uuid.New(),
						FeedID:      feed.ID,
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
						PublishedAt: pubTime,
						Title:       item.Title,
						Description: sql.NullString{
							String: item.Description,
						},
						Url: item.Link,
					})

					if err != nil {
						log.Fatalln(err.Error())
						log.Println("cound not save post in the db" + feed.Url)
					}
				}

				return
			}(feed)
		}

		time.Sleep(time.Second * 60)
	}
}
