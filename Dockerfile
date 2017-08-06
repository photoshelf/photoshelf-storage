FROM golang:alpine
MAINTAINER shunsuke maeda <duck8823@gmail.com>

RUN apk --update add --no-cache git

RUN mkdir -p /go/src/github.com/photoshelf/photoshelf-storage/images
WORKDIR /go/src/github.com/photoshelf/photoshelf-storage

ADD . .

RUN go get -u github.com/golang/dep/cmd/dep github.com/mattn/goveralls
RUN dep ensure

RUN go build

VOLUME /go/src/github.com/photoshelf/photoshelf-storage/images
EXPOSE 1323

CMD ["./photoshelf-storage"]
