package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/g0disd3ad/rbt/internal/dictionary"
)

type TranslationRequest struct {
	Eng string `json:"eng"`
	Rus string `json:"rus"`
}

type API struct {
	dict   *dictionary.Dictionary
	server *http.Server
}

func isStrictEnglish(word string) bool {
	if word == "" {
		return false
	}
	for _, r := range word {
		if r == '-' || r == ' ' || (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
			continue
		}
		return false
	}
	return true
}

func isStrictRussian(word string) bool {
	if word == "" {
		return false
	}
	for _, r := range word {
		if r == '-' || r == ' ' || (r >= 'а' && r <= 'я') || (r >= 'А' && r <= 'Я') || r == 'ё' || r == 'Ё' {
			continue
		}
		return false
	}
	return true
}

func NewAPI(d *dictionary.Dictionary) *API {
	return &API{dict: d}
}

func (a *API) handleSearch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	word := strings.TrimSpace(r.URL.Query().Get("word"))
	if word == "" {
		http.Error(w, "Missing 'word' parameter", http.StatusBadRequest)
		return
	}

	if !isStrictEnglish(word) {
		http.Error(w, "Invalid format: 'eng' word must be in English, only English letters allowed", http.StatusBadRequest)
		return
	}

	translations, err := a.dict.Search(word)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"word":         word,
		"translations": translations,
	})
}

func (a *API) handleAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req TranslationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	req.Eng = strings.TrimSpace(req.Eng)
	req.Rus = strings.TrimSpace(req.Rus)

	if !isStrictEnglish(req.Eng) || !isStrictRussian(req.Rus) {
		http.Error(w, "Invalid format: 'eng' word must be in English, 'rus' word must be in Russian", http.StatusBadRequest)
		return
	}

	if err := a.dict.Insert(req.Eng, req.Rus); err != nil {
		http.Error(w, "Failed to add word", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (a *API) StartServer(port string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/translate", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			a.handleSearch(w, r)
		case http.MethodPost:
			a.handleAdd(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	a.server = &http.Server{
		Addr:    port,
		Handler: mux,
	}

	return a.server.ListenAndServe()
}

func (a *API) Shutdown(ctx context.Context) error {
	if a.server != nil {
		return a.server.Shutdown(ctx)
	}
	return nil
}
