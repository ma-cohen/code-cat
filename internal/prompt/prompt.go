package prompt

import (
	"github.com/AlecAivazis/survey/v2"
)

// AskString prompts the user for a non-empty string.
func AskString(message, defaultVal string) (string, error) {
	var result string
	q := &survey.Input{
		Message: message,
		Default: defaultVal,
	}
	err := survey.AskOne(q, &result, survey.WithValidator(survey.Required))
	return result, err
}

// AskSelect prompts the user to pick one option from a list.
func AskSelect(message string, options []string) (string, error) {
	var result string
	q := &survey.Select{Message: message, Options: options}
	err := survey.AskOne(q, &result)
	return result, err
}

// AskConfirm prompts the user for a yes/no answer.
func AskConfirm(message string, defaultVal bool) (bool, error) {
	var result bool
	q := &survey.Confirm{
		Message: message,
		Default: defaultVal,
	}
	err := survey.AskOne(q, &result)
	return result, err
}
