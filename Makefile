build:
	docker build -t validation-job .

run:
	docker run -it validation-job

build-local:
	go build -o bin/job

run-local:
	./bin/job
