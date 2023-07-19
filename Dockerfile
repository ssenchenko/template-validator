# Use the official Golang image to install Golang
FROM golang:latest as go_image

RUN apt-get update
RUN apt-get -y install python3
RUN apt-get -y install python3-setuptools
RUN apt-get -y install python3-pip

# Install necessary dependencies for AWS CLI and AWS SAM CLI
RUN apt-get update && apt-get install -y \
    unzip

RUN pip install awscli --break-system-packages

# Download and install AWS SAM CLI
RUN curl -sSL "https://github.com/aws/aws-sam-cli/releases/latest/download/aws-sam-cli-linux-x86_64.zip" -o "aws-sam-cli-linux-x86_64.zip" \
    && unzip aws-sam-cli-linux-x86_64.zip -d /tmp/sam-installation \
    && /tmp/sam-installation/install \
    && rm -rf aws-sam-cli-linux-x86_64.zip /tmp/sam-installation

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY *.go ./

COPY test/ ./test/

RUN go build -o bin/job

CMD ["./bin/job"]

# CMD ["sam", "validate", "-t", "test/error1.yaml", "--region", "us-west-2"]
