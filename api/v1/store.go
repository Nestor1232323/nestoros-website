package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type App struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Version     string  `json:"version"`
	Description string  `json:"description"`
	DownloadURL string  `json:"download_url"`
	Price       float64 `json:"price"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	url := os.Getenv("SUPABASE_URL")
	key := os.Getenv("SUPABASE_SERVICE_ROLE_KEY")

	// Тянем список приложений из Supabase (сортировка по новым)
	reqURL := fmt.Sprintf("%s/rest/v1/apps?select=*&order=created_at.desc", url)
	
	req, _ := http.NewRequest("GET", reqURL, nil)
	req.Header.Set("apikey", key)
	req.Header.Set("Authorization", "Bearer "+key)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "DB Error", 500)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	// Отдаем JSON обратно в макрос PowerPoint
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}