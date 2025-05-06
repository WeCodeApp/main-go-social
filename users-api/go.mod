module users-api

go 1.24.2

require (
	common v0.0.0
	github.com/golang-jwt/jwt/v4 v4.5.2
	github.com/spf13/viper v1.20.1
	go.uber.org/zap v1.27.0
	golang.org/x/oauth2 v0.26.0
	google.golang.org/grpc v1.72.0
	gorm.io/driver/mysql v1.5.6
	gorm.io/gorm v1.25.8
)

replace common => ../common

require (
	cloud.google.com/go/compute/metadata v0.6.0 // indirect
	github.com/fsnotify/fsnotify v1.8.0 // indirect
	github.com/go-sql-driver/mysql v1.7.0 // indirect
	github.com/go-viper/mapstructure/v2 v2.2.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/pelletier/go-toml/v2 v2.2.3 // indirect
	github.com/sagikazarmark/locafero v0.7.0 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/afero v1.12.0 // indirect
	github.com/spf13/cast v1.7.1 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	golang.org/x/net v0.35.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/text v0.22.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250218202821-56aae31c358a // indirect
	google.golang.org/protobuf v1.36.6 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
