pjt-settings:
	mkdir -p cmd internal pkg vendor api configs build infra utils
	touch cmd/.keep internal/.keep pkg/.keep vendor/.keep api/.keep configs/.keep build/.keep infra/.keep utils/.kee

lint:
	clear
	golangci-lint run

run: lint
	go run main.go
	
test:
	go test -v ./...