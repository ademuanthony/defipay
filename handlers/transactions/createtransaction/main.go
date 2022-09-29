package main

import (
	"context"
	"deficonnect/defipayapi/app"
	"deficonnect/defipayapi/handlers"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	
	lambda.Start(handler)
}

func handler(ctx context.Context, r events.APIGatewayProxyRequest) (app.Response, error) {
	var input app.CreateTransactionInput
	if err := json.Unmarshal([]byte(r.Body), &input); err != nil {
		app.Log.Error("Login", "json::Decode", err)
		return app.SendErrorfJSON("cannot decode request")
	}

	val := true
	m, err := handlers.InitSlsApp(val)
	if err != nil {
		panic(err)
	}

	account, _ := m.CurrentAccount(ctx, r)

	return m.CreateTransaction(ctx, input, account)
}
