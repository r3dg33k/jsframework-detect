package detector

import (
	"encoding/json"
	"net/http"
)

func DetectBuiltWith(apiKey, url string) ([]string, error) {
	req, _ := http.NewRequest(
		"GET",
		"https://api.builtwith.com/v20/api.json?KEY="+apiKey+"&LOOKUP="+url,
		nil,
	)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data struct {
		Results []struct {
			Result struct {
				Paths []struct {
					Technologies []struct {
						Name string `json:"Name"`
					} `json:"Technologies"`
				} `json:"Paths"`
			} `json:"Result"`
		} `json:"Results"`
	}

	json.NewDecoder(resp.Body).Decode(&data)

	var out []string
	for _, r := range data.Results {
		for _, p := range r.Result.Paths {
			for _, t := range p.Technologies {
				out = append(out, t.Name)
			}
		}
	}

	return out, nil
}