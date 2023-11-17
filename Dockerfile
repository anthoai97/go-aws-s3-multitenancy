# syntax=docker/dockerfile:1
FROM golang:1.21.1
WORKDIR /app
ADD . /app
RUN go mod download
RUN go env -w GO111MODULE=on
RUN go build -o /docker-gs-ping
EXPOSE 8080
CMD [ "/docker-gs-ping" ] 