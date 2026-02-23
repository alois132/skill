package resources

type Reference struct {
	Name string `json:"name"`
	Body string `json:"body"`
}

// String returns the reference content
func (r *Reference) String() string {
	return r.Body
}

// Summary returns a brief summary of the reference
func (r *Reference) Summary() string {
	return r.Body
}
