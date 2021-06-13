package store

type User struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

type JwtToken struct {
    Token string `json:"token"`
}

type Exception struct {
    Message string `json:"message"`
}

type Product struct {
	ID int			`bson:"_id"`
	Title string	`bson:"title" json:"title"`
	Image string	`bson:"image" json:"image"`
	Price uint64	`bson:"price" json:"price"`
	Rating uint64	`bson:"rating" json:"rating"`
}
