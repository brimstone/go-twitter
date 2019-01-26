package twitter

import (
	"bufio"
	"fmt"
	"os"

	dt "github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

type Tokens struct {
	ConsumerKey    string
	ConsumerSecret string
	AccessToken    string
	AccessSecret   string
}

type AuthResponse struct {
	DidAuth      bool
	AccessToken  string
	AccessSecret string
}

func NewClient(t Tokens) (*dt.Client, *AuthResponse, error) {
	// Use the consumer key and secret
	config := oauth1.NewConfig(t.ConsumerKey, t.ConsumerSecret)
	config.CallbackURL = "oob"
	config.Endpoint = oauth1.Endpoint{
		RequestTokenURL: "https://api.twitter.com/oauth/request_token",
		AuthorizeURL:    "https://api.twitter.com/oauth/authorize",
		AccessTokenURL:  "https://api.twitter.com/oauth/access_token",
	}
	if t.AccessToken == "" {
		t.AccessToken = os.Getenv("TWITTER_TOKEN")
		t.AccessSecret = os.Getenv("TWITTER_SECRET")
	}
	resp := AuthResponse{}
	if t.AccessToken == "" {
		// Get a request token from twitter
		requestToken, requestSecret, err := config.RequestToken()
		// Let the user know how to get a verification code
		fmt.Printf("Now visit https://api.twitter.com/oauth/authorize?oauth_token=%s\n", requestToken)

		// Get verification code from the user
		fmt.Printf("Code: ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		vcode := scanner.Text()

		// Trade in the verification code for a proper access token
		t.AccessToken, t.AccessSecret, err = config.AccessToken(requestToken, requestSecret, vcode)
		if err != nil {
			return nil, nil, err
		}
		resp.DidAuth = true
		resp.AccessToken = t.AccessToken
		resp.AccessSecret = t.AccessSecret
	}

	token := oauth1.NewToken(t.AccessToken, t.AccessSecret)
	// http.Client will automatically authorize Requests
	httpClient := config.Client(oauth1.NoContext, token)

	// twitter client
	client := dt.NewClient(httpClient)
	return client, &resp, nil
}
