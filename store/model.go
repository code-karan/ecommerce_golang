type Product struct {
	ID int			`bson:"_id"`
	Title int		`json:"title"`
	Image string	`json:"image"`
	Price uint64	`json:"price"`
	Rating uint64	`json:"rating"`
}

// Array of product objects

type Products []Product