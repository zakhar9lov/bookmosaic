package models

type Book struct {
	ID         int    `json:"id"`
	ISBN       string `json:"isbn"`
	Title      string `json:"title"`
	Author     string `json:"author"`
	Year       int    `json:"year"`
	Summary    string `json:"summary"`
	CoverImage string `json:"cover"`
}

type SummaryBook struct {
	ID         int    `json:"id"`
	Title      string `json:"title"`
	CoverImage string `json:"cover"`
}
