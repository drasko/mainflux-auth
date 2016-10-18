package models

type Scope struct {
	Actions  string `json:"actions"`
	Resource string `json:"resource"`
	Id       string `json:"id"`
}

type Scopes struct {
	Items []Scope `json:"scopes"`
}
