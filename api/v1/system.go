package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type ConfigResult struct {
	Value string `json:"value"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	url := os.Getenv("SUPABASE_URL")
	key := os.Getenv("SUPABASE_SERVICE_ROLE_KEY")

	queryType := r.URL.Query().Get("get")
	if queryType == "" {
		queryType = "version"
	}

	reqURL := fmt.Sprintf("%s/rest/v1/system_config?key=eq.%s&select=value", url, queryType)
	
	req, _ := http.NewRequest("GET", reqURL, nil)
	req.Header.Set("apikey", key)
	req.Header.Set("Authorization", "Bearer "+key)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprint(w, "Error")
		return
	}
	defer resp.Body.Close()

	var results []ConfigResult
	json.NewDecoder(resp.Body).Decode(&results)

	if len(results) > 0 {
		fmt.Fprint(w, results[0].Value)
	} else {
		fmt.Fprint(w, "unknown")
	}
}