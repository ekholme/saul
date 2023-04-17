package saul

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/sashabaranov/go-openai"
)

// define a server type
type Server struct {
	Router             *mux.Router
	Srvr               *http.Server
	GPTClient          *openai.Client
	Templates          *template.Template
	PerformanceService *PerformanceService
	TestService        *TestService
	SchoolService      *SchoolService
}

func NewServer(r *mux.Router, client *openai.Client, t *template.Template, ps *PerformanceService, ts *TestService, ss *SchoolService) *Server {
	listenAddr := ":8080"

	return &Server{
		Router: r,
		Srvr: &http.Server{
			Addr: listenAddr,
		},
		GPTClient:          client,
		Templates:          t,
		PerformanceService: ps,
		TestService:        ts,
		SchoolService:      ss,
	}
}

// register routes
func (s *Server) registerRoutes() {
	s.Router.HandleFunc("/", s.handleIndex).Methods("GET")
	s.Router.HandleFunc("/free", s.handleRequestLesson).Methods("POST")
	s.Router.HandleFunc("/free", s.handleFree).Methods("GET")
	s.Router.HandleFunc("/guided", s.handleSchool).Methods("GET")
	s.Router.HandleFunc("/guided", s.handleTestRedirect).Methods("POST")
	s.Router.HandleFunc("/guided/{school}", s.handleGetTestsBySchool).Methods("GET")
	// s.Router.HandleFunc("/", s.handleMockLesson).Methods("POST")

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

func (s *Server) handleFree(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	s.Templates.ExecuteTemplate(w, "free.html", nil)

}

func (s *Server) handleRequestLesson(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	lr := &LessonRequest{
		Grade:          r.FormValue("grade"),
		ItemDescriptor: r.FormValue("itemDescriptor"),
		StudentPop:     r.FormValue("studentPop"),
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

// function to mock out creating a lesson
// useful for testing UI without making requests to OpenAI
func (s *Server) handleMockLesson(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	lr := &LessonRequest{
		Grade:          r.FormValue("grade"),
		ItemDescriptor: r.FormValue("itemDescriptor"),
		StudentPop:     r.FormValue("studentPop"),
	}

	var m string

	if lr.StudentPop == "all students" {
		m = "this is a mock response for " + lr.Grade + " graders and a lesson on " + lr.ItemDescriptor
	} else {
		m = "this is a mock response for " + lr.Grade + " grade " + lr.StudentPop + " and a lesson on " + lr.ItemDescriptor
	}

	l := NewLessonResponse(lr, m)

	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	s.Templates.ExecuteTemplate(w, "lesson_plan.html", l)
}

// get page for schools
func (s *Server) handleSchool(w http.ResponseWriter, r *http.Request) {

	ctx := context.Background()

	schs, err := s.SchoolService.GetAllSchools(ctx)

	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	s.Templates.ExecuteTemplate(w, "schools.html", schs)
}

// redirect to test select page
func (s *Server) handleTestRedirect(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	sch := r.FormValue("schName")

	su := url.QueryEscape(sch)

	u, err := url.JoinPath("/", "guided", su)

	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, err)
	}

	http.Redirect(w, r, u, http.StatusSeeOther)
}

// handle getting tests
func (s *Server) handleGetTestsBySchool(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	ctx := context.Background()

	su := vars["school"]

	sch, err := url.QueryUnescape(su)

	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, err)
	}

	tsts, err := s.TestService.GetTestBySchool(ctx, sch)

	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, err)
	}

	u, err := url.JoinPath("/", "guided", su)

	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, err)
	}

	tr := &TestRequest{
		URL:   u,
		Tests: tsts,
	}

	// WriteJSON(w, http.StatusOK, schs)
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	s.Templates.ExecuteTemplate(w, "tests.html", tr)

}

// writeJSON helper
func WriteJSON(w http.ResponseWriter, statusCode int, v any) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	json.NewEncoder(w).Encode(v)
}
