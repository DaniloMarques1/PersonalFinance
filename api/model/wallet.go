package model

type Wallet struct {
	Id          int64   `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Total       float64 `json:"total"`
	Client_id   int64   `json:"client_id"`
}

// TODO add removeWallet
type IWallet interface {
	SaveWallet(wallet *Wallet) error
	RemoveWallet(client_id, wallet_id int64) error
	GetWallets(client_id int64) ([]Wallet, error)
	//GetWallet(wallet_id, client_id int64) (*Wallet, error) // TODO will have to return an slice of movements
}
