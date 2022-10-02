.PHONY: build clean deploy gomodgen

build: gomodgen
	export GO111MODULE=on
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/GetTransaction handlers/transactions/GetTransaction/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/GetTransactions handlers/transactions/GetTransactions/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/CreateTransaction handlers/transactions/CreateTransaction/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/UpdateTransactionCurrency handlers/transactions/UpdateTransactionCurrency/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/CheckTransactionStatus handlers/transactions/CheckTransactionStatus/main.go

	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/Register handlers/auth/Register/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/Login handlers/auth/Login/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/Me handlers/auth/Me/main.go

	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/CreatePaymentLink handlers/paymentLink/CreatePaymentLink/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/GetPaymentLink handlers/paymentLink/GetPaymentLink/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/GetPaymentLinks handlers/paymentLink/GetPaymentLinks/main.go

	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/CreateBeneficiary handlers/beneficiary/CreateBeneficiary/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/GetBeneficiaries handlers/beneficiary/GetBeneficiaries/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/GetBeneficiary handlers/beneficiary/GetBeneficiary/main.go

	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/GetDfcEndpoint handlers/config/GetDfcEndpoint/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/SupportedCurrencies handlers/config/SupportedCurrencies/main.go
	
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/SlackCallback handlers/SlackCallback/main.go

	
	
	
	

clean:
	rm -rf ./bin ./vendor go.sum

deploy: clean build
	sls deploy --verbose

gomodgen:
	chmod u+x gomod.sh
	./gomod.sh
