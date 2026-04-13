package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type StatusResponse struct {
	Status string `json:"status"`
	IsWeb  bool   `json:"isweb"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	supabaseStorageURL := os.Getenv("SUPABASE_URL") + "/storage/v1/object/public/versions/"

	v := r.URL.Query().Get("v")
	if v == "" {
		v = "current"
	}

	// Расширение по умолчанию
	ext := ".pptm"

	// Если запрашивается current, проверяем статус
	if v == "current" {
		client := http.Client{Timeout: 5 * time.Second}
		// ОБЯЗАТЕЛЬНО СО СЛЕШЕМ В КОНЦЕ
		resp, err := client.Get("https://nestoros.vercel.app/settings/status/")
		
		if err == nil && resp.StatusCode == 200 {
			var status StatusResponse
			if err := json.NewDecoder(resp.Body).Decode(&status); err == nil {
				if status.IsWeb {
					ext = ".zip"
				} else {
					ext = ".pptm"
				}
			}
			resp.Body.Close()
		}
	} else {
		// Карта статичных расширений для старых версий
		extensions := map[string]string{
			"5.0b":  ".pptx",
			"5.0r":  ".pptm",
			"5.0u1": ".pptm",
			"5.0m":  ".pptm",
			"5.0t":  ".pptm",
			"5.0tv": ".pptm",
			"6.0":   ".pptx",
			"7.0":   ".pptm",
			"8.0":   ".pptm",
			"8.1":   ".pptm",
			"8.2":   ".pptm",
		}
		
		if val, ok := extensions[v]; ok {
			ext = val
		} else {
			http.Error(w, "Version not found", 404)
			return
		}
	}

	finalURL := fmt.Sprintf("%s%s%s", supabaseStorageURL, v, ext)
	http.Redirect(w, r, finalURL, http.StatusFound)
}