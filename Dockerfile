FROM golang:alpine
MAINTAINER shunsuke maeda <duck8823@gmail.com>

RUN apk --update add --no-cache git

WORKDIR /go/src/github.com/photoshelf/photoshelf-storage

ADD . .

RUN go get -u github.com/golang/dep/cmd/dep github.com/mattn/goveralls
RUN dep ensure

RUN go build -o photoshelf-storage ./cmd/*.go

RUN mkdir -p /photoshelf/photos
VOLUME /photoshelf/photos

EXPOSE 1323

CMD ["./photoshelf-storage", "-d", "/photoshelf/photos"]
