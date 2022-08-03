#build stage
#FROM golang:alpine AS builder
#RUN apk add --no-cache git
#WORKDIR /go/src/app
#COPY . .
#RUN go get -d -v ./...
#RUN go build -o /go/bin/app -v ./...

#final stage
#FROM alpine:latest
#RUN apk --no-cache add ca-certificates
#COPY --from=builder /go/bin/app /app
#ENTRYPOINT /app
#LABEL Name=bnmobackend Version=0.0.1
#EXPOSE 8000


FROM golang:alpine

RUN mkdir /bnmo-backend
RUN apk add --no-cache git

WORKDIR /bnmo-backend

ADD go.mod .
ADD go.sum .

RUN go mod download
ADD . .

RUN go install github.com/githubnemo/CompileDaemon

EXPOSE 8000

ENTRYPOINT CompileDaemon --build="go build main.go" --command="./main"