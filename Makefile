linux:
	GOARCH=amd64 GOOS=linux go build -o "bin/linux/unsubmail"

windows:
	GOARCH=amd64 GOOS=windows go build -o "bin/windows/unsubmail.exe"

clean:
	rm -r bin/*

all: linux windows