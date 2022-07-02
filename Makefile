linux:
	GOARCH=amd64 GOOS=linux go build -o "bin/linux/unsubmail"

windows:
	GOARCH=amd64 GOOS=windows go build -o "bin/windows/unsubmail.exe"

darwin:
	GOARCH=amd64 GOOS=darwin go build -o "bin/darwin/unsubmail"

clean:
	rm -r bin/*

all: linux windows darwin