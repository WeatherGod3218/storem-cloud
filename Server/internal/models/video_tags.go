package models

type VideoTag struct {
	TagID string `json:"tag_id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

type ModifyVideoTagRequest struct {
	VideoID string `json:"video_id"`
	TagID   string `json:"tag_id"`
	User    string `json:"user"`
}
