package types

type Name struct {
	Name      string  `json:"name"`
	FirstName string  `json:"firstName"`
	LastName  *string `json:"lastName,omitempty"`
	Type      string  `json:"type"`
}

type ContactInfo struct {
	Phone     string  `json:"phone"`
	FirstName string  `json:"firstName"`
	LastName  *string `json:"lastName,omitempty"`
}

type User struct {
	ID               FlexInt                `json:"id"`
	AccountStatus    *FlexInt               `json:"accountStatus,omitempty"`
	RegistrationTime *FlexInt               `json:"registrationTime,omitempty"`
	Country          *string                `json:"country,omitempty"`
	BaseRawURL       *string                `json:"baseRawUrl,omitempty"`
	BaseURL          *string                `json:"baseUrl,omitempty"`
	Names            []Name                 `json:"names"`
	Options          []string               `json:"options"`
	PhotoID          *FlexInt               `json:"photoId,omitempty"`
	UpdateTime       *FlexInt               `json:"updateTime,omitempty"`
	Phone            *FlexInt               `json:"phone,omitempty"`
	Status           *string                `json:"status,omitempty"`
	Description      *string                `json:"description,omitempty"`
	Gender           interface{}            `json:"gender,omitempty"`
	Link             *string                `json:"link,omitempty"`
	WebApp           interface{}            `json:"webApp,omitempty"`
	MenuButton       map[string]interface{} `json:"menuButton,omitempty"`
}
