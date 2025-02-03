package infrastructure

import (
	"Infocenter/Application"
	"Infocenter/Domain"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type HTTPHandler struct {
	service *application.MessageService
}

func NewHTTPHandler(service *application.MessageService) *HTTPHandler {
	return &HTTPHandler{service: service}
}

func (h *HTTPHandler) PostMessage(w http.ResponseWriter, r *http.Request) {
	topicName := r.URL.Path[len("/infocenter/"):]
	if topicName == "" {
		http.Error(w, "Please add a topic.", http.StatusBadRequest)
		return
	}

	data, messageError := io.ReadAll(r.Body)
	if messageError != nil || len(data) == 0 {
		http.Error(w, "Incorrect message body.", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if _, serverError := h.service.PublishMessage(topicName, string(data)); serverError != nil {
		http.Error(w, serverError.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *HTTPHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
	topicName := r.URL.Path[len("/infocenter/"):]
	if topicName == "" {
		http.Error(w, "Missing topic", http.StatusBadRequest)
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Date", time.Now().UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT"))

	topic := h.service.GetTopic(topicName)
	client := make(chan domain.Message, 10)
	topic.Subscribe(client)

	defer topic.Unsubscribe(client)

	timeout := time.NewTimer(30 * time.Second)
	defer timeout.Stop()

	for {
		select {
		case msg, ok := <-client:
			if !ok {
				return
			}
			fmt.Fprintf(w, "id: %d\nevent: msg\ndata: %s\n\n", msg.ID, msg.Data)
			flusher.Flush()

		case <-timeout.C:
			fmt.Fprintf(w, "event: timeout\ndata: 30s\n\n")
			flusher.Flush()
			return
		}
	}
}

func StartServer(service *application.MessageService) {
	handler := NewHTTPHandler(service)
	http.HandleFunc("/infocenter/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handler.PostMessage(w, r)
		case http.MethodGet:
			handler.GetMessages(w, r)
		default:
			http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
		}
	})

	log.Println("Server running on :8080")
	if serverError := http.ListenAndServe(":8080", nil); serverError != nil {
		log.Fatalf("Server failed: %v", serverError)
	}
}
