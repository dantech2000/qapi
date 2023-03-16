

build: clean
	go build -o ./bin/qapi main.go

clean:
	rm -rf ./bin