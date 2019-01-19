# prepare builder
FROM golang as builder
COPY . $GOPATH/src/github.com/hamilton-lima/robolucha-publisher/
WORKDIR $GOPATH/src/github.com/hamilton-lima/robolucha-publisher/

# Authorize SSH Host
RUN mkdir -p /root/.ssh && \
    chmod 0700 /root/.ssh && \
    ssh-keyscan gitlab.com > /root/.ssh/known_hosts

# Add the keys and set permissions
COPY .ssh/id_rsa /root/.ssh/id_rsa
RUN chmod 600 /root/.ssh/id_rsa 

# Set git to use SSH instead of HTTPS
RUN git config --global --add url."git@gitlab.com:robolucha/robolucha-publisher.git".insteadOf "https://gitlab.com/robolucha/robolucha-publisher.git"

# get dependencies
RUN go get -v

# build the binary static linked
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/publisher
RUN chmod +x /go/bin/publisher

# start from scratch
FROM scratch
EXPOSE 5000

# Copy our static executable
COPY --from=builder /go/bin/publisher /go/bin/publisher
ENTRYPOINT ["/go/bin/publisher"]
