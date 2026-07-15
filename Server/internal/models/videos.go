package models

import "time"

type GetVideoGroupPart struct {
	RowID string `json:"row_id"`
	S3Id  string `json:"s3_id"`

	CustomTitle       *string `json:"custom_title"`
	CustomDescription *string `json:"custom_description"`

	UserId       string    `json:"user_id"`
	Timestamp    time.Time `json:"timestamp"`
	ThumbnailURL string    `json:"thumbnail"`
	Filename     string    `json:"filename"`
}

type GetVideoGroupPartResponse struct {
	RowID string `json:"row_id"`
	S3Id  string `json:"s3_id"`

	CustomTitle       *string `json:"custom_title"`
	CustomDescription *string `json:"custom_description"`

	Username     string    `json:"username"`
	Filename     string    `json:"filename"`
	ThumbnailURL string    `json:"thumbnail"`
	Timestamp    time.Time `json:"timestamp"`
}

type GetVideoGroupCursor struct {
	Timestamp time.Time `json:"timestamp"`
	RowID     string    `json:"row_id"`
}

type GetVideoGroupRequest struct {
	Timestamp *time.Time `json:"timestamp"`
	RowID     string     `json:"row_id"`
}

type GetVideoDataResponse struct {
	RowID string `json:"row_id"`
	S3Id  string `json:"s3_id"`

	CustomTitle       *string `json:"custom_title"`
	CustomDescription *string `json:"custom_description"`

	Username  string    `json:"username"`
	Filename  string    `json:"filename"`
	VideoURL  string    `json:"video_url"`
	Timestamp time.Time `json:"timestamp"`
}

type GetVideoDataDatabase struct {
	RowID string `json:"row_id"`
	S3Id  string `json:"s3_id"`

	CustomTitle       *string `json:"custom_title"`
	CustomDescription *string `json:"custom_description"`

	UserId    string    `json:"user_id"`
	Filename  string    `json:"filename"`
	Timestamp time.Time `json:"timestamp"`
}

type GetVideoGroupResponse struct {
	Videos []GetVideoGroupPartResponse `json:"videos"`
	Cursor *GetVideoGroupCursor        `json:"cursor"`
}

type AccessVideoResponse struct {
	URL   string `json:"url"`
	Video string `json:"video"`
}

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

type VideoDatabaseEntry struct {
	Filename    string    `json:"filename"`
	FileLength  float64   `json:"file_length"`
	FileSize    int64     `json:"file_size"`
	FileModDate time.Time `json:"file_mod_date"`
	VideoId     string    `json:"video_id"`
	UserId      string    `json:"user_id"`
}

type VideoBackupProcessed struct {
	Filename       string  `json:"filename"`
	FileLength     float64 `json:"file_length"`
	FileSize       int64   `json:"file_size"`
	VideoS3URL     string  `json:"video_s3_url"`
	ThumbnailS3URL string  `json:"thumbnail_s3_url"`
}
