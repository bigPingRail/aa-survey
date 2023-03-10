package resources

import (
	"aa-survey/internal/utils"
	"aa-survey/internal/validators"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
)

func promptWithValidator(prompt survey.Prompt, answer interface{}, validatorFunc func(interface{}) error) error {
	if validatorFunc == nil {
		return survey.AskOne(prompt, answer)
	}
	return survey.AskOne(prompt, answer, survey.WithValidator(validatorFunc))
}

func handleQuestionError(err error) {
	if err == terminal.InterruptErr {
		log.Fatal("Survey Interrupted...")
	}
	log.Fatalf("error asking question: %v\n", err)
}

func globHomeDir(s string) string {
	re := regexp.MustCompile(`^~`)
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return s
	}
	return re.ReplaceAllString(s, homeDir)
}

func askConfirm(q Question) (string, error) {
	prompt := &survey.Confirm{
		Message: q.Prompt,
		Default: utils.ConvertBoolToStr(q.Default),
		Help:    q.Help,
	}
	var answer bool
	err := survey.AskOne(prompt, &answer)
	return strconv.FormatBool(answer), err
}

func askInput(q Question) (string, error) {
	prompt := &survey.Input{
		Message: q.Prompt,
		Default: q.Default,
		Help:    q.Help,
	}
	var answer string
	err := promptWithValidator(prompt, &answer, q.VFunc)
	return answer, err
}

func askPassword(q Question) (string, error) {
	prompt := &survey.Password{
		Message: q.Prompt,
		Help:    q.Help,
	}
	var answer string
	err := promptWithValidator(prompt, &answer, q.VFunc)
	return answer, err
}

func askFile(q Question) (string, error) {
	var cv survey.Validator

	switch q.Type {
	case "public_key":
		cv = survey.ComposeValidators(survey.Required, validators.ValidatePubKey)
	case "private_key":
		cv = survey.ComposeValidators(survey.Required, validators.ValidatePrivKey)
	case "file":
		cv = survey.ComposeValidators(survey.Required, validators.ValidateIsFile)
	case "dir":
		cv = survey.ComposeValidators(survey.Required, validators.ValidateIsDir)
	}

	prompt := &survey.Input{
		Message: q.Prompt,
		Suggest: func(toComplete string) []string {
			files, _ := filepath.Glob(globHomeDir(toComplete) + "*")
			return files
		},
		Help: q.Help,
	}
	var answer string
	err := survey.AskOne(prompt, &answer, survey.WithValidator(cv))
	path, _ := filepath.Abs(answer)
	return path, err
}

func askSelect(q Question) (string, error) {
	prompt := &survey.Select{
		Message: q.Prompt,
		Options: q.Options,
		Help:    q.Help,
	}
	var answer string
	err := promptWithValidator(prompt, &answer, q.VFunc)
	return answer, err
}

func askMultiSelect(q Question) ([]string, error) {
	prompt := &survey.MultiSelect{
		Message: q.Prompt,
		Options: q.Options,
		Help:    q.Help,
	}
	var answer []string
	err := promptWithValidator(prompt, &answer, q.VFunc)
	return answer, err
}

func askQuestion(question Question) (interface{}, error) {
	switch question.Type {
	case "confirm":
		return askConfirm(question)
	case "input":
		return askInput(question)
	case "password":
		return askPassword(question)
	case "public_key", "private_key", "file", "dir":
		return askFile(question)
	case "select":
		return askSelect(question)
	case "multiselect":
		return askMultiSelect(question)
	default:
		return nil, fmt.Errorf("unsupported question type: %s", question.Type)
	}
}
