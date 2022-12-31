package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mendesbarreto/friday/api/dto"
)

type TweetsRespose struct {
	Data []Tweet `json:"data"`
	Meta Meta    `json:"meta"`
}
type Tweet struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Text      string    `json:"text"`
}
type Meta struct {
	ResultCount int    `json:"result_count"`
	NextToken   string `json:"next_token"`
}

type TweetSentment struct {
	Tweet Tweet `json:"tweet"`
	Score int8  `json:"score"`
}

type TweetSentmentResponse struct {
	Data []TweetSentment `json:"data"`
}

const twitterUrl string = "https://api.twitter.com/2/lists/%s/tweets?tweet.fields=created_at"

var twitterToken string = os.Getenv("TWITTER_API_TOKEN")

func getTweetSentment(message string) uint8 {
	sentiment, err := GetTextSentiment(message)

	if err != nil {
		fmt.Println("Twitter was not able to score")
		return 3
	}

	return sentiment
}

func GetLatestTweetsByListId(listId string) (*TweetsRespose, error) {
	client := http.Client{}

	request, err := http.NewRequest("GET", fmt.Sprintf(twitterUrl, listId), nil)
	if err != nil {
		return nil, err
	}

	twitterBearer := fmt.Sprintf("Bearer %s", twitterToken)

	request.Header = http.Header{
		"Authorization": {twitterBearer},
	}

	res, err := client.Do(request)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var result TweetsRespose
	resBody, err := ioutil.ReadAll(res.Body)

	if err := json.Unmarshal(resBody, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func getLatestTweetsSentmentByListId(listId string) (*[]TweetSentment, error) {
	tweets, err := GetLatestTweetsByListId(listId)

	if err != nil {
		return nil, err
	}

	tweetSentmentList := make([]TweetSentment, len(tweets.Data))

	var wg sync.WaitGroup
	for index, data := range tweets.Data {
		wg.Add(1)
		go func(index int, tweet Tweet) {
			tweetSentmentList[index] = TweetSentment{
				Tweet: tweet,
				Score: int8(getTweetSentment(tweet.Text)),
			}
			wg.Done()
		}(index, data)
	}

	wg.Wait()

	return &tweetSentmentList, nil
}

func GetTweetsFromToday() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		listId := ctx.Params("id")

		if len(listId) == 0 {
			return dto.BadRequest(ctx, "Missing list id on query params")
		}

		tweetSentmentList, err := getLatestTweetsSentmentByListId(listId)

		if err != nil {
			return dto.InternalServerError(ctx, err.Error())
		}

		response := TweetSentmentResponse{Data: *tweetSentmentList}

		return ctx.Status(200).JSON(response)
	}
}

type AverageMood struct {
	Total      int32   `json:"total"`
	Percentage float32 `json:"percentage"`
}

type AverageTweetsMoodResponse struct {
	Positive AverageMood `json:"positive"`
	Negative AverageMood `json:"negative"`
	Total    int32       `json:"total"`
}

func GetAverageTweetsMoodByList() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		listId := ctx.Params("id")

		if len(listId) == 0 {
			return dto.BadRequest(ctx, "Missing list id on query Params")
		}

		tweetSentmentList, err := getLatestTweetsSentmentByListId(listId)

		if err != nil {
			return dto.InternalServerError(ctx, err.Error())
		}

		var positiveTweets []TweetSentment
		var negativeTweets []TweetSentment

		for _, value := range *tweetSentmentList {
			if value.Score == 0 {
				negativeTweets = append(negativeTweets, value)
			} else {
				positiveTweets = append(positiveTweets, value)
			}
		}

		totalTweetsCount := int32(len(*tweetSentmentList))
		positiveTweetsCount := int32(len(positiveTweets))
		negativeTweetsCount := int32(len(negativeTweets))

		response := AverageTweetsMoodResponse{
			Positive: AverageMood{
				Total:      positiveTweetsCount,
				Percentage: float32(positiveTweetsCount) / float32(totalTweetsCount),
			},
			Negative: AverageMood{
				Total:      negativeTweetsCount,
				Percentage: float32(negativeTweetsCount) / float32(totalTweetsCount),
			},
			Total: totalTweetsCount,
		}

		return ctx.Status(200).JSON(response)
	}
}
