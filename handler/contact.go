package handler

// Contact is the contact struct
type Contact struct {
	ID                string   `json:"contact_id,omitempty"`
	Email             string   `json:"Email,omitempty"`
	Twitter           string   `json:"Twitter,omitempty"`
	FirstName         string   `json:"FirstName,omitempty"`
	LastName          string   `json:"LastName,omitempty"`
	Salutation        string   `json:"Salutation,omitempty"`
	Company           string   `json:"Company,omitempty"`
	NumberOfEmployees int      `json:"NumberOfEmployees,omitempty"`
	Title             string   `json:"Title,omitempty"`
	Type              string   `json:"type,omitempty"`
	Industry          string   `json:"Industry,omitempty"`
	Phone             string   `json:"Phone,omitempty"`
	MobilePhone       string   `json:"MobilePhone,omitempty"`
	Fax               string   `json:"Fax,omitempty"`
	Website           string   `json:"Website,omitempty"`
	MailingStreet     string   `json:"MailingStreet,omitempty"`
	MailingCity       string   `json:"MailingCity,omitempty"`
	MailingState      string   `json:"MailingState,omitempty"`
	MailingPostalCode string   `json:"MailingPostalCode,omitempty"`
	MailingCountry    string   `json:"MailingCountry,omitempty"`
	OwnerName         string   `json:"owner_name,omitempty"`
	LeadSource        string   `json:"LeadSource,omitempty"`
	Status            string   `json:"Status,omitempty"`
	LinkedIn          string   `json:"LinkedIn,omitempty"`
	Unsubscribed      bool     `json:"unsubscribed,omitempty"`
	Custom            string   `json:"custom,omitempty"`
	SessionID         string   `json:"_autopilot_session_id,omitempty"`
	List              string   `json:"_autopilot_list,omitempty"`
	Notify            bool     `json:"notify,omitempty"`
	Lists             []string `json:"lists,omitempty"`
}
