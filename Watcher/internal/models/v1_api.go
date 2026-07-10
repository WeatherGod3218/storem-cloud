package models

import "time"

type VideoStartBackupRequest struct {
	FileName   string  `json:"filename"`
	FileLength float64 `json:"filelength"`
	FileSize   int64   `json:"filesize"`
}

type VideoStartBackupResponse struct {
	RowID           string         `json:"row_id"`
	VideoS3URLs     []VideoURLPart `json:"video_s3_urls"`
	VideoS3UploadID string         `json:"video_s3_upload_id"`
	ThumbnailS3URL  string         `json:"thumbnail_s3_url"`
}

type VideoCompleteBackupRequest struct {
	RowID           string               `json:"row_id"`
	Filename        string               `json:"filename"`
	VideoS3UploadID string               `json:"video_s3_upload_id"`
	CompletedParts  []VideoCompletedPart `json:"completed_parts"`
}

type VideoCompleteBackupResponse struct {
	RowID           string         `json:"row_id"`
	VideoS3URLs     []VideoURLPart `json:"video_s3_urls"`
	VideoS3UploadID string         `json:"video_s3_upload_id"`
	ThumbnailS3URL  string         `json:"thumbnail_s3_url"`
}

type VideoAbortBackupRequest struct {
	RowID           string `json:"row_id"`
	VideoS3UploadID string `json:"video_s3_upload_id"`
	Filename        string `json:"filename"`
}

type GetVideoGroupPart struct {
	Timestamp    time.Time `json:"offset"`
	RowId        string    `json:"row_id"`
	VideoURL     string    `json:"video_url"`
	ThumbnailURL string    `json:"thumbnail_url"`
}

type GetVideoGroupRequest struct {
	Timestamp *time.Time `json:"offset"`
	RowId     string     `json:"row_id"`
}

type AccessVideoResponse struct {
	URL   string `json:"url"`
	Video string `json:"video"`
}

type GetVideoGroupResponse struct {
	Videos []GetVideoGroupPart `json:"videos"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Success bool `json:"success"`
}
