package models

type Book struct {
	ID               string   `json:"id"`
	Title            string   `json:"title"`
	Authors          []string `json:"authors"`
	EditionCount     int      `json:"edition_count"`
	FirstPublishYear int      `json:"first_publish_year"`
	CoverImage       string   `json:"cover_image"`
	Genre            string   `json:"genre"`
}
