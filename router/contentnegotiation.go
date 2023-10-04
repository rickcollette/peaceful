package router

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"strings"
)

// Respond handles content negotiation and responds in the appropriate format
func Respond(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
	accept := r.Header.Get("Accept")

	var content string
	var marshalErr error

	if strings.Contains(accept, "application/xml") {
		w.Header().Set("Content-Type", "application/xml")
		content, marshalErr = xmlMarshalIndent(data)
	} else {
		w.Header().Set("Content-Type", "application/json")
		content, marshalErr = jsonMarshalIndent(data)
	}

	if marshalErr != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	w.Write([]byte(content))
}

func jsonMarshalIndent(data interface{}) (string, error) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func xmlMarshalIndent(data interface{}) (string, error) {
	xmlData, err := xml.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", err
	}
	return string(xmlData), nil
}
