package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	// cors
	w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, apikey")

    if r.Method == "OPTIONS" {
        w.WriteHeader(http.StatusOK)
        return
    }

	url := os.Getenv("SUPABASE_URL")
	key := os.Getenv("SUPABASE_SERVICE_ROLE_KEY")

	// Получаем версию из параметра ?v=
	versionParam := r.URL.Query().Get("v")

	// Базовый URL запроса с сортировкой
	// По умолчанию тянем всё, если параметр не передан
	filter := ""
	if versionParam == "8" {
		// Только нативные для nestorOS 8
		filter = "&is_n8=eq.true"
	} else if versionParam == "10" {
		// Только старые/сторонние приложения
		filter = "&is_n8=eq.false"
	}

	// Формируем финальный запрос к Supabase
	reqURL := fmt.Sprintf("%s/rest/v1/apps?select=*&order=created_at.desc%s", url, filter)
	
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

	// Читаем ответ от Supabase
	body, _ := io.ReadAll(resp.Body)

	// Отдаем JSON в PowerPoint
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}