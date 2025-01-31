package models

type DirectMessageRequest struct {
	SourceAccountId string
	TargetAccountId string
	Content         string
}
