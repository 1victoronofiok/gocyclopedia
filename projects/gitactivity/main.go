package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

type GitActivityResp struct {
	Actor struct {
		DisplayLogin string `json:"display_login"`
	} `json:"actor"`
}
func fetchUserActivity(u string) error {
	url := fmt.Sprintf("https://api.github.com/users/%s/events", u)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New("failed request")
	}

	var activities []GitActivityResp
	err = json.NewDecoder(resp.Body).Decode(&activities)

	if err != nil {
		return err
	}

	fmt.Println(activities)

	return nil
}


func main() {
	var username string
	fmt.Print("Enter github username: ")
	_, err := fmt.Scanf("%s", &username)

	if err != nil {
		log.Fatal(err)
	}

	if username == "" {
		log.Fatal("username response empty")
	}

	fetchUserActivity(username)
}

//1victoronofiok