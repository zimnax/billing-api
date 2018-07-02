package pub_sub

type MetadataEvent struct {
	Publisher string  `json:"Publisher"`
	EventType string  `json:"EventType"`
	Payload   Payload `json:"Payload"`
}

type Payload struct {
	Description string `json:"Description"`
	Enabled     bool   `json:"Enabled"`
	ID          string `json:"Id"`
	Links       []Link `json:"Links"`
	Name        string `json:"Name"`
	Version     int    `json:"Version"`
}

type Link struct {
	Href      string `json:"Href"`
	Rel       string `json:"Rel"`
	Templated bool   `json:"Templated"`
}
