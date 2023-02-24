package main

import (
	"aa-survey/internal/resources"
	"flag"
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

func main() {
	// Read cli arguments
	var surveyFile string
	flag.StringVar(&surveyFile, "survey", "./survey.yaml", "Path to survey yaml file")
	flag.Parse()

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

	// Ask the questions and print result to stdout
	answers := questionnaire.AskQuestions()
	for key, value := range answers {
		fmt.Printf("%s: %s\n", key, value)
	}
}
