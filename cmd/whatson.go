package main

import (
	"context"
	"flag"
	"fmt"
	"strings"
	"time"

	"gitlab.com/vextasy/claude/whatson/app/claude"
	"gitlab.com/vextasy/claude/whatson/app/tvdb"
)

func main() {
	var desire = "I'm feeling lucky."
	var descriptions = flag.Bool("d", false, "Show descriptions")
	flag.Parse()
	args := flag.Args()
	if len(args) > 0 {
		desire = strings.Join(args, " ")
	}
	tvDbServices := tvdb.NewServices("./TvDb.xml")
	claude := claude.NewServices(tvDbServices.TvDb)
	suggestions, err := claude.Suggestion.GetSuggestions(context.Background(), desire)
	if err != nil {
		fmt.Println(err)
	}
	if len(suggestions) == 0 {
		fmt.Println("Sorry, I couldn't find any suggestions for that.")
		return
	}
	for _, s := range suggestions {
		fmt.Println("- "+weekday(s.Date)[0:3], s.Date, "at", s.Time+": "+s.Name, "on", s.Channel)
	}
	if *descriptions {
		fmt.Println("\nDescriptions:")
		for _, s := range suggestions {
			fmt.Println("- "+s.Name, ":", s.Description)
		}
	}
}

// Return the weekday for the given date string argument.
// The date string should be in the format "YYYY-MM-DD".
func weekday(dateStr string) string {
	d, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return "Unk"
	}
	return d.Weekday().String()
}
