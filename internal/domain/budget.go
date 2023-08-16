package domain

type Budget struct {
	Base

	// budget fields
	Amount     float64 `json:"amount"`
	CategoryID uint    `json:"-"`
	CreatedBy  uint    `json:"-"`
}
