package cmd

import (
	"fmt"
	"ekhoes-ctl/config"
)

func GetCtlEndpoint(method string) string {
	endpoint := fmt.Sprintf("%s%s/ctl/%s", config.Conf.URL, config.Conf.RootPath, method)
	
	return endpoint
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

