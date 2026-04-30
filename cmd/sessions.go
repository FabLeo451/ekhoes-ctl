package cmd

import (
	"encoding/json"
	"errors"
	"io"
	"time"
	"net/http"
	"fmt"

	"ekhoes-ctl/gptable"
)

type Item struct {
	Id         string `json:"id"`
	Status     string `json:"status"`
	User       User   `json:"user"`
	Agent      string `json:"agent"`
	Platform   string `json:"platform"`
	DeviceType string `json:"deviceType"`
	Updated    string `json:"updated"`
	TTL        time.Duration `json:"ttl"`
}

func HumanizeDuration(d time.Duration) string {
	days := int(d.Hours()) / 24
	h := int(d.Hours()) % 24
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60

	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm %ds", days, h, m, s)
	}
	if h > 0 {
		return fmt.Sprintf("%dh %dm %ds", h, m, s)
	}
	if m > 0 {
		return fmt.Sprintf("%dm %ds", m, s)
	}
	return fmt.Sprintf("%ds", s)
}

func GetSessions(args []string) error {
	endpoint := GetCtlEndpoint("sessions")
	token, _ := GetToken()

	req, err := http.NewRequest("GET", endpoint, nil)
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

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.New("Unable to read body")
	}

	fmt.Println(string(bodyBytes))

	var items []Item
	if err := json.Unmarshal(bodyBytes, &items); err != nil {
		panic(err)
	}

	gptable.Init();
	gptable.SetHeader(
		"Id",
		"Status",
		"Name",
		"Email",
		"Agent",
		"Platform",
		"Device Type",
		"Updated",
		"TTL",
	)

	for _, item := range items {

		// parse RFC3339
		tm, err := time.Parse(time.RFC3339, item.Updated)
		if err != nil {
			panic(err)
		}

		updatedLocal := tm.In(time.Local)

		gptable.AppendRow(
			item.Id,
			item.Status,
			item.User.Name,
			item.User.Email,
			item.Agent,
			item.Platform,
			item.DeviceType,
			updatedLocal.Format("2006-01-02 15:04:05"),
			HumanizeDuration(item.TTL),
		)
	}

	gptable.Render()

	return nil
}

/*
func GetSessions(args []string) error {
	endpoint := fmt.Sprintf("%s/sessions", config.Conf.URL)
	token, _ := GetToken()

	req, err := http.NewRequest("GET", endpoint, nil)
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

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.New("Unable to read body")
	}

	//fmt.Println(string(bodyBytes))

	var items []Item
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
		"Id",
		"Status",
		"Name",
		"Email",
		"Agent",
		"Platform",
		"Device Type",
		"Updated",
	})

	for _, item := range items {

		// parse RFC3339
		tm, err := time.Parse(time.RFC3339, item.Updated)
		if err != nil {
			panic(err)
		}

		updatedLocal := tm.In(time.Local)

		t.AppendRow(table.Row{
			item.Id,
			item.Status,
			item.User.Name,
			item.User.Email,
			item.Agent,
			item.Platform,
			item.DeviceType,
			updatedLocal.Format("2006-01-02 15:04:05"),
		})
	}

	//fmt.Println()
	t.Render()

	return nil
}
*/
