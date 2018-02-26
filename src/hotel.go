package hotel

import "time"

type HotelsResponse struct {
	Hotels   []HotelResponse `json:"hotels"`
	CheckIn  time.Time `json:"checkIn"`
	CheckOut time.Time `json:"checkOut"`
}

type HotelResponse struct {
	HotelId            string        `json:"hotelId"`
	HotelName          string        `json:"name"`
	Brand              string        `json:"brand"`
	BrandId            string        `json:"brandId"`
	StarRating         int           `json:"starRating"`
	TotalReviewCount   int        `json:"totalReviewCount"`
	OverallGuestRating float64         `json:"overallGuestRating"`
	RatesSummary       map[string]interface{} `json:"ratesSummary"`
	HotelFeatures      map[string]interface{} `json:"hotelFeatures"`
	Location           map[string]interface{}        `json:"location"`
	HtlDealScore       int32        `json:"htlDealScore"`
}

type Hotels []Hotel

type Hotel struct {
	CheckIn            time.Time
	CheckOut           time.Time
	HotelId            string
	BrandId            string
	HotelName          string
	NeighborhoodId     string
	NeighborhoodName   string
	CityId             float64
	OverallGuestRating float64
	MinPrice           string
	StarRating         int
	TotalReviewCount   int
}