package models

type QuestionRequest struct {
	Question string `json:"question" binding:"required"`
	Model    string `json:"model" binding:"required"`
}

type UploadRequest struct {
	Filename string `json:"filename"`
	Content  string `json:"content"`
}
