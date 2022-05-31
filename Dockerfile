# syntax=docker/dockerfile:1

FROM golang:1.16-alpine

# Set destination for COPY
WORKDIR /app

ENV GOPROXY=https://goproxy.cn,direct

# Download Go modules
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY docker ./
COPY http ./http/
COPY logger ./logger/
COPY secret/secret.go ./secret/

# Build
RUN go build -o /docker-eds-logger

# This is for documentation purposes only.
# To actually open the port, runtime parameters
# must be supplied to the docker command.
EXPOSE 8080

# (Optional) environment variable that our dockerised
# application can make use of. The value of environment
# variables can also be set via parameters supplied
# to the docker command on the command line.
#ENV HTTP_PORT=8081

# Run
CMD [ "/docker-eds-logger" ]
