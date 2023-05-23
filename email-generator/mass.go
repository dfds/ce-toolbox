package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	awsHttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
	"log"
	"net/http"
	"os"
	"text/template"
)

type templateVars struct {
	Vars map[string]interface{}
}

type mass struct {
	Title   string      `json:"title"`
	Entries []massEntry `json:"entries"`
}

type massEntry struct {
	Name   string                 `json:"name"`
	Emails []string               `json:"emails"`
	Values map[string]interface{} `json:"values"`
}

type sesRequest struct {
	Msg    string
	Title  string
	From   string
	Emails []string
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
		err = sendEmail(context.Background(), sesRequest{
			Msg:    body.String(),
			Title:  massData.Title,
			From:   "noreply@dfds.cloud",
			Emails: entry.Emails,
		})
		if err != nil {
			log.Fatal(err)
		}
	}

}

func sendEmail(ctx context.Context, req sesRequest) error {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("eu-west-1"), config.WithHTTPClient(CreateHttpClientWithoutKeepAlive()))
	if err != nil {
		return err
	}

	sesClient := sesv2.NewFromConfig(cfg)

	input := &sesv2.SendEmailInput{
		FromEmailAddress: &req.From,
		Destination:      &types.Destination{BccAddresses: req.Emails},
		Content: &types.EmailContent{
			Simple: &types.Message{
				Body: &types.Body{Text: &types.Content{
					Data: &req.Msg,
				}},
				Subject: &types.Content{
					Data: &req.Title,
				},
			},
		},
	}

	output, err := sesClient.SendEmail(ctx, input)
	if err != nil {
		fmt.Println(output)
		return err
	}

	fmt.Println(output)

	return nil
}

func CreateHttpClientWithoutKeepAlive() *awsHttp.BuildableClient {
	client := awsHttp.NewBuildableClient().WithTransportOptions(func(transport *http.Transport) {
		transport.DisableKeepAlives = true
	})

	return client
}
