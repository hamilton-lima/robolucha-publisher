# prepare builder
FROM golang:1.12.6 as builder

RUN mkdir -p /usr/local/share/robolucha-publisher
WORKDIR /usr/local/share/robolucha-publisher/

COPY go.mod /usr/local/share/robolucha-publisher/

# Authorize SSH Host
# RUN mkdir -p /root/.ssh && \
#     chmod 0700 /root/.ssh && \
#     ssh-keyscan gitlab.com > /root/.ssh/known_hosts

# Add the keys and set permissions
# COPY .ssh/id_rsa /root/.ssh/id_rsa
# RUN chmod 600 /root/.ssh/id_rsa 

# # Set git to use SSH instead of HTTPS
# RUN git config --global --add url."git@gitlab.com:robolucha/robolucha-publisher.git".insteadOf "https://gitlab.com/robolucha/robolucha-publisher.git"

# # get dependencies
# RUN bash -c $GOPATH/src/gitlab.com/robolucha/robolucha-publisher/go-get.sh
RUN go get 
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
