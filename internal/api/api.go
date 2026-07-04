package api

import (
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
	dict *dictionary.Dictionary
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

	if err := a.dict.Insert(req.Eng, req.Rus); err != nil {
		http.Error(w, "Failed to add word", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (a *API) StartServer(port string) {
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

	go func() {
		http.ListenAndServe(port, mux)
	}()
}
