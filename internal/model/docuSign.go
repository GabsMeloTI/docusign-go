package model

type EnvelopeDefinition struct {
	EmailSubject string     `json:"emailSubject"`
	Documents    []Document `json:"documents"`
	Recipients   Recipients `json:"recipients"`
	Status       string     `json:"status"`
}

type Document struct {
	DocumentBase64 string `json:"documentBase64"`
	Name           string `json:"name"`
	FileExtension  string `json:"fileExtension"`
	DocumentID     string `json:"documentId"`
}

type Recipients struct {
	Signers []Signer `json:"signers"`
}

type Signer struct {
	Email        string `json:"email"`
	Name         string `json:"name"`
	RecipientID  string `json:"recipientId"`
	RoutingOrder string `json:"routingOrder"`
	Tabs         Tabs   `json:"tabs"`
}

type Tabs struct {
	SignHereTabs []SignHereTab `json:"signHereTabs"`
}

type SignHereTab struct {
	DocumentID  string `json:"documentId"`
	PageNumber  string `json:"pageNumber"`
	RecipientID string `json:"recipientId"`
	TabLabel    string `json:"tabLabel"`
	XPosition   string `json:"xPosition"`
	YPosition   string `json:"yPosition"`
}
