package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type Message struct {
	Text      string `json:"text"`
	CreatedAt string `json:"created_at,omitempty"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	url := os.Getenv("SUPABASE_URL")
	key := os.Getenv("SUPABASE_SERVICE_ROLE_KEY")

	switch r.Method {
	case http.MethodGet:
		// Аналог getLatestMessagesText()
		fetchMessages(w, url, key)
	case http.MethodPost:
		// Аналог postSendMessage()
		sendMessage(w, r, url, key)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func fetchMessages(w http.ResponseWriter, url, key string) {
	// Берем последние 50 сообщений
	reqURL := fmt.Sprintf("%s/rest/v1/messages?select=text&order=created_at.asc&limit=50", url)
	
	req, _ := http.NewRequest("GET", reqURL, nil)
	req.Header.Set("apikey", key)
	req.Header.Set("Authorization", "Bearer "+key)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Service Unavailable", 503)
		return
	}
	defer resp.Body.Close()

	var msgs []Message
	json.NewDecoder(resp.Body).Decode(&msgs)

	// Собираем всё в одну строку через перенос строки для TextBox2
	var builder strings.Builder
	for _, m := range msgs {
		builder.WriteString(m.Text + "\n")
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprint(w, builder.String())
}

func sendMessage(w http.ResponseWriter, r *http.Request, url, key string) {
	var msg Message
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil || msg.Text == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	jsonData, _ := json.Marshal(msg)
	reqURL := fmt.Sprintf("%s/rest/v1/messages", url)

	req, _ := http.NewRequest("POST", reqURL, bytes.NewBuffer(jsonData))
	req.Header.Set("apikey", key)
	req.Header.Set("Authorization", "Bearer "+key)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, _ := client.Do(req)
	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)
	fmt.Fprint(w, `{"status":"OK"}`)
}