package config

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type Configuration struct {
	Port     string `json:"port"`
	Host     string `json:"host"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

func GetConfig() *Configuration {
	file, err := ioutil.ReadFile("config.json")
	if err != nil {
		fmt.Println(err)
	}

	var config Configuration
	json.Unmarshal(file, &config)

	return &config

}

func GetEnvKey(keyname string) string {
	content, err := ioutil.ReadFile(".env")
	if err != nil {
		fmt.Println("Error reading .env file")
	}
	scanner := bufio.NewScanner(strings.NewReader(string(content)))
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}
		if line[0] == '#' {
			continue
		}
		keyVal := strings.SplitN(line, "=", 2)
		if len(keyVal) != 2 {
			continue
		}
		key := keyVal[0]
		val := keyVal[1]
		os.Setenv(key, val)
	}
	apiKey := os.Getenv(keyname)
	fmt.Println("API Key:", apiKey)

	return apiKey
}
