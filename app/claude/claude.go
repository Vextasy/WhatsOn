package claude

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/liushuangls/go-anthropic/v2"
	"github.com/liushuangls/go-anthropic/v2/jsonschema"
	"gitlab.com/vextasy/claude/whatson/domain"
)

type suggestionSvc struct {
	verbose bool
	tvDbSvc domain.TvDbSvc
}

func NewSuggestionSvc(verbose bool, tvDbSvc domain.TvDbSvc) domain.SuggestionSvc {
	return suggestionSvc{verbose: verbose, tvDbSvc: tvDbSvc}
}

func (c suggestionSvc) GetSuggestions(ctx context.Context, desire string) ([]domain.Suggestion, error) {
	return getSuggestions(ctx, desire, c.verbose, c.tvDbSvc)
}

func getSuggestions(ctx context.Context, desire string, verbose bool, tvDbSvc domain.TvDbSvc) ([]domain.Suggestion, error) {
	client := anthropic.NewClient(os.Getenv("ANTHROPIC_API_KEY"))
	request := anthropic.MessagesRequest{
		Model:     anthropic.ModelClaude3Dot5Sonnet20240620,
		MaxTokens: 1024,
		ToolChoice: &anthropic.ToolChoice{
			Type: "auto",
		},
		Tools: []anthropic.ToolDefinition{
			{
				Name: "get_tv_programmes",
				Description: `Gets information about the TV programmes that will be aired between the given from_date and to_date inclusive.
					If to_date is not provided, it will default to the same date as the from_date.
					It should be used when the user asks about TV programmes and we want to know which programs will be shown in the near future including today.
					Programme information will be returned in XML tags with the following format:
						<Programmes>
							<Programme>
								<Channel>Channel Name</Channel>
								<Name>Programme Name</Name>
								<Date>YYYY-MM-DD</Date>
								<Time>HH:MM</Time>
								<Description>Programme Description</Description>
							</Programme>
						</Programmes>
					The outer level <Programmes></Programmes> tag will be omitted if the result is empty,
					but will otherwise contain one or more <Programme></Programme> tags.
					The <Location></Location> tag contains the TV channel on which the program will be shown.
					`,
				InputSchema: jsonschema.Definition{
					Type: jsonschema.Object,
					Properties: map[string]jsonschema.Definition{
						"from_date": {
							Type:        jsonschema.String,
							Description: "The start date of the TV programmes search. Format: YYYY-MM-DD",
						},
						"to_date": {
							Type:        jsonschema.String,
							Description: "The end date of the TV programmes search. Format: YYYY-MM-DD",
						},
					},
					Required: []string{"from_date"},
				},
			}},
		Messages: []anthropic.Message{
			anthropic.NewUserTextMessage(`
			You are given information about TV programmes enclosed between <Programmes></Programmes> tags.
			The date today is ` + time.Now().Format("2006-01-02") + `, a ` + time.Now().Weekday().String() + `.
			You are an expert at suggesting TV programmes.
			The user would like to have some suggestions about what TV programmes to watch but has the following strict request:<Request>` + desire + `</Request>.
			Your task is to provide a list of programme suggestions in XML format using the following tag structure for each selection:
			<Suggestion><Channel>Channel Name</Channel><Name>Programme Name</Name <Date>YYYY-MM-DD</Date><Time>HH:MM</Time><Description>Programme Description</Description></Suggestion>
			Wrap the list of <Suggestion></Suggestion> tags in a <Suggestions></Suggestions> tag.
			Do not output any other information either before or after the <Suggestions></Suggestions> tags.
			If the user does not request that programs begin within a particular date range then use a range that begins today and ends 1 week later.
			Return at most 10 suggestions.
			If you have no suggestions then return the <Suggestions></Suggestions> tags with no content.
			Before answering re-check that all of the user's requests have been met.
			If you are not sure that a particular suggestion satisfies the user's Requests then do not return that suggestion.
			If the user has requested programmes on a particular day of the week or days of the week check that only programmes on those days are returned.
			`),
		},
	}
	// Make the initial request.
	resp, err := client.CreateMessages(ctx, request)
	if err != nil {
		var e *anthropic.APIError
		if errors.As(err, &e) {
			fmt.Printf("Messages error, type: %s, message: %s", e.Type, e.Message)
		} else {
			fmt.Printf("Messages error: %v\n", err)
		}
		return nil, err
	}
	request.Messages = append(request.Messages, anthropic.Message{
		Role:    anthropic.RoleAssistant,
		Content: resp.Content,
	})

	// Look for the tool use response.
	var toolUse *anthropic.MessageContentToolUse
	for _, content := range resp.Content {
		if content.Type == anthropic.MessagesContentTypeToolUse {
			toolUse = content.MessageContentToolUse
			break
		} else {
			if verbose {
				fmt.Printf("\n%s\n", *content.Text)
			}
		}
	}
	xmlUp := func(xmlText string) ([]domain.Suggestion, error) {
		container := SuggestionsContainer{}
		err := xml.Unmarshal([]byte(xmlText), &container)
		if err != nil {
			return nil, err
		}
		return container.Suggestions, nil
	}
	if toolUse == nil {
		// Perhaps Claude already has some suggestions.
		if resp.Content[0].Text != nil {
			return xmlUp(*resp.Content[0].Text)
		}
		// Maybe not.
		return []domain.Suggestion{}, nil
	}

	// Extract the tool use input parameters.
	type GetTvProgrammesInput struct {
		FromDate string `json:"from_date"`
		ToDate   string `json:"to_date"`
	}
	input := GetTvProgrammesInput{}
	err = json.Unmarshal([]byte(toolUse.Input), &input)
	if err != nil {
		return nil, err
	}
	if input.FromDate == "" {
		todayString := time.Now().Format("2006-01-02")
		input.FromDate = todayString
	}
	if input.ToDate == "" {
		input.ToDate = input.FromDate
	}
	if verbose {
		fmt.Printf("\nUsing the tool date range: %s to %s\n", input.FromDate, input.ToDate)
	}

	// Call the tool to get TV programmes.
	toolResults, nprogs, err := tvDbSvc.GetTvProgrammesXml(ctx, input.FromDate, input.ToDate)
	if nprogs == 0 {
		return []domain.Suggestion{}, nil
	}
	// Enqueue the tool results message.
	request.Messages = append(request.Messages, anthropic.NewToolResultsMessage(toolUse.ID, string(toolResults), false))
	// Call the Agent with the tool results.
	resp, err = client.CreateMessages(ctx, request)
	if err != nil {
		var e *anthropic.APIError
		if errors.As(err, &e) {
			fmt.Printf("Messages error, type: %s, message: %s", e.Type, e.Message)
		} else {
			fmt.Printf("Messages error: %v\n", err)
		}
		return nil, err
	}
	//fmt.Printf("Response: %+v\n", resp)

	if len(resp.Content) == 0 || resp.Content[0].Text == nil {
		return []domain.Suggestion{}, nil
	}
	return xmlUp(*resp.Content[0].Text)
}

type SuggestionsContainer struct {
	Suggestions []domain.Suggestion `xml:"Suggestion"`
}
