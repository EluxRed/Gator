package main

import (
	"context"
	"fmt"
	"time"

	"github.com/EluxRed/Gator/internal/database"
	"github.com/google/uuid"
)

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("parameters passed: %v\nparameters expected: %v - Name of the feed and URL", len(cmd.Args), 2)
	}
	feedName := cmd.Args[0]
	feedURL := cmd.Args[1]
	currentUsr, err := s.dbPtr.GetUser(context.Background(), s.configPtr.Current_User_Name)
	if err != nil {
		return err
	}
	newFeed := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      feedName,
		Url:       feedURL,
		UserID:    currentUsr.ID,
	}
	feed, err := s.dbPtr.CreateFeed(context.Background(), newFeed)
	if err != nil {
		return err
	}
	fmt.Println("Feed created successfully:")
	printFeed(feed)
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	if len(cmd.Args) > 0 {
		return fmt.Errorf("parameters passed: %v\nparameters expected: %v", len(cmd.Args), 0)
	}
	feeds, err := s.dbPtr.GetFeeds(context.Background())
	if err != nil {
		return err
	}
	printFeeds(feeds)
	return nil
}

func printFeed(feed database.Feed) {
	fmt.Printf("* ID:            %s\n", feed.ID)
	fmt.Printf("* Created:       %v\n", feed.CreatedAt)
	fmt.Printf("* Updated:       %v\n", feed.UpdatedAt)
	fmt.Printf("* Name:          %s\n", feed.Name)
	fmt.Printf("* URL:           %s\n", feed.Url)
	fmt.Printf("* UserID:        %s\n", feed.UserID)
}

func printFeeds(feeds []database.GetFeedsRow) {
	for i := range feeds {
		fmt.Printf("* FeedName:       %s\n", feeds[i].FeedName)
		fmt.Printf("* FeedUrl:        %s\n", feeds[i].FeedUrl)
		fmt.Printf("* UserName:       %s\n", feeds[i].UserName)
		fmt.Println("==================================================")
	}
}
