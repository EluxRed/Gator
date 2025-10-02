package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/EluxRed/Gator/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

func handlerAgg(s *state, cmd command) error {
	if err := checkParams(cmd.Args, 1); err != nil {
		return err
	}
	timeBetweenRequests, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return err
	}
	ticker := time.NewTicker(timeBetweenRequests)
	fmt.Printf("Collecting feeds every %v\n\n", timeBetweenRequests)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func scrapeFeeds(s *state) error {
	feed, err := s.dbPtr.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}
	if err := s.dbPtr.MarkFeedFetched(context.Background(), feed.ID); err != nil {
		return err
	}
	rssFeed, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		return err
	}
	for i := range rssFeed.Channel.Item {
		item := rssFeed.Channel.Item[i]
		var pubDate sql.NullTime
		if date, err := parseDate(item.PubDate); err == nil {
			pubDate = sql.NullTime{
				Time:  date,
				Valid: true,
			}
		} else {
			return err
		}
		post := database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Title:     item.Title,
			Url:       item.Link,
			Description: sql.NullString{
				String: item.Description,
				Valid:  true,
			},
			PublishedAt: pubDate,
			FeedID:      feed.ID,
		}
		_, err = s.dbPtr.CreatePost(context.Background(), post)
		if err != nil {
			if e, ok := err.(*pq.Error); ok && string(e.Code) == "23505" { //23505 is the unique-violation code in PostgreSql
				fmt.Printf("duplicate URL (%s) for post with title '%s': post creation skipped\n", post.Url, post.Title)
				err = nil
			} else {
				return err
			}
		}
		fmt.Printf("Post with title '%s' from feed '%s' saved successfuly\n", post.Title, feed.Name)
	}
	return nil
}

func parseDate(date string) (time.Time, error) {
	layouts := []string{time.RFC1123Z, time.RFC1123, time.RFC850, time.RFC3339, time.RFC822, time.RFC822Z, time.Layout, time.ANSIC, time.UnixDate, time.RubyDate, time.RFC3339Nano}
	var err error
	var parsedDate time.Time
	for i := range layouts {
		parsedDate, err = time.Parse(layouts[i], date)
		if err == nil {
			return parsedDate, err
		} else {
			err = fmt.Errorf("error using layout %s: %w", layouts[i], err)
		}
	}
	return parsedDate, err
}
