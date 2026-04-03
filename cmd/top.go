package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"ekhoes-ctl/gptable"
)

type ProcInfo struct {
    PID  int32
    User string
    CPU  float64
    Name string
}

func TopCpuProcesses(args []string) error {
	endpoint := GetCtlEndpoint("top")
	
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

	var items []ProcInfo
	if err := json.Unmarshal(bodyBytes, &items); err != nil {
		panic(err)
	}
	
	//fmt.Printf("%v\n", info);
	
	gptable.Init();
	gptable.SetHeader(
		"PID",
		"User",
		"CPU (%)",
		"Process",
	)
	
	for _, item := range items {
		pidStr := fmt.Sprintf("%d", item.PID);
		cpuStr := fmt.Sprintf("%.2f", item.CPU);

		gptable.AppendRow(
			pidStr,
			item.User,
			cpuStr,
			item.Name,
		)
	}
	
	gptable.Render()

	return nil
}
