package handlers

type createShortURLRequest struct {
	OriginalURL string `json:"original_url"`
}

type createShortURLResponse struct {
	ShortURL string `json:"short_url"`
}

type getStatsByShortURLResponse struct {
	RequestNumber int64
}
