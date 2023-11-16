package entity

import "time"

type S3Object struct {
	Key          *string    `json:"name"`
	LastModified *time.Time `json:"last_modified"`
	Size         int64      `json:"size"`
}

type S3ObjectTree struct {
	CommonPrefixes        []*string   `json:"subfolders"`
	Contents              []*S3Object `json:"files"`
	Prefix                *string     `json:"root_path"`
	NextContinuationToken *string     `json:"next_continuation_token"`
}
