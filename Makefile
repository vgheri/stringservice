#CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o testzombie .
default: stringservice
	docker build -f Dockerfile -t valeriogheri/stringservice:latest .

stringservice:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o stringservice

clean:
	rm stringservice
