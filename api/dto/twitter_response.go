package dto

import "github.com/mendesbarreto/friday/pkg/twitter"

type TweetsResposeBody struct {
	Data []twitter.Tweet `json:"data"`
	Meta twitter.Meta    `json:"meta"`
}

type TweetSentmentResponseBody struct {
	Data []twitter.TweetSentment `json:"data"`
}

type AverageMood struct {
	Total      int32   `json:"total"`
	Percentage float32 `json:"percentage"`
}

type AverageTweetsMoodResponseBody struct {
	Positive AverageMood `json:"positive"`
	Negative AverageMood `json:"negative"`
	Total    int32       `json:"total"`
}
