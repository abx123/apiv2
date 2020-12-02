package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/db"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"google.golang.org/api/option"
)

var client *db.Client

// Novel ...
type Novel struct {
	Name string
}

func main() {
	fmt.Println("lambda start")
	Init()
	lambda.Start(handleRequest)
}

// The input type and the output type are defined by the API Gateway.
func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println(fmt.Sprintf("qs:%+v", request.QueryStringParameters))
	var err error
	var resp interface{}
	headers := map[string]string{
		"Access-Control-Allow-Headers": "Content-Type",
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods": "GET",
	}

	resp, err = getNovelList()

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers:    headers,
			Body:       err.Error(),
		}, err
	}
	formattedResp := formatResp(resp)

	response := events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers:    headers,
		Body:       formattedResp,
	}
	return response, nil
}

// Init ...
func Init() {
	ctx := context.Background()
	conf := &firebase.Config{
		DatabaseURL: "https://novel-fac48.firebaseio.com",
	}
	opt := option.WithCredentialsJSON([]byte(os.Getenv("FIREBASE")))
	app, err := firebase.NewApp(ctx, conf, opt)
	if err != nil {
		log.Fatalf("firebase.NewApp: %v", err)
	}
	client, err = app.Database(ctx)
	if err != nil {
		log.Fatalf("app.Firestore: %v", err)
	}
}

func getNovelList() ([]Novel, error) {
	var shallowNovels map[string]bool
	if err := client.NewRef("novels").GetShallow(context.Background(), &shallowNovels); err != nil {
		return nil, err
	}
	novels := make([]Novel, 0, len(shallowNovels))
	for k := range shallowNovels {
		novels = append(novels, Novel{Name: k})
	}
	return novels, nil
}

func formatResp(input interface{}) string {
	bytesBuffer := new(bytes.Buffer)
	json.NewEncoder(bytesBuffer).Encode(input)

	responseBytes := bytesBuffer.Bytes()

	var prettyJSON bytes.Buffer
	error := json.Indent(&prettyJSON, responseBytes, "", "  ")
	if error != nil {
		log.Println("JSON parse error: ", error)
	}
	formattedResp := string(prettyJSON.Bytes())
	return formattedResp
}
