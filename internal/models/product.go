package models

// Product models
type Product struct {
	ProductID   string   `json:"productId"`
	CategoryID  string   `json:"categoryId"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       float64  `json:"price"`
	ImageURL    string   `json:"imageUrl"`
	Photos      []string `json:"photos"`
	Quantity    int64    `json:"quantity"`
	Rating      int      `json:"rating"`
}
