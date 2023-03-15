package saul

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sashabaranov/go-openai"
)

// define a server type
type Server struct {
	Router    *mux.Router
	Srvr      *http.Server
	GPTClient *openai.Client
}

func NewServer(r *mux.Router, client *openai.Client) *Server {
	listenAddr := ":8080"

	return &Server{
		Router: r,
		Srvr: &http.Server{
			Addr: listenAddr,
		},
		GPTClient: client,
	}
}

// register routes
func (s *Server) registerRoutes() {
	s.Router.HandleFunc("/", s.handleIndex).Methods("GET")
	s.Router.HandleFunc("/lesson", s.handleRequestLesson).Methods("POST")
}

// method to run the server
func (s *Server) Run() {
	s.registerRoutes()

	fmt.Printf("Saul running on port %s", s.Srvr.Addr)

	s.Srvr.Handler = s.Router

	s.Srvr.ListenAndServe()
}

// define some handlers
func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	msg := make(map[string]string)

	msg["welcome"] = "Welcome to Saul"

	WriteJSON(w, 200, msg)
}

func (s *Server) handleRequestLesson(w http.ResponseWriter, r *http.Request) {
	var lr *LessonRequest

	err := json.NewDecoder(r.Body).Decode(&lr)

	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, err)
		return
	}

	m := lr.CreateGPTMessage()

	req := openai.ChatCompletionRequest{
		Model:    openai.GPT3Dot5Turbo,
		Messages: m,
	}

	ctx := context.Background()

	resp, err := s.GPTClient.CreateChatCompletion(ctx, req)

	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, err)
		return
	}

	l := NewLessonResponse(resp.Choices[0].Message.Content)

	WriteJSON(w, http.StatusOK, l)
}

// writeJSON helper
func WriteJSON(w http.ResponseWriter, statusCode int, v any) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	json.NewEncoder(w).Encode(v)
}
