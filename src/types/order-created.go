package types

type OrderCreatedEvent struct {
	SpecVersion     string `json:"specversion"`
	Type            string `json:"type"`
	Source          string `json:"source"`
	Subject         string `json:"subject"`
	Id              string `json:"id"`
	Time            string `json:"time"`
	DataContentType string `json:"datacontenttype"`
	Data            Order  `json:"data"`
}
