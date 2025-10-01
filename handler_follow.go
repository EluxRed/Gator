package main

import (
	"context"
	"fmt"
	"time"

	"github.com/EluxRed/Gator/internal/database"
	"github.com/google/uuid"
)

func handlerFollow(s *state, cmd command, currentUsr database.User) error {
	if err := checkParams(cmd.Args, 1); err != nil {
		return err
	}
	url := cmd.Args[0]
	feed, err := s.dbPtr.GetFeedFromURL(context.Background(), url)
	if err != nil {
		return err
	}
	feedFollow := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    currentUsr.ID,
		FeedID:    feed.ID,
	}
	row, err := s.dbPtr.CreateFeedFollow(context.Background(), feedFollow)
	if err != nil {
		return err
	}
	fmt.Printf("Feed followed successfuly by %s:\n", row.UserName)
	printFeedFollow(row)
	return nil
}

func handlerFollowing(s *state, cmd command, currentUsr database.User) error {
	if err := checkParams(cmd.Args, 0); err != nil {
		return err
	}
	feeds, err := s.dbPtr.GetFeedFollowsForUser(context.Background(), currentUsr.Name)
	if err != nil {
		return err
	}
	if len(feeds) == 0 {
		fmt.Printf("The current user %s is not following any feed\n", currentUsr.Name)
		return nil
	}
	fmt.Printf("The current user %s is following:\n", currentUsr.Name)
	for i := range feeds {
		fmt.Println(feeds[i].FeedName)
	}
	return nil
}

func printFeedFollow(feedFollow database.CreateFeedFollowRow) {
	fmt.Printf("* ID:            %s\n", feedFollow.ID)
	fmt.Printf("* Feed Name:       %v\n", feedFollow.FeedName)
	fmt.Printf("* User Name:       %v\n", feedFollow.UserName)
}

func handlerUnfollow(s *state, cmd command, currentUsr database.User) error {
	if err := checkParams(cmd.Args, 1); err != nil {
		return err
	}
	url := cmd.Args[0]
	deleteParams := database.DeleteFeedFollowParams{UserID: currentUsr.ID, Url: url}
	if err := s.dbPtr.DeleteFeedFollow(context.Background(), deleteParams); err != nil {
		return err
	}
	fmt.Println("Feed unfollowed successfuly")
	return nil
}
