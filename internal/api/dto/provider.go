package dto

type ProvidersQueryParams struct {
	Limit  int    `schema:"limit"`
	Offset int    `schema:"offset"`
	Name   string `schema:"name"`
}
