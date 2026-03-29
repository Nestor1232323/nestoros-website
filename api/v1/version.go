package handler

import (
	"fmt"
	"net/http"
	"os"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	supabaseStorageURL := os.Getenv("SUPABASE_URL") + "/storage/v1/object/public/versions/"

	v := r.URL.Query().Get("v")
	if v == "" {
		v = "current"
	}

	extensions := map[string]string{
		"5.0b":    ".pptx",
		"5.0r":    ".pptm",
		"5.0u1":   ".pptm",
		"5.0m":    ".pptm",
		"5.0t":    ".pptm",
		"5.0tv":   ".pptm",
		"6.0":     ".pptx",
		"7.0":     ".pptm",
		"8.0":     ".pptm",
		"8.1":     ".pptm",
		"current": ".pptm",
	}

	ext, ok := extensions[v]
	if !ok {
		http.Error(w, "Version not found", 404)
		return
	}

	finalURL := fmt.Sprintf("%s%s%s", supabaseStorageURL, v, ext)

	http.Redirect(w, r, finalURL, http.StatusFound)
}