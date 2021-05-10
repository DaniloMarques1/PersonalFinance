package model

type Client struct {
	Id           int64   `json:"id,omitempty"`
	Name         string  `json:"name,omitempty"`
	Email        string  `json:"email,omitempty"`
	Total        float64 `json:"total,omitempty"`
	PasswordHash []byte  `json:"password_hash,omitempty"`
}

type IClient interface {
	SaveClient(client *Client) error
	FindById(id int64) (*Client, error)
	FindByEmail(email string) (*Client, error)
}
