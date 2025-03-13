package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
)

const (
	redirectURL    = "http://localhost:8080/callback"
	stravaAuthURL  = "https://www.strava.com/oauth/authorize"
	stravaTokenURL = "https://www.strava.com/oauth/token"
)

var (
	clientID     = os.Getenv("STRAVA_CLIENT_ID")
	clientSecret = os.Getenv("STRAVA_CLIENT_SECRET")
)

var oauthConfig = &oauth2.Config{
	ClientID:     clientID,
	ClientSecret: clientSecret,
	RedirectURL:  redirectURL,
	Endpoint: oauth2.Endpoint{
		AuthURL:  stravaAuthURL,
		TokenURL: stravaTokenURL,
	},
}

type Activity struct {
	ID                 int64   `json:"id"`
	Name               string  `json:"name"`
	Distance           float64 `json:"distance"`
	MovingTime         int     `json:"moving_time"`
	ElapsedTime        int     `json:"elapsed_time"`
	TotalElevationGain float64 `json:"total_elevation_gain"`
	Type               string  `json:"type"`
	StartDate          string  `json:"start_date"`
}

func main() {
	http.HandleFunc("/", handleMain)
	http.HandleFunc("/callback", handleCallback)
	http.HandleFunc("/activities", handleActivities)
	log.Println("Server started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleMain(w http.ResponseWriter, r *http.Request) {
	url := oauthConfig.AuthCodeURL("state")
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	token, err := oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		log.Fatal(err)
	}

	// Store the token or use it to fetch data
	http.Redirect(w, r, "/activities?access_token="+token.AccessToken, http.StatusTemporaryRedirect)
}

func handleActivities(w http.ResponseWriter, r *http.Request) {
	accessToken := r.URL.Query().Get("access_token")
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://www.strava.com/api/v3/athlete/activities", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	fmt.Println("accessToken: ", accessToken)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var activities []Activity
	if err := json.NewDecoder(resp.Body).Decode(&activities); err != nil {
		log.Fatal(err)
	}

	for _, activity := range activities {
		fmt.Fprintf(w, "Activity: %s, Distance: %.2f km, D+: %.2f m \n", activity.Name, activity.Distance/1000, activity.TotalElevationGain)
	}
}
