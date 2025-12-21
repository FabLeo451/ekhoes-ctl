package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"

	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/term"
)

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Token string `json:"token"`
}

var _token string

func getToken() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	tokenPath := filepath.Join(home, ".ekhoes", "token")

	info, err := os.Stat(tokenPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}

	if !info.Mode().IsRegular() {
		return "", errors.New("not a regular file")
	}

	data, err := os.ReadFile(tokenPath)
	if err != nil {
		return "", err
	}

	_token = string(data)

	return string(_token), nil
}

func saveToken(token string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	filePath := filepath.Join(home, ".ekhoes/token")

	// 0600 → solo l'utente può leggere/scrivere
	return os.WriteFile(filePath, []byte(token), 0600)
}

func login(args []string) error {
	var creds Credentials
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Email: ")
	username, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	creds.Email = strings.TrimSpace(username)

	fmt.Print("Password: ")
	bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println() // newline dopo input password
	if err != nil {
		return err
	}

	creds.Password = strings.TrimSpace(string(bytePassword))

	if creds.Email == "" || creds.Password == "" {
		return errors.New("empty credentials")
	}

	// Convert credentials to json
	jsonData, err := json.Marshal(creds)
	if err != nil {
		return err
	}

	// Init POST request

	params := url.Values{}
	params.Set("nosession", "1")

	endpoint := fmt.Sprintf("%s/login?%s", conf.URL, params.Encode())

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	// Create client and call
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	//fmt.Println("Status:", resp.Status)

	if resp.StatusCode == 200 {
		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)

		tokenInterface, ok := result["token"]
		if !ok {
			log.Fatal("Key 'token' nmissing")
		}

		token, ok := tokenInterface.(string)
		if !ok {
			log.Fatal("Value 'token' is not a string")
		}

		// Save the token
		err = saveToken(token)
		if err != nil {
			return err
		}

		fmt.Printf("\nHello, %s. You have successfully logged in!\n\n", result["name"])

	} else {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return errors.New("Unable to read body")
		}

		return errors.New(string(bodyBytes))
	}

	return nil
}
