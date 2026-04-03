package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"ekhoes-ctl/gptable"
)

type SystemInfo struct {
    Hostname     string  `json:"hostname"`
    OS           string  `json:"os"`
    Platform     string  `json:"platform"`
    PlatformVer  string  `json:"platform_version"`
    KernelVersion       string  `json:"kernel_version"`
    CPULoad      float64 `json:"cpu_load"`
    RAMUsed      uint64  `json:"ram_used"`
    RAMTotal     uint64  `json:"ram_total"`
    DiskUsed     uint64  `json:"disk_used"`
    DiskTotal    uint64  `json:"disk_total"`
}

func GetSystemInfo(args []string) error {
	endpoint := GetCtlEndpoint("system")
	
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

	var info SystemInfo
	if err := json.Unmarshal(bodyBytes, &info); err != nil {
		panic(err)
	}
	
	//fmt.Printf("%v\n", info);

	fmt.Println();
	fmt.Printf(" Hostname        : %s\n", info.Hostname);
	fmt.Printf(" Platform        : %s\n", info.OS);
	fmt.Printf(" Operating System: %s %s %s\n\n", info.Platform, info.PlatformVer, info.KernelVersion);
	fmt.Printf(" CPU Load        : %.2f%%\n\n", info.CPULoad)

	pctRAM := fmt.Sprintf("%.2f%%", float64(info.RAMUsed)/float64(info.RAMTotal)*100)
	pctDisk := fmt.Sprintf("%.2f%%", float64(info.DiskUsed)/float64(info.DiskTotal)*100)

	gptable.Init();
	gptable.SetHeader(
		"Resource",
		"Total",
		"Used",
		"Usage percentage",
	)
	
	gptable.AppendRow(
		"Memory",
		HumanBytes(info.RAMTotal),
		HumanBytes(info.RAMUsed),
		pctRAM,
	)
	
	gptable.AppendRow(
		"Disk",
		HumanBytes(info.DiskTotal),
		HumanBytes(info.DiskUsed),
		pctDisk,
	)
	
	gptable.Render()
	
	fmt.Println();

	return nil
}
