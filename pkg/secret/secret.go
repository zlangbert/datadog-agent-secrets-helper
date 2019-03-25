package secret

// Result is the result of resolving a secret handle using a provider
type Result struct {
	Error string `json:"error,omitempty"`
	Value string `json:"value,omitempty"`
}
