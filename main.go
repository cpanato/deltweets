package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

func main() {

	apiKey := os.Getenv("TWITTER_API_KEY")
	if apiKey == "" {
		log.Fatal("Twitter API key required")
		os.Exit(2)
	}

	apiSecret := os.Getenv("TWITTER_API_SECRET")
	if apiSecret == "" {
		log.Fatal("Twitter API secret required")
		os.Exit(2)
	}

	accessToken := os.Getenv("TWITTER_ACCESS_TOKEN")
	if accessToken == "" {
		log.Fatal("Twitter access token required")
		os.Exit(2)
	}

	accessTokenSecret := os.Getenv("TWITTER_ACCESS_TOKEN_SECRET")
	if accessTokenSecret == "" {
		log.Fatal("Twitter access token secret required")
		os.Exit(2)
	}

	twitterUserName := os.Getenv("TWITTER_USERNAME")
	if twitterUserName == "" {
		log.Fatal("Twitter username required")
		os.Exit(2)
	}

	config := oauth1.NewConfig(apiKey, apiSecret)
	token := oauth1.NewToken(accessToken, accessTokenSecret)

	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)

	timeToDelete := time.Now().AddDate(0, 0, -3).Format("2006-01-02")
	fmt.Printf("will delete tweets older then %+v\n", timeToDelete)
	query := fmt.Sprintf("from:%s", twitterUserName)

	// search tweets
	searchTweetParams := &twitter.SearchTweetParams{
		Query: query,
		Until: timeToDelete,
		Count: 500,
	}
	search, _, _ := client.Search.Tweets(searchTweetParams)
	fmt.Printf("TWEETS:\n%+v\n", len(search.Statuses))

	for _, status := range search.Statuses {
		fmt.Printf("Will delete: %+v - %+v", status.ID, status)
		// status destroy
		params := &twitter.StatusDestroyParams{TrimUser: twitter.Bool(false)}
		tweet, resp, err := client.Statuses.Destroy(status.ID, params)
		fmt.Printf("STATUSES DESTROY:\n%+v\n", tweet)
		fmt.Printf("RESP:\n%+v\n", resp)
		fmt.Printf("Err:\n%+v\n", err)
	}

	os.Exit(0)
}
