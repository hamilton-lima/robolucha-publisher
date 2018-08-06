FROM golang
COPY . $GOPATH/src/github.com/hamilton-lima/robolucha-services/
WORKDIR $GOPATH/src/github.com/hamilton-lima/robolucha-services/
RUN go get -d -v

RUN go build -o /go/bin/publisher
RUN chmod +x /go/bin/publisher

EXPOSE 5000
ENTRYPOINT ["/go/bin/publisher"]
