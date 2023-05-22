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

type mass struct {
	Entries []massEntry `json:"entries"`
}

type massEntry struct {
	Name   string                 `json:"name"`
	Emails []string               `json:"emails"`
	Values map[string]interface{} `json:"values"`
}

func main() {
	massRaw, err := os.ReadFile("mass.json")
	if err != nil {
		log.Fatal("Unable to read mass.json")
	}

	var massData mass
	err = json.Unmarshal(massRaw, &massData)
	if err != nil {
		log.Println(err)
		log.Fatal("Unable to unmarshal mass data")
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

	for _, entry := range massData.Entries {
		var body bytes.Buffer
		entry.Values["RootId"] = entry.Name
		err = templateParsed.Execute(&body, templateVars{Vars: entry.Values})
		if err != nil {
			log.Println("Unable to generate template")
			log.Fatal(err)
		}

		fmt.Println(body.String())
	}

}
