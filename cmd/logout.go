package cmd

import (
	"fmt"
	"os"
	"path/filepath"
)

func Logout(args []string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	
	filePath := filepath.Join(home, ".ekhoes/token")
	
	err = os.Remove(filePath)
	if err != nil {
		return err
	}
	
	fmt.Println("Authentication token deleted")
	
	return nil
}
