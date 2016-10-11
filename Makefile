#CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o testzombie .
default: stringservice
	docker build -f Dockerfile -t valeriogheri/stringservice:latest .
	docker login -u valeriogheri -p MetopA_2016
	docker tag valeriogheri/stringservice:latest valeriogheri/stringservice:1.1.2
	docker push valeriogheri/stringservice

stringservice:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o stringservice

clean:
	rm stringservice
