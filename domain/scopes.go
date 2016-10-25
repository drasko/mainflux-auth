package domain

// Scope represents a resource(s) access scope definition. Each definition
// consists of allowed actions, resource type and its identifier. Keep in
// mind that the '*' is treated as wild card - it can be used to indicate
// "all" available resource types and "all" resource identifiers.
type Scope struct {
	Actions  string `json:"actions"`
	Resource string `json:"resource"`
	Id       string `json:"id"`
}

// Payload represents a collection of resource access scope. It will be
// embedded into the generated API key.
type Payload struct {
	Scopes []Scope `json:"scopes"`
}
