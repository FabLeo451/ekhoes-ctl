package cmd

import (
	"encoding/json"
	"errors"
	"io"
	"time"
	"net/http"

	"ekhoes-ctl/gptable"
)

type ConnectionItem struct {
	SessionId 	 string    `json:"sessionId"`
	Email     	 string    `json:"email"`
	Created   	 time.Time `json:"created"`
	LastActivity     string    `json:"lastActivity"`
	LastActivityTime time.Time `json:"lastActivityTime"`
}

func GetWebsocketConnections(args []string) error {
	endpoint := GetCtlEndpoint("ws")
	
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

	var items []ConnectionItem
	if err := json.Unmarshal(bodyBytes, &items); err != nil {
		panic(err)
	}
	
	gptable.Init();
	gptable.SetHeader(
		"Session Id",
		"User",
		"Created",
		"Last activity",
		"Last activity time",
	)

	for _, item := range items {

		tUTC := item.Created.UTC()
		tLocal := tUTC.In(time.Local)
		created := tLocal.Format("2006-01-02 15:04:05")

		tUTC = item.LastActivityTime.UTC()
		tLocal = tUTC.In(time.Local)
		lastActivityTime := tLocal.Format("2006-01-02 15:04:05")		
		
		gptable.AppendRow(
			item.SessionId,
			item.Email,
			created,
			item.LastActivity,
			lastActivityTime,
		)
	}

	gptable.Render()

	return nil
}

/*
func __GetWebsocketConnections(args []string) error {
	endpoint := fmt.Sprintf("%s/connections", config.Conf.URL)
	
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

	var items []ConnectionItem
	if err := json.Unmarshal(bodyBytes, &items); err != nil {
		panic(err)
	}
	
	gptable.Init();
	gptable.SetHeader(
		"Session Id",
		"User",
		"Created",
		"Last activity",
		"Last activity time",
	)
	
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
		"User",
		"Created",
		"Last activity",
		"Last activity time",
	})

	for _, item := range items {

		tUTC := item.Created.UTC()
		tLocal := tUTC.In(time.Local)
		created := tLocal.Format("2006-01-02 15:04:05")

		tUTC = item.LastActivityTime.UTC()
		tLocal = tUTC.In(time.Local)
		lastActivityTime := tLocal.Format("2006-01-02 15:04:05")		

		t.AppendRow(table.Row{
			item.SessionId,
			item.Email,
			created,
			item.LastActivity,
			lastActivityTime,
		})
		
		gptable.AppendRow(
			item.SessionId,
			item.Email,
			created,
			item.LastActivity,
			lastActivityTime,
		)
	}

	//fmt.Println()
	t.Render()
	gptable.Render()

	return nil
}
*/
