package token

type TokenUser struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Image      string `json:"image"`
	IsVerified bool   `json:"is_verified"`
}
