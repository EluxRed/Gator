package main

import (
	"context"
	"fmt"
)

func handlerAgg(s *state, cmd command) error {
	if err := checkParams(cmd.Args, 0); err != nil {
		return err
	}
	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}
	fmt.Println(feed)
	return nil
}
