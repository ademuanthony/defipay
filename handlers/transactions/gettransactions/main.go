package main

import (
	"deficonnect/defipayapi/handlers"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	app, err := handlers.InitSlsApp()
	if err != nil {
		panic(err)
	}
	lambda.Start(app.GetTransactions)
}
