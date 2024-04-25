package model

type ModifyRole struct {
	User  string   `json:"user"`
	Roles []string `json:"roles"`
	Del   string   `json:"del"`
}
