package ratelimiter

import (
	"encoding/json"
	"time"
	"x-straight-check/cache"
	"x-straight-check/pkg/gemini"
	"x-straight-check/pkg/xscraper"
)

func GetUserAnalysis(username, lang string) (*gemini.UserStraightnessAnalysis, error) {
	var res *gemini.UserStraightnessAnalysis
	key := username + "#" + lang

	cachedResp, err := cache.Get(key)
	if err == nil {
		err = json.Unmarshal([]byte(cachedResp), res)
		if err != nil {
			return nil, err
		}
		return res, nil
	}

	posts, user, err := xscraper.GetUserPosts(username, lang, 35)
	if err != nil {
		return nil, err
	}

	res, err = gemini.CheckUserStraightness(posts, user)
	if err != nil {
		return nil, err
	}

	resJson, err := json.Marshal(*res)
	if err != nil {
		return nil, err
	}

	err = cache.SetWithTTL(key, string(resJson), time.Hour*3)
	if err != nil {
		return nil, err
	}

	return res, nil
}
