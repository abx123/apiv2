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

// Chapter ...
type Chapter struct {
	Title   string `json:"title"`
	Text    string `json:"text"`
	Link    string `json:"link"`
	Chapter int64  `json:"chapter"`
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

	if request.QueryStringParameters["novel"] == "" {
		err = fmt.Errorf("Missing novel")
	}
	resp, err = getChapterList(request.QueryStringParameters["novel"])

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

func getChapterList(novel string) ([]Chapter, error) {
	q := client.NewRef("novels/" + novel).OrderByChild("chapter")
	result, err := q.GetOrdered(context.Background())
	if err != nil {
		return nil, err
	}
	var chapters []Chapter

	var ch Chapter
	for _, c := range result {
		c.Unmarshal(&ch)
		chapters = append(chapters, ch)
	}
	return chapters, nil
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
