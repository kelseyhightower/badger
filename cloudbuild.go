package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/kelseyhightower/run"
)

const (
	userAgent = "badger/0.0.1"
)

const cloudbuildEndpoint = "https://cloudbuild.googleapis.com/v1"
var scopes = []string{"https://www.googleapis.com/auth/cloud-platform"}

type BuildList struct {
	Builds []Build `json:"builds"`
}

type Build struct {
	Status string `json:"status"`
}

func getBuildStatus(project, id string) (string, error) {
	endpoint := fmt.Sprintf("%s/projects/%s/builds", cloudbuildEndpoint, project)

	u, err := url.Parse(endpoint)
	if err != nil {
		return "", err
	}

	q := u.Query()
	q.Set("filter", fmt.Sprintf("trigger_id=\"%s\"", id))
	q.Set("pageSize", "1")
	u.RawQuery = q.Encode()

	request, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return "", err
	}

	token, err := run.Token(scopes)
	if err != nil {
		return "", err
	}

	request.Header.Set("User-Agent", userAgent)
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

	timeout := time.Duration(5) * time.Second
	client := http.Client{Timeout: timeout}

	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	switch s := response.StatusCode; s {
	case 200:
		break
	case 401:
		return "", errors.New("cloud build api unauthorized")
	case 403:
		return "", errors.New("cloud build api permission denied")
	case 404:
		return "", errors.New("build not found")
	default:
		return "", errors.New("error retrieving build status")
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("error reading build response body", err)
	}

	var buildList BuildList
	err = json.Unmarshal(data, &buildList)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling build status: %w", err)
	}

	return buildList.Builds[0].Status, nil
}
