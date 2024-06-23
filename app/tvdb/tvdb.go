package tvdb

import (
	"context"
	"encoding/xml"
	"log"
	"os"
	"time"

	"gitlab.com/vextasy/claude/whatson/domain"
)

type dbProgrammes struct {
	XMLName    xml.Name      `xml:"Programmes"`
	Programmes []dbProgramme `xml:"Programme"`
}
type dbProgramme struct {
	Channel     string `xml:"Channel"`
	Name        string `xml:"Name"`
	Date        string `xml:"Date"`
	Time        string `xml:"Time"`
	Description string `xml:"Description"`
}

type tvDbSvc struct {
	pathToTvDb string
}

func NewTvDbSvc(pathToTvDb string) domain.TvDbSvc {
	return tvDbSvc{pathToTvDb: pathToTvDb}
}

// Return TV programs that are aired between the given dates as an XML string
func (svc tvDbSvc) GetTvProgrammesXml(ctx context.Context, dateFrom string, dateTo string) (string, error) {
	programmes, err := getTvProgrammes(svc.pathToTvDb, dateFrom, dateTo)
	if err != nil {
		return "", err
	}
	xmlbytes, err := xml.MarshalIndent(programmes, "", "  ")
	if err != nil {
		return "", err
	}
	return string(xmlbytes), err
}

// Get programmes from the tvdb xml file and return as a slice of TvProgramme.
// To avoid having to continually rebuild the database of TV programmes we will
// assume that a date in the database of 2024-06-21 represents today.
func getTvProgrammes(dbpath string, dateFrom string, dateTo string) ([]domain.TvProgramme, error) {
	xmlFile, err := os.ReadFile(dbpath)
	if err != nil {
		log.Fatalf("Error reading XML file: %v", err)
	}

	// Create a Programmes struct to unmarshal the XML into
	var xprogrammes dbProgrammes

	// Unmarshal the XML into the programmes struct
	err = xml.Unmarshal(xmlFile, &xprogrammes)
	if err != nil {
		log.Fatalf("Error unmarshalling XML: %v", err)
	}

	tvProgrammes := []domain.TvProgramme{}

	for _, prog := range xprogrammes.Programmes {
		prog.Date = shiftDate(prog.Date)
		if prog.Date > dateTo || prog.Date < dateFrom {
			continue
		}
		tvProgrammes = append(tvProgrammes, domain.TvProgramme{
			Channel:     prog.Channel,
			Name:        prog.Name,
			Date:        prog.Date,
			Time:        prog.Time,
			Description: prog.Description,
		})
	}
	return tvProgrammes, nil
}

// Shift the date string to be bigger by the difference
// in days between 2024-06-21 and today's date.
func shiftDate(date string) string {
	// calculate the number of days difference between 2024-06-21 and today.
	days := int(time.Now().AddDate(0, 0, 0).Sub(time.Date(2024, 6, 21, 0, 0, 0, 0, time.UTC)).Hours() / 24)
	t, _ := time.Parse("2006-01-02", date)
	return t.AddDate(0, 0, days).Format("2006-01-02")
}
