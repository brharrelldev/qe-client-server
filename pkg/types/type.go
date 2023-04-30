package types

type Payload struct {
	Id     string `json:"id"`
	Method string `json:"method"`
	Data   *Data  `json:"data"`
}

type Data struct {
	Name       string `json:"name"`
	PartNumber int    `json:"partNumber"`
	Quantity   int    `json:"quantity"`
}

type Response struct {
	Message string
	Results []Payload
}
