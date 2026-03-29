package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type ConfigResult struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	// CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, apikey")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	url := os.Getenv("SUPABASE_URL")
	key := os.Getenv("SUPABASE_SERVICE_ROLE_KEY")

	// Получаем значение параметра "get"
	queryType := r.URL.Query().Get("get")
	if queryType == "" {
		queryType = "version"
	}

	var reqURL string
	
	// Если запрашивают 'all', убираем фильтр по key и запрашиваем оба поля
	if queryType == "all" {
		reqURL = fmt.Sprintf("%s/rest/v1/system_config?select=key,value", url)
	} else {
		// Для одиночных запросов (version, vercodename, verdesc и т.д.)
		reqURL = fmt.Sprintf("%s/rest/v1/system_config?key=eq.%s&select=value", url, queryType)
	}

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

	// Разделяем логику вывода
	if queryType == "all" {
		// Возвращаем полноценный JSON для списка всех настроек
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(results)
	} else {
		// Возвращаем просто строку для обратной совместимости с существующими вызовами
		if len(results) > 0 {
			fmt.Fprint(w, results[0].Value)
		} else {
			fmt.Fprint(w, "unknown")
		}
	}
}