package ratelimiter

import (
	"encoding/base64"
	"encoding/json"
	"time"
	"x-straight-check/cache"
	"x-straight-check/log"
	"x-straight-check/pkg/gemini"
	"x-straight-check/pkg/xscraper"
)

func GetUserAnalysis(username, lang string) (*gemini.UserStraightnessAnalysis, error) {
	var res *gemini.UserStraightnessAnalysis
	key := username + "#" + lang

	cachedRes, err := cache.Get(key)
	if err == nil {
		decRes, err := base64.StdEncoding.DecodeString(cachedRes)
		if err != nil {
			log.Errorf("Failed decoding cached response, %v\n", err)
			goto noResp
		}

		err = json.Unmarshal([]byte(decRes), res)
		if err != nil {
			log.Errorf("Failed unmarshalling decoded cached response, %v\n", err)
			goto noResp
		}
		return res, nil
	}
noResp:

	posts, user, err := xscraper.GetUserPosts(username, lang, 35)
	if err != nil {
		log.Errorf("X Scrapper failed, %v\n", err)
		return nil, err
	}

	res, err = gemini.CheckUserStraightness(posts, user)
	if err != nil {
		log.Errorf("AI failed, %v\n", err)
		return nil, err
	}

	resJson, err := json.Marshal(*res)
	if err != nil {
		log.Errorf("Failed marshalling json, %v\n", err)
		return nil, err
	}

	encResJson := base64.StdEncoding.EncodeToString(resJson)

	err = cache.SetWithTTL(key, string(encResJson), time.Hour*3)
	if err != nil {
		log.Errorf("Failed setting cache, %v\n", err)
		return nil, err
	}

	return res, nil
}
