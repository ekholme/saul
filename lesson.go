package main

import "github.com/sashabaranov/go-openai"

type LessonRequest struct {
	Grade          string `json:"grade"`
	ItemDescriptor string `json:"itemDescriptor"`
}

// add a constructor
//I don't think I need this for now
// func NewLessonRequest(grade string, id string) *LessonRequest {
// 	return &LessonRequest{
// 		Grade:          grade,
// 		ItemDescriptor: id,
// 	}
// }

// method to create a ChatGPT message request from a lesson request type
func (lr *LessonRequest) CreateGPTMessage() []openai.ChatCompletionMessage {

	s := "Plan a lesson for " + lr.Grade + " graders on " + lr.ItemDescriptor

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
