FROM golang:1.14.0-alpine as builder

# Enable Go Modules
ENV GO111MODULE=on

# package update & install
RUN apk update --no-cache && \
    apk add git make build-base

# work copy
WORKDIR /go/src/github.com/YumaUeno123/linebot_go
COPY . .

# download modules
RUN go mod download

# generate swagger document
RUN make install-tools

# go build
#RUN GOOS=linux GOARCH=amd64 go build -o main cmd/main.go
RUN go build -o main cmd/main.go

#exec container image
FROM alpine:3.10

# ca-certificates
RUN apk update && apk --no-cache add ca-certificates curl

#timezone set
RUN apk --no-cache add tzdata && \
    cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime && \
    apk del tzdata

WORKDIR /cmd
COPY --from=builder /go/src/github.com/YumaUeno123/linebot_go/main .
CMD /cmd/main -port $PORT