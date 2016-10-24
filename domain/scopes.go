package domain

type Scope struct {
	Actions  string `json:"actions"`
	Resource string `json:"resource"`
	Id       string `json:"id"`
}

type Payload struct {
	Scopes []Scope `json:"scopes"`
}
