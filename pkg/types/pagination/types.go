package pagination

type Pagination struct {
	Page  int `query:"page" json:"page"`
	Limit int `query:"limit" json:"limit"`
}

type Metadata struct {
	Pagination
	Total   int64 `json:"total"`
	Count   int   `json:"count"`
	HasPrev bool  `json:"hasPrev"`
	HasNext bool  `json:"hasNext"`
}

type PaginatedResult[T any] struct {
	Metadata Metadata `json:"metadata"`
	Result   []*T     `json:"result"`
}
