package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/EluxRed/Gator/internal/database"
)

func handlerBrowse(s *state, cmd command, currentUsr database.User) error {
	limit := 2
	if len(cmd.Args) != 0 {
		if len(cmd.Args) != 1 {
			return fmt.Errorf("parameters passed: %v\nparameters expected: 0 or 1", len(cmd.Args))
		} else {
			arg, err := strconv.Atoi(cmd.Args[0])
			if err != nil {
				return err
			}
			if arg == 0 {
				fmt.Println("0 posts fetched")
				return nil
			}
			limit = arg
		}
	}
	postsParams := database.GetPostsForUserParams{UserID: currentUsr.ID, Limit: int32(limit)}
	posts, err := s.dbPtr.GetPostsForUser(context.Background(), postsParams)
	if err != nil {
		return err
	}
	for i := range posts {
		printPost(posts[i])
	}
	return nil
}

func printPost(post database.GetPostsForUserRow) {
	fmt.Printf("* Title:           %s\n", post.Title)
	fmt.Printf("* URL:             %s\n", post.Url)
	fmt.Printf("* From feed:       %s\n", post.FeedName)
	fmt.Printf("* Published at:    %s\n", post.PublishedAt.Time)
	fmt.Printf("* Description:     %s\n", post.Description.String)
	fmt.Println()
}
