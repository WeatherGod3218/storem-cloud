package models

type Tag struct {
	TagID     string `json:"tag_id"`
	Name      string `json:"name"`
	Color     string `json:"color"`
	CreatedBy string `json:"created_by"`
}

type CreateTagRequest struct {
	Name string `json:"name"`
}
