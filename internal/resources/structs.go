package resources

import (
	"aa-survey/internal/utils"
	"aa-survey/internal/validators"
	"log"

	"github.com/AlecAivazis/survey/v2"
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

func (q *Questionnaire) checkValidators() {
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
	// Check the initial validator setup.
	q.checkValidators()

	for _, question := range q.Questions {
		answer, err := askQuestion(question)
		if err != nil {
			handleQuestionError(err)
		}
		answers[question.Target] = answer
	}

	return answers
}
