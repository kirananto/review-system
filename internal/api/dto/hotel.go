package dto

type HotelsQueryParams struct {
	Limit  int    `schema:"limit"`
	Offset int    `schema:"offset"`
	Name   string `schema:"name"`
}
