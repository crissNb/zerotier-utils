package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Failed to get user home directory:", err)
		return
	}

	// Read token from a file
	token, err := ioutil.ReadFile(filepath.Join(homeDir, "/.config/zerotier-utils/token"))
	if err != nil {
		fmt.Println("Failed to read token from file:", err)
		return
	}

	trimmedToken := strings.TrimSpace(string(token))

	getNetwork(trimmedToken)
}

func getNetwork(token string) {
	// Make API request
	req, err := http.NewRequest("GET", "https://api.zerotier.com/api/v1/network", nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	req.Header.Add("Authorization", "token "+token)

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(body))
}
