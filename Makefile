.PHONY: build clean deploy gomodgen

build: gomodgen
	export GO111MODULE=on
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/gettransaction handlers/transactions/gettransaction/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/gettransactions handlers/transactions/gettransactions/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/createtransaction handlers/transactions/createtransaction/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/updatetransactioncurrency handlers/transactions/updatetransactioncurrency/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/checktransactionstatus handlers/transactions/checktransactionstatus/main.go

	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/register handlers/auth/register/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/login handlers/auth/login/main.go
	
	

clean:
	rm -rf ./bin ./vendor go.sum

deploy: clean build
	sls deploy --verbose

gomodgen:
	chmod u+x gomod.sh
	./gomod.sh
