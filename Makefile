TARGET=goRedisJieba

all: mac

linux: 
	GOOS=linux GOARCH=amd64 go build -o ./bin/${TARGET}_${@} ./src

mac: 
	GOOS=darwin GOARCH=amd64 go build -o ./bin/${TARGET}_${@} ./src
	
clean:
	rm -rf ./bin/${TARGET}_*	
