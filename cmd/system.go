package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	//"time"

	"net/http"
	//"os"

	//"github.com/jedib0t/go-pretty/v6/table"
	"ekhoes-ctl/config"
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

func HumanBytes(b uint64) string {
    const (
        KB = 1 << 10
        MB = 1 << 20
        GB = 1 << 30
        TB = 1 << 40
    )

    format := func(v float64, unit string) string {
        if v == float64(uint64(v)) {
            return fmt.Sprintf("%.0f%s", v, unit)
        }
        return fmt.Sprintf("%.2f%s", v, unit)
    }

    switch {
    case b >= TB:
        return format(float64(b)/TB, "TB")
    case b >= GB:
        return format(float64(b)/GB, "GB")
    case b >= MB:
        return format(float64(b)/MB, "MB")
    case b >= KB:
        return format(float64(b)/KB, "KB")
    default:
        return fmt.Sprintf("%dB", b)
    }
}

func GetSystemInfo(args []string) error {
	endpoint := fmt.Sprintf("%s/ctl/system", config.Conf.URL)
	
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
