package hmapi

type Link struct {
	Href     string    `json:"href,omitempty"`
	Type     MediaType `json:"type,omitempty"`
	Encoding MediaType `json:"encoding,omitempty"`
}

type LinkRequest interface{}
