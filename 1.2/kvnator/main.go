package main

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type kv struct {
	CreateIndex int
	ModifyIndex int
	LockIndex   int
	Flags       int
	Key         string
	Value       string
}

func main() {
	envFlag := flag.String("env", "STACKENGINE_IP", "The key of the env variable you are wanting to write")
	flag.Parse()

	apiKey := os.Getenv("STACKENGINE_API_TOKEN")
	leaderIP := os.Getenv("STACKENGINE_LEADER_IP")
	discoveryKey := os.Getenv("STACKENGINE_SERVICE_DISCOVERY_KEY")

	if apiKey == "" {
		fmt.Println("STACKENGINE_API_TOKEN was not available in the ENV")
		os.Exit(1)
	}
	if leaderIP == "" {
		fmt.Println("STACKENGINE_LEADER_IP was not available in the ENV")
		os.Exit(1)
	}
	if discoveryKey == "" {
		fmt.Println("STACKENGINE_SERVICE_DISCOVERY_KEY was not available in the ENV")
		os.Exit(1)
	}

	url := "https://" + leaderIP + ":8443/api/kv/" + discoveryKey + "?recurse"

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Errorf("Error while calling %s", url)
		panic(err)
	}
	defer resp.Body.Close()

	//fmt.Println("response Status:", resp.Status)
	//fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body[:]))
	var response []kv
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("Error while unmarshaling json body.")
		fmt.Println("kvnator did not set a value. Returning.")
		return
	}

	if len(response) > 0 {
		data, err := base64.StdEncoding.DecodeString(response[0].Value)
		if err == nil {
			var stringToWrite = "export " + *envFlag + "=" + string(data[:]) + "\n"
			err = ioutil.WriteFile("/tmp/kvnator.txt", []byte(stringToWrite), 0664)
			if err != nil {
				fmt.Println("Could not write to file")
				panic(err)
				return
			}

			err = os.Setenv(*envFlag, string(data[:]))
			if err != nil {
				panic(err)
				return
			}
			fmt.Println("export " + *envFlag + "=" + os.Getenv(*envFlag))
		} else {
			fmt.Println("kvnator did not set a value. Returning.")
		}
	} else {
		fmt.Println("kvnator did not set a value. Returning.")
	}
}
