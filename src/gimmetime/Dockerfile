FROM golang:latest

RUN mkdir /go/src/gimmetime
WORKDIR /go/src/gimmetime
COPY . .
ENV GOPATH=/go

RUN go get -d -v ./
RUN go install -v ./

EXPOSE 80

CMD ["/go/bin/gimmetime"]