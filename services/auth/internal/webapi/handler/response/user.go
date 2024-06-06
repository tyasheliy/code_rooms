package response

type UserResponse struct {
	Id     int64  `json:"id"`
	Login  string `json:"login"`
	RoleId int    `json:"role_id"`
}
