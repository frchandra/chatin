package validation

type RegisterValidation struct {
	Name     string `json:"name" binding:"required,min=3,max=36"`
	Email    string `json:"email" binding:"required,email,min=5,max=36"`
	Password string `json:"password" binding:"required,min=7"`
}

type LoginValidation struct {
	Name     string `json:"name" binding:"required_without=Email,min=0,max=36"`
	Email    string `json:"email" binding:"required_without=Username,min=0,max=36"`
	Password string `json:"password" binding:"required,min=7"`
}
