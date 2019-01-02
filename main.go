package main

import (
	"fmt"
	"log"
	"os"
	"strings"
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

	tweetsToIgnore := strings.Split(os.Getenv("TWEETS_IGNORE"), ",")

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
	fmt.Printf("TOTAL TWEETS:\n%+v\n", len(search.Statuses))

	for _, status := range search.Statuses {
		flag := false
		for _, tweet := range tweetsToIgnore {
			if status.IDStr == tweet {
				fmt.Printf("tweet is in the whitelist - %v\n", status.ID)
				flag = true
				break
			}
		}
		if flag {
			continue
		}
		fmt.Printf("%v\n", status.ID)
		fmt.Printf("Will delete: %+v - %+v\n", status.ID, status)
		// status destroy
		params := &twitter.StatusDestroyParams{TrimUser: twitter.Bool(false)}
		tweet, resp, err := client.Statuses.Destroy(status.ID, params)
		fmt.Printf("STATUSES DESTROY:\n%+v\n", tweet)
		fmt.Printf("RESP:\n%+v\n", resp)
		fmt.Printf("Err:\n%+v\n", err)
	}

	search, _, _ = client.Search.Tweets(searchTweetParams)
	fmt.Printf("TOTAL TWEETS AFTER DELETION:\n%+v\n", len(search.Statuses))

	os.Exit(0)
}
