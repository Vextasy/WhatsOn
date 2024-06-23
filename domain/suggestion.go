package domain

type Suggestion struct {
	Channel     string `json:"channel" xml:"Channel"` // BBC1, BBC2, BBC3 ...
	Name        string `json:"name" xml:"Name"`       // The programme or event name
	Date        string `json:"date" xml:"Date"`       // In YYYY-MM-DD format.
	Time        string `json:"time" xml:"Time"`       // In HH:MM format.
	Description string `json:"description" xml:"Description"`
}
