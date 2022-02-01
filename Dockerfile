# syntax=docker/dockerfile:1

FROM golang:1.17-alpine
WORKDIR /user-balance
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .
RUN go build -o ./app ./cmd/app
EXPOSE 8080
CMD ["/user-balance/app"]

