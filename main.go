package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type APIResponse struct {
	Data string `json:"data"`
}

const greenAPIBaseURL = "https://api.green-api.com"

// Шаблон HTML страницы
func serveTemplate(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(w, "Unable to load template", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

// Обработчик для вызова методов GREEN-API
func handleAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()
	idInstance := r.FormValue("idInstance")
	apiToken := r.FormValue("apiToken")
	method := r.FormValue("method")
	phone := r.FormValue("phone")
	message := r.FormValue("message")
	fileURL := r.FormValue("fileURL")

	// Сбор URL и данных для запроса
	var apiURL string
	var payload string

	switch method {
	case "getSettings":
		apiURL = fmt.Sprintf("%s/waInstance%s/getSettings/%s", greenAPIBaseURL, idInstance, apiToken)
	case "getStateInstance":
		apiURL = fmt.Sprintf("%s/waInstance%s/getStateInstance/%s", greenAPIBaseURL, idInstance, apiToken)
	case "sendMessage":
		apiURL = fmt.Sprintf("%s/waInstance%s/sendMessage/%s", greenAPIBaseURL, idInstance, apiToken)
		payload = fmt.Sprintf(`{"chatId": "%s@c.us", "message": "%s"}`, phone, message)
	case "sendFileByUrl":
		apiURL = fmt.Sprintf("%s/waInstance%s/sendFileByUrl/%s", greenAPIBaseURL, idInstance, apiToken)
		payload = fmt.Sprintf(`{"chatId": "%s@c.us", "urlFile": "%s", "fileName": "file.jpg"}`, phone, fileURL)
	default:
		http.Error(w, "Invalid method", http.StatusBadRequest)
		return
	}

	// Выполнение HTTP-запроса к GREEN-API
	var req *http.Request
	var err error
	if payload != "" {
		req, err = http.NewRequest("POST", apiURL, strings.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequest("GET", apiURL, nil)
	}

	if err != nil {
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Error executing request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error reading response", http.StatusInternalServerError)
		return
	}

	// Возвращаем ответ клиенту
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func main() {
	http.HandleFunc("/", serveTemplate)
	http.HandleFunc("/api", handleAPI)

	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
