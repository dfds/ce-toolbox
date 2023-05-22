package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"text/template"
)

type templateVars struct {
	Vars map[string]interface{}
}

func main() {
	varsRaw, err := os.ReadFile("vars.json")
	if err != nil {
		log.Fatal("Unable to read vars.json")
	}

	var vars map[string]interface{}
	err = json.Unmarshal(varsRaw, &vars)
	if err != nil {
		log.Fatal("Unable to unmarsharl vars")
	}

	templateRaw, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal("Unable to read template file")
	}

	templateContainer := template.New("gen")
	templateParsed, err := templateContainer.Parse(string(templateRaw))
	if err != nil {
		log.Fatal("Unable to parse template file")
	}

	var body bytes.Buffer
	err = templateParsed.Execute(&body, templateVars{Vars: vars})
	if err != nil {
		log.Println("Unable to generate template")
		log.Fatal(err)
	}

	fmt.Println(body.String())
}
