package client

type ClientRepr struct {
	Id       int64  `json:"id"`
	FullName string `json:"fullName"`
}

type CreateClient struct {
	FullName string `json:"fullName"`
}
