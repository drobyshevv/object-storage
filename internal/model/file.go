package model

import "time"

type File struct {
	ID          int
	Filename    string
	Size        int64
	ContentType string
	S3Key       string
	CreatedAt   time.Time
}
