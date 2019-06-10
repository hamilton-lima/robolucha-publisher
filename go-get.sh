set -o xtrace

go get -v gitlab.com/robolucha/robolucha-publisher/session
go get -v gitlab.com/robolucha/robolucha-publisher/redis

go get -v github.com/gomodule/redigo/redis
go get -v github.com/gin-contrib/cors
go get -v github.com/gin-gonic/gin
go get -v github.com/sirupsen/logrus
go get -v gopkg.in/olahol/melody.v1
go get -v github.com/onsi/ginkgo
go get -v github.com/hpcloud/tail
go get -v github.com/onsi/gomega
