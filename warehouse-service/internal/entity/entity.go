package entity

type Product struct {
	ID       string  `db:"id" json:"id"`
	Name     string  `db:"name" json:"name"`
	Quantity int     `db:"quantity" json:"quantity"`
	Price    float32 `db:"price" json:"price"`
}
