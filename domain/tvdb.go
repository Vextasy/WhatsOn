package domain

type TvProgramme struct {
	Channel     string `json:"channel" xml:"Channel"`
	Name        string `json:"name" xml:"Name"`
	Date        string `json:"date" xml:"Date"`
	Time        string `json:"time" xml:"Time"`
	Description string `json:"description" xml:"Description"`
}
