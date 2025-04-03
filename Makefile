deps:
	go mod tidy

build:
	go build -o ./build/exchange ./main.go

clean:
	rm -rf ./build/*