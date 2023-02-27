package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
	orgTypes "github.com/aws/aws-sdk-go-v2/service/organizations/types"
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
		runAthenaQuery(*acc.Id)
		time.Sleep(time.Duration(10) * time.Second) // Try not to hit rate limiting
	}

}

func runAthenaQuery(acc string) {

	awscfg := &aws.Config{}
	awscfg.WithRegion("eu-west-1")

	sess := session.Must(session.NewSession())

	svc := athena.New(sess, aws.NewConfig().WithRegion("eu-central-1"))
	var s athena.StartQueryExecutionInput
	querystring := fmt.Sprintf("SELECT useridentity.arn, eventtime FROM default.cloudtrail WHERE accountid = '%s' AND eventTime > to_iso8601(current_timestamp - interval '30' day) AND useridentity.arn LIKE '%%/s.%%'", acc)
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
		return
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
			return
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
			return
		}
		fmt.Printf("%+v", op)
	} else {
		fmt.Println(*qrop.QueryExecution.Status.State)

	}
}
