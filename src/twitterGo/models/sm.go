package models

type Secret struct {
	Host     string `json:"host"`
	UserName string `json:"username"`
	Password string `json:"password"`
	JWTSing  string `json:"jwtsing"`
	Database string `json:"database"`
}
