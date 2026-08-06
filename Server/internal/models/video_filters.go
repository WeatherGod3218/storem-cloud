package models

type FilterDirection string

const (
	Ascending  FilterDirection = "ascending"
	Descending FilterDirection = "descending"
)

type FilterElement string

const (
	DateCreated  FilterElement = "date_created"
	DateUploaded FilterElement = "date_uploaded"
	Filename     FilterElement = "filename"
	Title        FilterElement = "title"
)

type Filter struct {
	Title           *string         `json:"title"`
	FilterElement   FilterElement   `json:"filter_element"`
	FilterDirection FilterDirection `json:"filter_direction"`
	FilterTags      []string        `json:"filter_tags"`
}
