package dto

type HotelRequestBody struct {
	HotelName string `json:"hotel_name" validate:"required"`
}

type HotelsQueryParams struct {
	Limit  int    `schema:"limit"`
	Offset int    `schema:"offset"`
	Name   string `schema:"name"`
}
