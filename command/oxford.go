package command

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/gurparit/slackbot/util"
)

// OxfordDictionaryURL base Oxford Dictionary API URL
const OxfordDictionaryURL = "https://od-api.oxforddictionaries.com/api/v1/entries/en/%s"

// OxfordResponse the default response format
const OxfordResponse = "%s - %s"

const oxfordStandardResponse = "Oxford Dict.: no results found."

// OxfordDictionaryCommand dictionary command implementation
type OxfordDictionaryCommand struct {
	Etymology bool
}

// OxfordResult oxford dictionary result
type OxfordResult struct {
	Results []struct {
		LexicalEntries []struct {
			Entries []struct {
				Etymologies []string `json:"etymologies"`
				Senses      []struct {
					Definitions []string `json:"definitions"`
				} `json:"senses"`
			} `json:"entries"`
		} `json:"lexicalEntries"`
	} `json:"results"`
}

func (oxford OxfordResult) hasEtyEntry() bool {
	isValid := (len(oxford.Results) > 0 && len(oxford.Results[0].LexicalEntries) > 0 && len(oxford.Results[0].LexicalEntries[0].Entries) > 0 && len(oxford.Results[0].LexicalEntries[0].Entries[0].Etymologies) > 0)

	return isValid
}

func (oxford OxfordResult) hasDefinitionEntry() bool {
	isValid := (len(oxford.Results) > 0 && len(oxford.Results[0].LexicalEntries) > 0 && len(oxford.Results[0].LexicalEntries[0].Entries) > 0 && len(oxford.Results[0].LexicalEntries[0].Entries[0].Senses) > 0 && len(oxford.Results[0].LexicalEntries[0].Entries[0].Senses[0].Definitions) > 0)

	return isValid
}

func (oxford OxfordResult) getEty() string {
	return oxford.Results[0].LexicalEntries[0].Entries[0].Etymologies[0]
}

func (oxford OxfordResult) getDefinition() string {
	return oxford.Results[0].LexicalEntries[0].Entries[0].Senses[0].Definitions[0]
}

func (oxford OxfordDictionaryCommand) search(searchString string) (OxfordResult, error) {
	var err error

	queryString := url.QueryEscape(searchString)
	targetURL := fmt.Sprintf(OxfordDictionaryURL, queryString)

	httpCommand := HTTPCommand{URL: targetURL}
	httpCommand.Headers = make(map[string]string)
	httpCommand.Headers["Accept"] = "application/json"
	httpCommand.Headers["app_id"] = util.Config.OxfordID
	httpCommand.Headers["app_key"] = util.Config.OxfordKey

	body, err := httpCommand.Result()

	var result OxfordResult
	err = json.Unmarshal(body, &result)

	return result, err
}

// Execute OxfordDictionaryCommand implementation
func (oxford OxfordDictionaryCommand) Execute(respond func(string), query string) {
	result, err := oxford.search(query)
	if util.IsError(err) {
		respond("Oxford Dict.: time to upskill that spelling game.")
		return
	}

	resultCount := len(result.Results)

	if resultCount > 0 {
		var definition string

		if oxford.Etymology && result.hasEtyEntry() {
			definition = result.getEty()
		} else if result.hasDefinitionEntry() {
			definition = result.getDefinition()
		} else {
			definition = oxfordStandardResponse
		}

		message := fmt.Sprintf(OxfordResponse, query, definition)

		respond(message)
	} else {
		respond(oxfordStandardResponse)
	}
}
