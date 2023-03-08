package types

type Payload struct {
	RequestID string `json:"requestID"`
	Method    string `json:"method"`
	Data      *Data  `json:"data"`
}

type Data struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	PartNumber int    `json:"partNumber"`
	Quantity   int    `json:"quantity"`
}
