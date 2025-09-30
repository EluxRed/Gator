package main

import (
	"context"
	"fmt"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) > 0 {
		return fmt.Errorf("parameters passed: %v\nparameters expected: %v", len(cmd.Args), 0)
	}
	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}
	fmt.Println(feed)
	return nil
}
