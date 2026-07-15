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

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Success bool `json:"success"`
}
