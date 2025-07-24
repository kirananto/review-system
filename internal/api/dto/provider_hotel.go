package dto

type ProviderHotelsQueryParams struct {
	Limit      int  `schema:"limit"`
	Offset     int  `schema:"offset"`
	HotelID    uint `schema:"hotel_id"`
	ProviderID uint `schema:"provider_id"`
}
