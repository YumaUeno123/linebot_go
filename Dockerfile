FROM golang:1.14
# work copy
WORKDIR /go/src/github.com/YumaUeno123/linebot_go
COPY . .

# generate swagger document
RUN make install-tools

# go build
RUN go build -o main cmd/main.go