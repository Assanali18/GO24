package models

// User represents a user in the system.
// swagger:model
type User struct {
	// The user's ID.
	// example: 1
	ID int `json:"id"`
	// The user's name.
	// example: John Doe
	Name string `json:"name"`
	// The user's age.
	// example: 30
	Age int `json:"age"`
}

// UserListResponse represents a paginated list of users.
// swagger:model
type UserListResponse struct {
	// The list of users.
	Users []User `json:"users"`
	// The total number of users.
	// example: 100
	TotalItems int `json:"total_items"`
	// The current page number.
	// example: 1
	Page int `json:"page"`
	// The size of each page.
	// example: 10
	PageSize int `json:"page_size"`
	// The total number of pages.
	// example: 10
	TotalPages int `json:"total_pages"`
}
