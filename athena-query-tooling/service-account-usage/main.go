package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
	orgTypes "github.com/aws/aws-sdk-go-v2/service/organizations/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/athena"
)

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("eu-west-1"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	svc := organizations.NewFromConfig(cfg)

	var maxResults int32 = 20
	var accounts []orgTypes.Account

	var Ids = make([]string, 1)

	input := organizations.ListAccountsInput{MaxResults: &maxResults}

	resps := organizations.NewListAccountsPaginator(svc, &input)
	for resps.HasMorePages() {
		page, err := resps.NextPage(context.TODO())
		if err != nil {
			log.Fatal(err)
		}

		accounts = append(accounts, page.Accounts...)
	}

	for _, acc := range accounts {
		fmt.Println(*acc.Id)
		result, err := runAthenaQuery(*acc.Id)
		if err != nil {
			fmt.Printf("Error running Athena Query for account %v", *acc.Id)
		}
		Ids = append(Ids, result)
		time.Sleep(time.Duration(10) * time.Second) // Try not to hit rate limiting
	}

	fmt.Println(Ids)

	// Cycle through exceution ID's
	// Check for status is SUCCEEDED
	// Grab .csv from bucket and place it into a folder
	// Open each CSV and pull unique service account and append to list
	// Display a list of unique service accounts in use (these contain ARN with account)

	for _, f := range Ids {
		file, err := DownloadFile("dfds-audit", f+".csv", f+".csv")
		if err != nil {
			fmt.Println("Error downloading file " + f)
		}
		fmt.Println(file + ".csv")

	}

}

func DownloadFile(bucketName string, objectKey string, fileName string) (file string, err error) {

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("eu-central-1"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	svc := s3.NewFromConfig(cfg)

	result, err := svc.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		log.Printf("Couldn't get object %v:%v. Here's why: %v\n", bucketName, objectKey, err)
		return "", err
	}
	defer result.Body.Close()

	var f *os.File
	f, err = os.Create(fileName)
	if err != nil {
		log.Printf("Couldn't create file %v. Here's why: %v\n", fileName, err)
		return "", err
	}
	defer f.Close()
	body, err := io.ReadAll(result.Body)
	if err != nil {
		log.Printf("Couldn't read object body from %v. Here's why: %v\n", objectKey, err)
	}
	_, err = f.Write(body)
	return f.Name(), err
}

func runAthenaQuery(acc string) (id string, err error) {

	awscfg := &aws.Config{}
	awscfg.WithRegion("eu-west-1")

	sess := session.Must(session.NewSession())

	svc := athena.New(sess, aws.NewConfig().WithRegion("eu-central-1"))
	var s athena.StartQueryExecutionInput
	querystring := fmt.Sprintf("SELECT useridentity.arn, eventtime FROM default.cloudtrail WHERE accountid = '%s' AND eventTime > to_iso8601(current_timestamp - interval '7' day) AND useridentity.arn LIKE '%%/s.%%'", acc)
	s.SetQueryString(querystring)
	fmt.Println(querystring)

	var q athena.QueryExecutionContext
	q.SetDatabase("default")
	s.SetQueryExecutionContext(&q)

	var r athena.ResultConfiguration
	r.SetOutputLocation("s3://dfds-audit/")
	s.SetResultConfiguration(&r)

	result, err := svc.StartQueryExecution(&s)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	fmt.Println("StartQueryExecution result:")
	fmt.Println(result.GoString())

	var qri athena.GetQueryExecutionInput
	qri.SetQueryExecutionId(*result.QueryExecutionId)

	var qrop *athena.GetQueryExecutionOutput
	duration := time.Duration(5) * time.Second // Pause for 5 seconds

	for {
		qrop, err = svc.GetQueryExecution(&qri)
		if err != nil {
			fmt.Println(err)
			return "", err
		}
		if *qrop.QueryExecution.Status.State != "RUNNING" {
			break
		}
		fmt.Println("waiting.")
		time.Sleep(duration)

	}
	if *qrop.QueryExecution.Status.State == "SUCCEEDED" {

		var ip athena.GetQueryResultsInput
		ip.SetQueryExecutionId(*result.QueryExecutionId)

		op, err := svc.GetQueryResults(&ip)
		if err != nil {
			fmt.Println(err)
			return "", err
		}
		fmt.Printf("%+v", op)
	} else {
		fmt.Println(*qrop.QueryExecution.Status.State)

	}

	return *result.QueryExecutionId, nil
}
