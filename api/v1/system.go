package handler

import (
	"fmt"
	"net/http"
	"os"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	url := os.Getenv("SUPABASE_URL")
	key := os.Getenv("SUPABASE_SERVICE_ROLE_KEY")

	// Получаем то, что просит макрос: version, vercodename или verdesc
	queryType := r.URL.Query().Get("get")
	if queryType == "" {
		queryType = "version" // по умолчанию
	}

	// Запрос в Supabase к таблице system_config
	// Используем фильтр по ключу .eq
	reqURL := fmt.Sprintf("%s/rest/v1/system_config?key=eq.%s&select=value", url, queryType)
	
	req, _ := http.NewRequest("GET", reqURL, nil)
	req.Header.Set("apikey", key)
	req.Header.Set("Authorization", "Bearer "+key)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		fmt.Fprint(w, "Error")
		return
	}
	defer resp.Body.Close()

	// Supabase REST возвращает массив: [{"value": "8.5"}]
	// Для простоты макроса вытащим только чистое значение, если нужно.
	// Но проще всего сделать так:
	
	type ConfigResult struct {
		Value string `json:"value"`
	}
	var results []ConfigResult
	import "encoding/json"
	json.NewDecoder(resp.Body).Decode(&results)

	if len(results) > 0 {
		fmt.Fprint(w, results[0].Value) // Отдаем чистый текст для совместимости с VBA
	} else {
		fmt.Fprint(w, "unknown")
	}
}