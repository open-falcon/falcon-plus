package portal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var HostnamesUrl = "http://127.0.0.1:5050/api/group/%s/hosts.json"

type hostnamesDto struct {
	Msg  string   `json:"msg"`
	Data []string `json:"data"`
}

func Hostnames(groupName string, hostnamesUrl ...string) ([]string, error) {
	pattern := HostnamesUrl
	if len(hostnamesUrl) > 0 {
		pattern = hostnamesUrl[0]
	}

	url := fmt.Sprintf(pattern, groupName)

	resp, err := http.Get(url)
	if err != nil {
		log.Println("[E]", err)
		return []string{}, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("[E]", err)
		return []string{}, err
	}

	if resp.StatusCode != 200 {
		log.Printf("[E] status code: %d != 200, response: %s", resp.StatusCode, string(body))
		return []string{}, fmt.Errorf(string(body))
	}

	var dto hostnamesDto
	err = json.Unmarshal(body, &dto)
	if err != nil {
		log.Println("[E]", err)
		return []string{}, err
	}

	return dto.Data, nil
}
