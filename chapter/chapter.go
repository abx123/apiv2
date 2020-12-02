package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/db"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"google.golang.org/api/option"
)

var client *db.Client

// Novel ...
type Novel struct {
	Name     string
	Chapters []Chapter
}

// Chapter ...
type Chapter struct {
	Title   string `json:"title"`
	Text    string `json:"text"`
	Link    string `json:"link"`
	Chapter int64  `json:"chapter"`
}

func main() {
	Init()
	lambda.Start(handleRequest)
}

// The input type and the output type are defined by the API Gateway.
func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var err error
	var resp interface{}
	headers := map[string]string{
		"Access-Control-Allow-Headers": "Content-Type",
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods": "GET",
	}

	chapter, _ := strconv.ParseInt(request.QueryStringParameters["chapter"], 10, 64)
	if chapter == 0 || request.QueryStringParameters["novel"] == "" {
		err = fmt.Errorf("Missing novel")

	}

	resp, err = getChapter(request.QueryStringParameters["novel"], chapter)

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

func getChapter(novel string, chapter int64) (Chapter, error) {
	var chMap map[string]interface{}
	if err := client.NewRef("novels/"+novel+fmt.Sprintf("/%d", chapter)).Get(context.Background(), &chMap); err != nil {
		return Chapter{}, err
	}
	ch := Chapter{
		Title:   chMap["title"].(string),
		Text:    chMap["text"].(string),
		Link:    chMap["link"].(string),
		Chapter: int64(chMap["chapter"].(float64)),
	}
	fmt.Println(ch)
	return ch, nil
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
