build:
	go build -o ./bin/app ./cmd/app

run: build
	./bin/app

test:
	go test ./... -v -coverprofile=c.out && cat c.out | grep -v mock_ > filtered.c.out

cover-html:
	go tool cover -html=filtered.c.out

cover-txt:
	go tool cover -func=filtered.c.out >> coverage.txt

test-cover: test cover-html cover-txt
