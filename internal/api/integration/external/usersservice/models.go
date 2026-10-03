package usersservice

type User struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CompanyID string `json:"company_id"`
	Active    bool   `json:"active"`
}
