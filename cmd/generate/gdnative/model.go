package gdnative

// API is a single API definition in `gdnative_api.json`
type API struct {
	Name    string `json:"name,omitempty"`
	Type    string `json:"type"`
	Version struct {
		Major int `json:"major"`
		Minor int `json:"minor"`
	} `json:"version"`
	API []struct {
		Name       string     `json:"name"`
		ReturnType string     `json:"return_type"`
		Arguments  [][]string `json:"arguments"`
	} `json:"api"`
}

// APIs is a structure based on `gdnative_api.json` in `godot_headers`.
type APIs struct {
	Core       API   `json:"core"`
	Extensions []API `json:"extensions"`
	// Extensions map[string]API `json:"extensions"`
}
