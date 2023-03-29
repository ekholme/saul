package saul

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sashabaranov/go-openai"
)

// define a server type
type Server struct {
	Router    *mux.Router
	Srvr      *http.Server
	GPTClient *openai.Client
	Templates *template.Template
}

func NewServer(r *mux.Router, client *openai.Client, t *template.Template) *Server {
	listenAddr := ":8080"

	return &Server{
		Router: r,
		Srvr: &http.Server{
			Addr: listenAddr,
		},
		GPTClient: client,
		Templates: t,
	}
}

// register routes
func (s *Server) registerRoutes() {
	s.Router.HandleFunc("/", s.handleIndex).Methods("GET")
	s.Router.HandleFunc("/", s.handleMockLesson).Methods("POST")
	// s.Router.HandleFunc("/", s.handleRequestLesson).Methods("POST")
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
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	s.Templates.ExecuteTemplate(w, "index.html", nil)
}

func (s *Server) handleRequestLesson(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	lr := &LessonRequest{
		Grade:          r.FormValue("grade"),
		ItemDescriptor: r.FormValue("itemDescriptor"),
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

	l := NewLessonResponse(lr, resp.Choices[0].Message.Content)

	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	s.Templates.ExecuteTemplate(w, "lesson_plan.html", l)

}

func (s *Server) handleMockLesson(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	lr := &LessonRequest{
		Grade:          r.FormValue("grade"),
		ItemDescriptor: r.FormValue("itemDescriptor"),
	}

	m := "this is a mock response for " + lr.Grade + " graders and a lesson on " + lr.ItemDescriptor

	l := NewLessonResponse(lr, m)

	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	s.Templates.ExecuteTemplate(w, "lesson_plan.html", l)
}

// writeJSON helper
func WriteJSON(w http.ResponseWriter, statusCode int, v any) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	json.NewEncoder(w).Encode(v)
}
