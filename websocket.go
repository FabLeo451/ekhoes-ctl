package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	"net/http"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
)

type ConnectionItem struct {
	SessionId string    `json:"sessionId"`
	Email     string    `json:"email"`
	Created   time.Time `json:"created"`
}

func getWebsocketConnections(args []string) error {
	endpoint := fmt.Sprintf("%s/connections", conf.URL)

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", _token)

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

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.New("Unable to read body")
	}

	//fmt.Println(string(bodyBytes))

	var items []ConnectionItem
	if err := json.Unmarshal(bodyBytes, &items); err != nil {
		panic(err)
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	// stile compatto tipo kubectl / ps
	t.SetStyle(table.StyleLight)
	t.Style().Options.DrawBorder = false
	t.Style().Options.SeparateColumns = false
	t.Style().Options.SeparateHeader = false
	t.Style().Options.SeparateRows = false

	t.AppendHeader(table.Row{
		"Session Id",
		"User Email",
		"Created",
	})

	for _, item := range items {

		t.AppendRow(table.Row{
			item.SessionId,
			item.Email,
			item.Created.Format("2006-01-02 15:04:05"),
		})
	}

	//fmt.Println()
	t.Render()

	return nil
}
