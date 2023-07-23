package api

type AddProject struct {
	ID      string   `json:"ID,omitempty"`
	Name    string   `json:"Name,omitempty"`
	Members []string `json:"Members,omitempty"`
}
