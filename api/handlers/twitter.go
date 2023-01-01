package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/mendesbarreto/friday/api/dto"
	"github.com/mendesbarreto/friday/pkg/twitter"
)

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

func GetLatestTweetsByListId(listId string) (*dto.TweetsResposeBody, error) {
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

	var result dto.TweetsResposeBody
	resBody, err := ioutil.ReadAll(res.Body)

	if err := json.Unmarshal(resBody, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func getLatestTweetsSentmentByListId(listId string) (*[]twitter.TweetSentment, error) {
	tweets, err := GetLatestTweetsByListId(listId)

	if err != nil {
		return nil, err
	}

	tweetSentmentList := make([]twitter.TweetSentment, len(tweets.Data))

	var wg sync.WaitGroup
	for index, data := range tweets.Data {
		wg.Add(1)
		go func(index int, tweet twitter.Tweet) {
			tweetSentmentList[index] = twitter.TweetSentment{
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

		response := dto.TweetSentmentResponseBody{Data: *tweetSentmentList}

		return ctx.Status(200).JSON(response)
	}
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

		var positiveTweets []twitter.TweetSentment
		var negativeTweets []twitter.TweetSentment

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

		response := dto.AverageTweetsMoodResponseBody{
			Positive: dto.AverageMood{
				Total:      positiveTweetsCount,
				Percentage: float32(positiveTweetsCount) / float32(totalTweetsCount),
			},
			Negative: dto.AverageMood{
				Total:      negativeTweetsCount,
				Percentage: float32(negativeTweetsCount) / float32(totalTweetsCount),
			},
			Total: totalTweetsCount,
		}

		return ctx.Status(200).JSON(response)
	}
}
