package models

type VideoURLPart struct {
	RequestURL string `json:"request_url"`
	Offset     int64  `json:"offset"`
	Size       int64  `json:"size"`
	PartNumber int32  `json:"part_number"`
}

type VideoCompletedPart struct {
	PartNumber int    `json:"part_number"`
	ETag       string `json:"etag"`
}

type VideoPart struct {
	PartNumber int32
	Offset     int64
	Size       int64
}

type VideoBackupProcessed struct {
	FileName       string  `json:"filename"`
	FileLength     float64 `json:"filelength"`
	FileSize       int64   `json:"filesize"`
	VideoS3URL     string  `json:"video_s3_url"`
	ThumbnailS3URL string  `json:"thumbnail_s3_url"`
}
