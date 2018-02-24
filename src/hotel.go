package hotel

type Hotels struct {
	Hotels []Hotel `json:"hotels"`
}

type Hotel struct {
	HotelId       string `json:"hotelId"`
	HotelName     string `json:"name"`
	StarRating    int `json:"starRating"`
	RatesSummary  map[string]interface{} `json:"ratesSummary"`
	HotelFeatures map[string]interface{} `json:"hotelFeatures"`
}