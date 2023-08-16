package domain

type Account struct {
	Base

	// account fields
	Name      string  `json:"name"`
	Balance   float64 `json:"balance"`
	Note      string  `json:"note,omitempty"`
	CreatedBy uint    `json:"-"`
}
