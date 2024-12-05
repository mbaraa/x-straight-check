package gemini

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"x-straight-check/config"
	"x-straight-check/pkg/xscraper"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

var (
	gemini *genai.GenerativeModel
	ctx    = context.Background()
)

func init() {
	client, err := genai.NewClient(ctx, option.WithAPIKey(config.Env().GeminiToken))
	if err != nil {
		panic(err)
	}
	gemini = client.GenerativeModel("gemini-1.5-flash")
}

type UserStraightnessAnalysis struct {
	Straightness  float32 `json:"straightness"`
	ReasonOfScore string  `json:"reason_of_score"`
}

func CheckUserStraightness(posts []xscraper.XPost, user xscraper.XUser) (*UserStraightnessAnalysis, error) {
	follwRatio := float64(user.Followers) / float64(user.Following)
	if user.Description == "" {
		user.Description = "user has no description"
	}

	postsJson, err := json.Marshal(posts)
	if err != nil {
		return nil, err
	}

	prompt := fmt.Sprintf(
		`Hey Gemini, take a look at this profile from x.com:
    - name: %s
    - username: %s
    - description: %s
    - profile picture link: %s
    - cover picture link: %s
    - followers: %d
    - following: %d
    - posts count: %d
    - follow ratio: %f
    - last %d posts (pre-formatted JSON): %s

Now with all this data, I need you to analyze it and report how straight the owner of this account appears to be. Make use of all the data provided above.

Make the output in strict JSON format with the following structure:

{"straightness": X.XXX, "reason_of_score": "your explanation"}

where:
 - "straightness" is a numeric value between 0.000 and 1.000 (with three decimal places).
 - "reason_of_score" is a string providing a PG-13 explanation for the score under 240 character.`,
		user.Name,
		user.UserName,
		user.Description,
		user.ProfilePicture,
		user.CoverPicture,
		user.Followers,
		user.Following,
		user.PostsCount,
		follwRatio,
		len(posts),
		postsJson,
	)

	resp, err := gemini.GenerateContent(context.Background(), genai.Text(prompt))
	if err != nil {
		return nil, err
	}
	if len(resp.Candidates) == 0 {
		return nil, errors.New("ai didn't respond well")
	}
	if len(resp.Candidates[0].Content.Parts) == 0 {
		return nil, errors.New("ai didn't respond well")
	}

	respJson, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}

	mdJsonStart := strings.Index(string(respJson), "```json")
	if mdJsonStart < 0 {
		return nil, errors.New("ai didn't respond well")
	}
	mdJsonEnd := strings.LastIndex(string(respJson), "```")
	if mdJsonEnd < 0 {
		return nil, errors.New("ai didn't respond well")
	}

	respJsonFixed := string(respJson[mdJsonStart+len("```json") : mdJsonEnd])
	respJsonFixed = strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(respJsonFixed, "\\n", ""), "\\t", ""), "\\\"", "\""), "\\\\\"", "\\\"")
	fmt.Println(respJsonFixed)

	var analysis UserStraightnessAnalysis
	err = json.Unmarshal([]byte(respJsonFixed), &analysis)
	if err != nil {
		return nil, err
	}

	return &analysis, nil
}
