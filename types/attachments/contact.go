package attachments

type ContactAttachment struct {
	Type      string `json:"type"`
	ContactID int    `json:"contactId"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Name      string `json:"name"`
	PhotoURL  string `json:"photoUrl"`
}
