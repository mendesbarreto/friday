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
	Meta Meta   `json:"meta"`
}
type Tweet struct {
	ID                  string    `json:"id"`
	CreatedAt           time.Time `json:"created_at"`
	Text                string    `json:"text"`
}
type Meta struct {
	ResultCount int    `json:"result_count"`
	NextToken   string `json:"next_token"`
}


type TweetSentment struct {
    Tweet Tweet `json:"tweet"` 
    Score int8 `json:"score"`
}

type TweetSentmentResponse struct {
    Data []TweetSentment `json:"data"`
}

const twitterUrl string = "https://api.twitter.com/2/lists/1530602256525041665/tweets"
var twitterToken string = os.Getenv("TWITTER_API_TOKEN")

func getTweetSentment(message string) uint8 {
    sentiment, err := GetTextSentiment(message)

    if err != nil {
        fmt.Println("Twitter was not able to score")
        return 3
    }

    return sentiment
}

func GetTweetsFromToday() fiber.Handler {
    return func (ctx *fiber.Ctx) error {

        listId := ctx.Params("id")

        if len(listId) == 0 {
            return dto.BadRequest(ctx, "The url should contain the string")
        }
 
        client := http.Client{}

        request, err := http.NewRequest("GET", twitterUrl, nil)

        if err != nil {
            return dto.InternalServerError(ctx, "The request could not be created")
        }
        
        twitterBearer := fmt.Sprintf("Bearer %s", twitterToken)

        request.Header = http.Header{
            "Authorization": { twitterBearer },
        }

        res, err := client.Do(request)

        if err != nil {
            return dto.InternalServerError(ctx, err.Error())
        }
        
        defer res.Body.Close()


        var result TweetsRespose
        resBody, err := ioutil.ReadAll(res.Body)

        if err := json.Unmarshal(resBody, &result); err != nil {
            return dto.InternalServerError(ctx, err.Error())
        }

        response := TweetSentmentResponse{
            Data: make([]TweetSentment, len(result.Data)), 
        }

        var wg sync.WaitGroup
        for index, data := range result.Data {
            wg.Add(1)
            go func(index int, tweet Tweet) {
                response.Data[index] = TweetSentment{
                    Tweet: tweet,
                    Score: int8(getTweetSentment(tweet.Text)),
                }
                wg.Done()
            }(index, data)
        }

        wg.Wait()
        
        return ctx.Status(200).JSON(response)
    }
}
