FROM golang:1.14
RUN mkdir /go/src/linebot_api
WORKDIR /go/src/linebot_api
ADD . /go/src/linebot_api