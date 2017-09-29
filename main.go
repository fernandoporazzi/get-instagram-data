package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var usernames = []string{"fernandoporazzi", "danielrohersphotos", "elisabeteporazzi"}

// User represents the user object in response
type User struct {
	FullName   string         `json:"full_name"`
	FollowedBy map[string]int `json:"followed_by"`
	Follows    map[string]int `json:"follows"`
}

// Response represents the whole response
type Response struct {
	User User `json:"user"`
}

func asyncHTTPGets(usernames []string) []*Response {
	ch := make(chan *Response, len(usernames))

	responses := []*Response{}

	for _, u := range usernames {
		go func(un string) {
			res, err := http.Get("https://www.instagram.com/" + un + "/?__a=1")

			if err != nil {
				log.Fatal(err)
			}

			data, err := ioutil.ReadAll(res.Body)

			res.Body.Close()

			if err != nil {
				log.Fatal(err)
			}

			response := Response{}

			json.Unmarshal(data, &response)

			ch <- &response
		}(u)
	}

	for {
		select {
		case r := <-ch:
			responses = append(responses, r)
			if len(responses) == len(usernames) {
				return responses
			}
		case <-time.After(50 * time.Millisecond):
			fmt.Printf(".")
		}
	}
}

func main() {
	results := asyncHTTPGets(usernames)

	fmt.Println()
	fmt.Println()

	for _, result := range results {
		fmt.Printf("Name: %s\n", result.User.FullName)
		fmt.Printf("Followed by: %d\n", result.User.FollowedBy["count"])
		fmt.Printf("Follows: %d\n", result.User.Follows["count"])
		fmt.Println()
	}
}
