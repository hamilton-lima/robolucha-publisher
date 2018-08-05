# STEP 1 build executable binary
FROM golang as builder
COPY . $GOPATH/src/github.com/hamilton-lima/robolucha-services/
WORKDIR $GOPATH/src/github.com/hamilton-lima/robolucha-services/
#get dependancies
#you can also use dep
RUN go get -d -v

#build the binary
RUN go build publisher.go -o /go/bin/publisher
# STEP 2 build a small image
# start from scratch
FROM scratch
# Copy our static executable
COPY --from=builder /go/bin/publisher /go/bin/publisher
ENTRYPOINT ["/go/bin/publisher"]
