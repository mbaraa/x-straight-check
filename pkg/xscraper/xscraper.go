package xscraper

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"x-straight-check/config"
)

type ApifyScrapeXActorRequest struct {
	Author             string   `json:"author"`
	CustomMapFunction  string   `json:"customMapFunction"`
	End                string   `json:"end"`
	IncludeSearchTerms bool     `json:"includeSearchTerms"`
	MaxItems           int16    `json:"maxItems"`
	MinimumFavorites   int16    `json:"minimumFavorites"`
	MinimumReplies     int16    `json:"minimumReplies"`
	MinimumRetweets    int16    `json:"minimumRetweets"`
	PlaceObjectId      string   `json:"placeObjectId"`
	Sort               string   `json:"sort"`
	Start              string   `json:"start"`
	TweetLanguage      string   `json:"tweetLanguage"`
	TwitterHandles     []string `json:"twitterHandles"`
}

type XUser struct {
	Type               string `json:"type"`
	UserName           string `json:"userName"`
	Id                 string `json:"id"`
	Name               string `json:"name"`
	IsVerified         bool   `json:"isVerified"`
	IsBlueVerified     bool   `json:"isBlueVerified"`
	ProfilePicture     string `json:"profilePicture"`
	CoverPicture       string `json:"coverPicture"`
	Description        string `json:"description"`
	Location           string `json:"location"`
	Followers          int    `json:"followers"`
	Following          int    `json:"following"`
	CreatedAt          string `json:"createdAt"`
	LikesCount         int    `json:"favouritesCount"`
	HasCustomTimelines bool   `json:"hasCustomTimelines"`
	IsTranslator       bool   `json:"isTranslator"`
	MediaCount         int    `json:"mediaCount"`
	PostsCount         int    `json:"statusesCount"`
}

type XPost struct {
	Url           string `json:"url,omitempty"`
	Id            string `json:"id"`
	Text          string `json:"text"`
	FullText      string `json:"fullText"`
	Source        string `json:"source"`
	RetweetCount  int    `json:"retweetCount"`
	ReplyCount    int    `json:"replyCount"`
	LikeCount     int    `json:"likeCount"`
	QuoteCount    int    `json:"quoteCount"`
	CreatedAt     string `json:"createdAt"`
	BookmarkCount int    `json:"bookmarkCount"`
	IsRetweet     bool   `json:"isRetweet"`
	IsQuote       bool   `json:"isQuote"`
	Author        *XUser `json:"author,omitempty"`
}

func GetUserPosts(username, lang string, postsCount int16) ([]XPost, XUser, error) {
	reqBody := ApifyScrapeXActorRequest{
		Author:             "apify",
		CustomMapFunction:  "(object) => { return {...object} }",
		Start:              "2021-07-02",
		End:                "2021-07-02",
		IncludeSearchTerms: false,
		MaxItems:           postsCount,
		MinimumFavorites:   5,
		MinimumReplies:     5,
		MinimumRetweets:    5,
		PlaceObjectId:      "96683cc9126741d1",
		Sort:               "Latest",
		TweetLanguage:      lang,
		TwitterHandles:     []string{username},
	}

	respBodyBuf := bytes.NewBuffer(nil)
	err := json.NewEncoder(respBodyBuf).Encode(reqBody)
	if err != nil {
		return nil, XUser{}, err
	}

	resp, err := http.Post(fmt.Sprintf("https://api.apify.com/v2/acts/apidojo~tweet-scraper/run-sync-get-dataset-items?token=%s&timeout=30&memory=256", config.Env().ApifyToken), "application/json", respBodyBuf)
	if err != nil {
		return nil, XUser{}, err
	}
	if resp.StatusCode != http.StatusCreated {
		return nil, XUser{}, errors.New("scraping status != 201")
	}

	posts := []XPost{}
	err = json.NewDecoder(resp.Body).Decode(&posts)
	if err != nil {
		return nil, XUser{}, err
	}

	var user XUser
	for i := range posts {
		if posts[i].Author != nil {
			user = *posts[i].Author
		}
		posts[i].Author = nil
	}

	return posts, user, nil
}
