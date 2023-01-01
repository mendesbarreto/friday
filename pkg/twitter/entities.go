package twitter

import "time"

type Tweet struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
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
