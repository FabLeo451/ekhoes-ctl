package cmd

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

func KillSession(args []string) error {

	if len(args) < 2 {
		return errors.New("Missing session id")
	}

	endpoint := GetCtlEndpoint("session")
	endpoint = fmt.Sprintf("%s/%s", endpoint, args[1])
	token, _ := GetToken()

	req, err := http.NewRequest("DELETE", endpoint, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	// Create client and call
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	//fmt.Println("Status:", resp.Status)

	if resp.StatusCode != 200 {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return errors.New("Unable to read body")
		}

		return errors.New(string(bodyBytes))
	}

	return nil
}
