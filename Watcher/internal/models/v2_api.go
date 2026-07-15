package models

import "time"

type GetVideoGroupPart struct {
	RowID        string    `json:"row_id"`
	VideoS3Id    string    `json:"video_s3_id"`
	Timestamp    time.Time `json:"timestamp"`
	ThumbnailURL string    `json:"thumbnail"`
	Filename     string    `json:"filename"`
}

type TusMetadata struct {
	FileName    string `json:"filename"`
	FileSize    int    `json:"filesize"`
	VideoLength int    `json:"filelength"`
}
