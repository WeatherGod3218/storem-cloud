package models

import "time"

type GetVideoGroupPart struct {
	RowID string `json:"row_id"`
	S3ID  string `json:"s3_id"`

	CustomTitle       *string `json:"custom_title"`
	CustomDescription *string `json:"custom_description"`

	Visibility string `json:"visibility"`

	UserId       string `json:"user_id"`
	Timestamp    any    `json:"timestamp"`
	ThumbnailURL string `json:"thumbnail"`
	Filename     string `json:"filename"`
}

type GetVideoGroupPartResponse struct {
	RowID string `json:"row_id"`
	S3ID  string `json:"s3_id"`

	CustomTitle       *string `json:"custom_title"`
	CustomDescription *string `json:"custom_description"`

	Visibility string `json:"visibility"`

	Username     string `json:"username"`
	Filename     string `json:"filename"`
	ThumbnailURL string `json:"thumbnail"`
	Timestamp    any    `json:"timestamp"`
}

type GetVideoGroupCursor struct {
	Timestamp any    `json:"timestamp"`
	RowID     string `json:"row_id"`
}

type GetVideoDataResponse struct {
	RowID string `json:"row_id"`
	S3ID  string `json:"s3_id"`

	CustomTitle       *string `json:"custom_title"`
	CustomDescription *string `json:"custom_description"`

	Visibility string `json:"visibility"`

	Username string `json:"username"`
	Filename string `json:"filename"`

	VideoURL     string `json:"video_url"`
	ThumbnailURL string `json:"thumbnail_url"`

	Timestamp any `json:"timestamp"`

	CanModify bool `json:"can_modify"`
}

type GetVideoDataDatabase struct {
	RowID string `json:"row_id"`
	S3ID  string `json:"s3_id"`

	CustomTitle       *string `json:"custom_title"`
	CustomDescription *string `json:"custom_description"`

	Visibility string `json:"visibility"`

	UserID    string `json:"user_id"`
	Filename  string `json:"filename"`
	Timestamp any    `json:"timestamp"`
}

type GetVideoGroupRequest struct {
	Timestamp *any    `json:"timestamp"`
	RowID     *string `json:"row_id"`
	Filters   Filter  `json:"filter"`
}

type GetVideoGroupResponse struct {
	Videos []GetVideoGroupPartResponse `json:"videos"`
	Cursor *GetVideoGroupCursor        `json:"cursor"`
}

type GetRandomVideoResponse struct {
	RowID string `json:"row_id"`
}

type ChangeVideoTitleRequest struct {
	RowID string `json:"row_id"`
	Title string `json:"title"`
}

type ChangeVideoDescriptionRequest struct {
	RowID       string `json:"row_id"`
	Description string `json:"description"`
}

type ChangeVideoVisibilityRequest struct {
	RowID      string `json:"row_id"`
	Visibility string `json:"visibility"`
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
	VideoID     string    `json:"video_id"`
	UserID      string    `json:"user_id"`
}

type VideoBackupProcessed struct {
	Filename       string  `json:"filename"`
	FileLength     float64 `json:"file_length"`
	FileSize       int64   `json:"file_size"`
	VideoS3URL     string  `json:"video_s3_url"`
	ThumbnailS3URL string  `json:"thumbnail_s3_url"`
}
