package animal

type animal struct {
	Name       string `json:"name" binding:"required"`
	Species    string `json:"species" binding:"required"`
	Breed      string `json:"breed"`
	Age        int    `json:"age"`
	OwnerEmail string `json:"owner_email" binding:"required,email"`
}

// outputs
type CreateAnimalOutput struct {
	Message string `json:"message"`
}
