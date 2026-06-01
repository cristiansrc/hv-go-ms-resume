package response

// BlogTypeResponse represents the API response for a BlogType.
type BlogTypeResponse struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	NameEng string `json:"nameEng"`
}

// BlogPageResponse represents a paginated response for Blog entries.
type BlogPageResponse struct {
	Content          []BlogResponse `json:"content"`
	Pageable         PageableInfo   `json:"pageable"`
	Last             bool           `json:"last"`
	TotalPages       int            `json:"totalPages"`
	TotalElements    int64          `json:"totalElements"`
	Size             int            `json:"size"`
	Number           int            `json:"number"`
	Sort             SortInfo       `json:"sort"`
	First            bool           `json:"first"`
	NumberOfElements int            `json:"numberOfElements"`
	Empty            bool           `json:"empty"`
}

// PageableInfo represents pagination metadata.
type PageableInfo struct {
	Sort       SortInfo `json:"sort"`
	Offset     int64    `json:"offset"`
	PageSize   int      `json:"pageSize"`
	PageNumber int      `json:"pageNumber"`
	Paged      bool     `json:"paged"`
	Unpaged    bool     `json:"unpaged"`
}

// SortInfo represents sort metadata.
type SortInfo struct {
	Empty    bool `json:"empty"`
	Sorted   bool `json:"sorted"`
	Unsorted bool `json:"unsorted"`
}
