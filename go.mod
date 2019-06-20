module gin_project_starter

go 1.12

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/fsnotify/fsnotify v1.4.7
	github.com/gin-contrib/pprof v1.2.0
	github.com/gin-gonic/gin v1.4.0
	github.com/go-pg/pg v8.0.4+incompatible
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/onsi/ginkgo v1.8.0
	github.com/onsi/gomega v1.5.0
	github.com/pkg/errors v0.8.1
	github.com/rs/zerolog v1.14.3
	github.com/spf13/viper v1.4.0
	golang.org/x/crypto v0.0.0-20190308221718-c2843e01d9a2
	mellium.im/sasl v0.2.1 // indirect
)

replace github.com/ugorji/go v1.1.4 => github.com/ugorji/go/codec v0.0.0-20190204201341-e444a5086c43
