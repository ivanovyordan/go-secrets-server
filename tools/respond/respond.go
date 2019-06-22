package respond

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"strings"
)

func Success(response http.ResponseWriter, request *http.Request, message interface{}) {
	accept := strings.ToLower(request.Header.Get("Accept"))
	var content []byte

	if accept == "application/xml" {
		content = formatXML(response, message)
	} else if accept == "application/json" {
		content = formatXML(response, message)
	} else {
		Fail(response, http.StatusNotAcceptable, "")
	}

	response.Write(content)
}

func Fail(response http.ResponseWriter, code int, message string) {
	http.Error(response, message, code)
}

func formatJSON(response http.ResponseWriter, message interface{}) []byte {
	response.Header().Set("Content-Type", "application/json")
	content, _ := json.Marshal(message)

	return content
}

func formatXML(response http.ResponseWriter, message interface{}) []byte {
	response.Header().Set("Content-Type", "application/xml")

	XMLHeader := `<?xml version="1.0" encoding="UTF-8"?>` + "\n"
	content, _ := xml.MarshalIndent(message, "", "  ")
	content = []byte(XMLHeader + string(content))

	return content
}
