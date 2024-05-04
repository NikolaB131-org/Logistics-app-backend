package entity

type Order struct {
	ID    string `json:"id"`
	Items []Item `json:"items"`
}

type Item struct {
	ID       string `json:"id"`
	Quantity int    `json:"quantity"`
}
