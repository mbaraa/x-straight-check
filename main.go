package main

import (
	"os"
	"x-straight-check/pkg/gemini"
	"x-straight-check/pkg/xscraper"
)

func main() {
	posts, user, err := xscraper.GetUserPosts(os.Args[1], "en", 42)
	if err != nil {
		panic(err)
	}

	_, err = gemini.CheckUserStraightness(posts, user)
	if err != nil {
		panic(err)
	}
}
