package domain

type Category struct {
	Base

	// category fields
	Name      string `json:"name"`
	Note      string `json:"note,omitempty"`
	CreatedBy *uint  `json:"-"`
}
