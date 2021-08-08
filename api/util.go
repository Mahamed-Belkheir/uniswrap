package api

import (
	"encoding/json"
	"net/http"
)

type m map[string]interface{}

func (m m) send(code int, w http.ResponseWriter) {
	data, err := json.Marshal(m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(code)
		w.Write(data)
	}
}

func checkQuery(query string, w http.ResponseWriter, r *http.Request) (string, bool) {
	q := r.URL.Query().Get(query)
	if q == "" {
		m{
			"status":  "error",
			"message": "missing query parameter " + query,
		}.send(400, w)
		return "", false
	}
	return q, true
}
