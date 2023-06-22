package saul

import "github.com/sashabaranov/go-openai"

type LessonRequest struct {
	Grade          string `json:"grade"`
	ItemDescriptor string `json:"itemDescriptor"`
	BestPractice   string `json:"bestPractice"`
	StudentPop     string `json:"studentPop"`
}

// method to create a ChatGPT message request from a lesson request type
func (lr *LessonRequest) CreateGPTMessage() []openai.ChatCompletionMessage {

	var s string

	if lr.StudentPop == "all students" {
		s = "Plan a lesson for " + lr.Grade + " students on " + lr.ItemDescriptor
	} else {
		s = "Plan a lesson for " + lr.Grade + " grade " + lr.StudentPop + " on " + lr.ItemDescriptor
	}

	s = s + " Please " + lr.BestPractice + " in the lesson plan."

	m := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleUser,
			Content: s,
		},
	}

	return m
}

// define a lesson response type
type LessonResponse struct {
	LessonRequest *LessonRequest
	Lesson        string `json:"lesson"`
}

func NewLessonResponse(lr *LessonRequest, l string) *LessonResponse {
	return &LessonResponse{
		LessonRequest: lr,
		Lesson:        l,
	}
}
