package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const apiEndpoint = "https://api.zerotier.com/api/v1/network"

type network struct {
	Config struct {
		Name              string `json:"name"`
		Private           bool   `json:"private"`
		IpAssignmentPools []struct {
			IpRangeStart string `json:"ipRangeStart"`
			IpRangeEnd   string `json:"ipRangeEnd"`
		} `json:"ipAssignmentPools"`
		Routes []struct {
			Target string `json:"target"`
		}
	} `json:"config"`
}

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

	deleteNetwork(trimmedToken, "9bee8941b5bc0d24")
}

func getNetwork(token string) {
	// Make API request
	req, err := http.NewRequest("GET", apiEndpoint, nil)
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

func createNetwork(token string) {
	// Create network
	net := &network{
		Config: struct {
			Name              string `json:"name"`
			Private           bool   `json:"private"`
			IpAssignmentPools []struct {
				IpRangeStart string `json:"ipRangeStart"`
				IpRangeEnd   string `json:"ipRangeEnd"`
			} `json:"ipAssignmentPools"`
			Routes []struct {
				Target string `json:"target"`
			}
		}{
			Name:    "MyNetwork",
			Private: true,
			IpAssignmentPools: []struct {
				IpRangeStart string `json:"ipRangeStart"`
				IpRangeEnd   string `json:"ipRangeEnd"`
			}{
				{
					IpRangeStart: "172.25.0.1",
					IpRangeEnd:   "172.25.255.254",
				},
			},
			Routes: []struct {
				Target string `json:"target"`
			}{
				{
					Target: "172.25.0.0/16",
				},
			},
		},
	}

	body, err := json.Marshal(net)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Make API request
	req, err := http.NewRequest("POST", apiEndpoint, bytes.NewBuffer(body))
	if err != nil {
		fmt.Println(err)
		return
	}

	req.Header.Add("Authorization", "token "+token)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()
	fmt.Println(string(body))
}

func deleteNetwork(token string, networkId string) {
	// Make API request
	req, err := http.NewRequest("DELETE", apiEndpoint+"/"+networkId, nil)
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
}

