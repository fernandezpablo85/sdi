package api

type AssetResponse struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type RateLimitHeaders struct {
	Limit     int
	Remaining int
	Reset     int64 // Unix timestamp
}
