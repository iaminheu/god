module git.zc0901.com/go/god

go 1.15

require (
	github.com/alicebob/miniredis/v2 v2.14.1
	github.com/beanstalkd/go-beanstalk v0.1.0
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/go-sql-driver/mysql v1.5.0
	github.com/go-xorm/builder v0.3.4
	github.com/golang/mock v1.4.4
	github.com/golang/protobuf v1.4.3
	github.com/google/uuid v1.0.0
	github.com/justinas/alice v1.2.0
	github.com/logrusorgru/aurora v2.0.3+incompatible
	github.com/onsi/ginkgo v1.14.2 // indirect
	github.com/onsi/gomega v1.10.3 // indirect
	github.com/prometheus/client_golang v1.8.0
	github.com/prometheus/common v0.14.0
	github.com/spaolacci/murmur3 v1.1.0
	github.com/stretchr/testify v1.6.1
	github.com/urfave/cli v1.22.5
	github.com/xwb1989/sqlparser v0.0.0-20180606152119-120387863bf2
	github.com/yuin/gopher-lua v0.0.0-20200816102855-ee81675732da // indirect
	go.etcd.io/etcd v0.0.0-20200402134248-51bdeb39e698
	golang.org/x/net v0.0.0-20201006153459-a7d1128ccaa0
	google.golang.org/grpc v1.29.1
	google.golang.org/protobuf v1.25.0 // indirect
	gopkg.in/yaml.v2 v2.3.0
)

//replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
