package detector

import (
	"encoding/json"
	"net/http"
)

func DetectWappalyzer(apiKey, url string) ([]string, error) {
	req, _ := http.NewRequest(
		"GET",
		"https://api.wappalyzer.com/v2/lookup/?url="+url,
		nil,
	)

	req.Header.Set("x-api-key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data struct {
		Technologies []struct {
			Name string `json:"name"`
		} `json:"technologies"`
	}

	json.NewDecoder(resp.Body).Decode(&data)

	var out []string
	for _, t := range data.Technologies {
		out = append(out, t.Name)
	}

	return out, nil
}