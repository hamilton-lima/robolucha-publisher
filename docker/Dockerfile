# prepare builder
FROM golang:1.12.6 as builder

RUN mkdir -p /usr/local/share/robolucha-publisher
WORKDIR /usr/local/share/robolucha-publisher/

# get dependencies
COPY go.mod /usr/local/share/robolucha-publisher/
RUN go get 

# copy source code
COPY . /usr/local/share/robolucha-publisher/

# build the binary static linked
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /tmp/publisher
RUN chmod +x /tmp/publisher

# start from scratch
FROM alpine
EXPOSE 5000

RUN mkdir -pv /usr/src/app

# Copy our static executable
COPY --from=builder /tmp/publisher /usr/src/app
RUN ls -alh /usr/src/app
ENTRYPOINT ["/usr/src/app/publisher"]
