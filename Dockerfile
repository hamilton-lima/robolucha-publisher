# prepare builder
FROM golang as builder
COPY . $GOPATH/src/github.com/hamilton-lima/robolucha-publisher/
WORKDIR $GOPATH/src/github.com/hamilton-lima/robolucha-publisher/

# get dependancies
RUN go get -d -v

# build the binary static linked
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/publisher
RUN chmod +x /go/bin/publisher

# start from scratch
FROM scratch
EXPOSE 5000

# Copy our static executable
COPY --from=builder /go/bin/publisher /go/bin/publisher
ENTRYPOINT ["/go/bin/publisher"]
