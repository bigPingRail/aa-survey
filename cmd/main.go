package main

import (
	"aa-survey/internal/resources"
	"aa-survey/internal/utils"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"gopkg.in/yaml.v2"
)

func main() {
	// Read cli arguments
	var surveyFile, outFile string
	var postCheck bool
	flag.StringVar(&surveyFile, "survey", "Surveyfile.yaml", "Path to survey yaml file")
	flag.StringVar(&outFile, "output", "", "Path to output file. If no output file provided, result will be written to STDOUT")
	flag.BoolVar(&postCheck, "check", false, "Additionally prompts the user to check the correctness of entered data")
	flag.Parse()

	// Check file extension and create file if is not exists
	if outFile != "" {
		utils.CheckFileExt(outFile)
	}

	// Read the YAML file with questions
	data, err := os.ReadFile(surveyFile)
	if err != nil {
		log.Fatalf("error reading questions file: %v", err)
	}

	// Unmarshal the YAML data into a Questionnaire struct
	var questionnaire resources.Questionnaire
	err = yaml.Unmarshal(data, &questionnaire)
	if err != nil {
		log.Fatalf("error unmarshaling questions: %v", err)
	}

	// Ask the questions and print result to stdout or file
	answers := questionnaire.AskQuestions()
	if outFile == "" {
		for key, value := range answers {
			fmt.Printf("%s: %s\n", key, value)
		}
	} else {
		utils.CreateIfNotExists(outFile)
		utils.WriteToFile(outFile, answers)
		if postCheck {
			var ans bool
			p := &survey.Confirm{
				Message: "All Correct?",
				Default: true,
			}
			err := survey.AskOne(p, &ans)
			if err != nil {
				log.Fatalf("%v", err)
			}
			if !ans {
				utils.RunEditor(outFile)
			}
		}
	}
}
