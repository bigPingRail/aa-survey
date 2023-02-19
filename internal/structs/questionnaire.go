package structs

import (
	"aa-survey-v2/internal/utils"
	"aa-survey-v2/internal/validators"
	"log"
	"path/filepath"
	"strconv"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
)

type Question struct {
	Prompt  string   `yaml:"prompt"`
	Type    string   `yaml:"type"`
	Default string   `yaml:"default"`
	Target  string   `yaml:"target"`
	Answer  string   `yaml:"answer"`
	Options []string `yaml:"options,omitempty"`

	Validate string                      `yaml:"validate"`
	VFunc    func(ans interface{}) error `yaml:"-"`
}

type Questionnaire struct {
	Questions []Question `yaml:"questions"`
}

func (q *Questionnaire) SetupValidators() {
	validatorsMap := map[string]struct {
		ValidatorFunc func(interface{}) error
		AllowedTypes  []string
	}{
		"password": {
			ValidatorFunc: validators.ValidatePassword,
			AllowedTypes:  []string{"password"},
		},
		"many": {
			ValidatorFunc: validators.ValidateMany,
			AllowedTypes:  []string{"multiselect"},
		},
		"required": {
			ValidatorFunc: survey.Required,
			AllowedTypes:  []string{},
		},
	}

	for i := range q.Questions {
		question := &q.Questions[i]
		if v, ok := validatorsMap[question.Validate]; ok {
			if len(v.AllowedTypes) == 0 || utils.Contains(v.AllowedTypes, question.Type) {
				question.VFunc = v.ValidatorFunc
			} else {
				log.Fatalf("Validator type: \"%s\" is not allowed for \"%s\" field", question.Validate, question.Type)
			}
		}
	}
}

func (q *Questionnaire) AskQuestions() map[string]interface{} {
	answers := make(map[string]interface{})
	for _, question := range q.Questions {
		var err error
		var ans string

		switch question.Type {
		case "confirm":
			var ans bool
			prompt := &survey.Confirm{
				Message: question.Prompt,
				Default: utils.ConvertBootToStr(question.Default),
			}
			err = survey.AskOne(prompt, &ans)
			answers[question.Target] = strconv.FormatBool(ans)

		case "input":
			prompt := &survey.Input{
				Message: question.Prompt,
				Default: question.Default,
			}
			err = utils.AskQuestion(prompt, &ans, question.VFunc)
			answers[question.Target] = ans

		case "password":
			prompt := &survey.Password{
				Message: question.Prompt,
			}
			err = utils.AskQuestion(prompt, &ans, question.VFunc)
			answers[question.Target] = ans

		case "public_key", "private_key", "file":
			var cv survey.Validator
			switch question.Type {
			case "public_key":
				cv = survey.ComposeValidators(survey.Required, validators.ValidatePubKey)
			case "private_key":
				cv = survey.ComposeValidators(survey.Required, validators.ValidatePrivKey)
			case "file":
				cv = survey.ComposeValidators(survey.Required, validators.ValidateIsFile)
			}

			prompt := &survey.Input{
				Message: question.Prompt,
				Suggest: func(toComplete string) []string {
					files, _ := filepath.Glob(toComplete + "*")
					return files
				},
			}
			err = survey.AskOne(prompt, &ans, survey.WithValidator(cv))
			answers[question.Target], _ = filepath.Abs(ans)

		case "select":
			prompt := &survey.Select{
				Message: question.Prompt,
				Options: question.Options,
			}
			err = utils.AskQuestion(prompt, &ans, question.VFunc)
			answers[question.Target] = ans

		case "multiselect":
			var ans []string

			prompt := &survey.MultiSelect{
				Message: question.Prompt,
				Options: question.Options,
			}
			err = utils.AskQuestion(prompt, &ans, question.VFunc)
			answers[question.Target] = ans
		}

		if err != nil {
			// Handle crtl+c
			if err == terminal.InterruptErr {
				log.Fatal("Survey Interrupted...")
			}
			log.Fatalf("error asking question: %v\n", err)
		}
	}
	return answers
}
