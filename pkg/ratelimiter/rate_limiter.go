package ratelimiter

import (
	"encoding/json"
	"time"
	"x-straight-check/cache"
	"x-straight-check/log"
	"x-straight-check/pkg/gemini"
	"x-straight-check/pkg/xscraper"
)

func GetUserAnalysis(username, lang string) (*gemini.UserStraightnessAnalysis, error) {
	var res gemini.UserStraightnessAnalysis
	key := username + "#" + lang

	cachedRes, err := cache.Get(key)
	if err == nil {
		err = json.Unmarshal([]byte(cachedRes), &res)
		if err != nil {
			log.Errorf("Failed decoded cached response, %v\n", err)
			goto noResp
		}
		return &res, nil
	}
noResp:

	posts, user, err := xscraper.GetUserPosts(username, lang, 35)
	if err != nil {
		log.Errorf("X Scrapper failed, %v\n", err)
		return nil, err
	}

	res1, err := gemini.CheckUserStraightness(posts, user)
	if err != nil {
		log.Errorf("AI failed, %v\n", err)
		return nil, err
	}
	res = *res1

	resJson, err := json.Marshal(res)
	if err != nil {
		log.Errorf("Failed marshalling json, %v\n", err)
		return nil, err
	}

	err = cache.SetWithTTL(key, string(resJson), time.Hour*3)
	if err != nil {
		log.Errorf("Failed setting cache, %v\n", err)
		return nil, err
	}

	return &res, nil
}
