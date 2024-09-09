module github.com/bitcoin-sv/spv-wallet

go 1.22.5

// NOTE: The following replace directives are essential for maintaining the cohesion and functionality of this project.
// We are using the packages github.com/bitcoin-sv/spv-wallet/models and github.com/bitcoin-sv/spv-wallet/engine directly
// to facilitate the seamless integration of features across various components of our application.
// Removing these replaces could disrupt the interdependency between modules and hinder our ability to build cohesive features
// that often require modifications across multiple packages. Please refrain from removing these directives.
replace github.com/bitcoin-sv/spv-wallet/models => ./models

require (
	github.com/99designs/gqlgen v0.17.49
	github.com/DATA-DOG/go-sqlmock v1.5.2
	github.com/bitcoin-sv/go-broadcast-client v0.21.0
	github.com/bitcoin-sv/go-paymail v0.20.1
	github.com/bitcoin-sv/go-sdk v1.1.5
	github.com/bitcoin-sv/spv-wallet/models v0.28.0
	github.com/bitcoinschema/go-bitcoin/v2 v2.0.5
	github.com/bitcoinschema/go-map v0.1.1
	github.com/coocood/freecache v1.2.4
	github.com/fergusstrange/embedded-postgres v1.29.0
	github.com/gin-contrib/pprof v1.5.0
	github.com/gin-gonic/gin v1.10.0
	github.com/go-ozzo/ozzo-validation v3.6.0+incompatible
	github.com/go-redis/redis/v8 v8.11.5
	github.com/go-redis/redis_rate/v9 v9.1.2
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/google/uuid v1.6.0
	github.com/iancoleman/strcase v0.3.0
	github.com/jarcoal/httpmock v1.3.1
	github.com/libsv/go-bc v0.1.29
	github.com/libsv/go-bk v0.1.6
	github.com/libsv/go-bt v1.0.8
	github.com/libsv/go-bt/v2 v2.2.5
	github.com/mrz1836/go-cache v0.11.0
	github.com/mrz1836/go-cachestore v0.5.0
	github.com/mrz1836/go-logger v0.3.3
	github.com/mrz1836/go-sanitize v1.3.2
	github.com/mrz1836/go-validate v0.2.1
	github.com/newrelic/go-agent/v3 v3.34.0
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.20.3
	github.com/rafaeljusto/redigomock v2.4.0+incompatible
	github.com/robfig/cron/v3 v3.0.1
	github.com/rs/zerolog v1.33.0
	github.com/spf13/viper v1.19.0
	github.com/stretchr/testify v1.9.0
	github.com/swaggo/swag v1.16.3
	github.com/vmihailenco/taskq/v3 v3.2.9
	go.elastic.co/ecszerolog v0.2.0
	go.mongodb.org/mongo-driver v1.16.0
	gorm.io/driver/postgres v1.5.9
	gorm.io/driver/sqlite v1.5.6
	gorm.io/gorm v1.25.12
	gorm.io/plugin/dbresolver v1.5.2
)

require (
	github.com/bytedance/sonic/loader v0.2.0 // indirect
	github.com/cloudwego/base64x v0.1.4 // indirect
	github.com/cloudwego/iasm v0.2.0 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/xi2/xz v0.0.0-20171230120015-48954b6210f8 // indirect
)

require (
	github.com/KyleBanks/depth v1.2.1 // indirect
	github.com/asaskevich/govalidator v0.0.0-20230301143203-a9d515a09cc2 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bitcoinschema/go-bpu v0.1.3 // indirect
	github.com/bitcoinsv/bsvd v0.0.0-20190609155523-4c29707f7173 // indirect
	github.com/bitcoinsv/bsvutil v0.0.0-20181216182056-1d77cf353ea9 // indirect
	github.com/bsm/redislock v0.9.4 // indirect
	github.com/bytedance/sonic v1.12.0 // indirect
	github.com/capnm/sysinfo v0.0.0-20130621111458-5909a53897f3 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/dgryski/go-farm v0.0.0-20200201041132-a6ae2369ad13 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.5 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-openapi/jsonpointer v0.21.0 // indirect
	github.com/go-openapi/jsonreference v0.21.0 // indirect
	github.com/go-openapi/spec v0.21.0 // indirect
	github.com/go-openapi/swag v0.23.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.22.0 // indirect
	github.com/go-resty/resty/v2 v2.14.0 // indirect
	github.com/goccy/go-json v0.10.3 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/hashicorp/golang-lru v1.0.2 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgx/v5 v5.6.0 // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.17.9 // indirect
	github.com/klauspost/cpuid/v2 v2.2.8 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/libsv/go-p2p v0.3.2 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-sqlite3 v1.14.22 // indirect
	github.com/miekg/dns v1.1.62 // indirect
	github.com/mitchellh/mapstructure v1.5.0
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/montanaflynn/stats v0.7.1 // indirect
	github.com/pelletier/go-toml/v2 v2.2.2 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.55.0 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	github.com/sagikazarmark/locafero v0.6.0 // indirect
	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
	github.com/sosodev/duration v1.3.1 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/afero v1.11.0 // indirect
	github.com/spf13/cast v1.6.0 // indirect
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/objx v0.5.2 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	github.com/swaggo/files v1.0.1
	github.com/swaggo/gin-swagger v1.6.0
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.2.12 // indirect
	github.com/vektah/gqlparser/v2 v2.5.16 // indirect
	github.com/vmihailenco/msgpack/v5 v5.4.1 // indirect
	github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20240726163527-a2c0da244d78 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/arch v0.8.0 // indirect
	golang.org/x/crypto v0.26.0 // indirect
	golang.org/x/exp v0.0.0-20240719175910-8a7402abbf56 // indirect
	golang.org/x/mod v0.19.0 // indirect
	golang.org/x/net v0.28.0 // indirect
	golang.org/x/sync v0.8.0 // indirect
	golang.org/x/sys v0.24.0 // indirect
	golang.org/x/text v0.17.0 // indirect
	golang.org/x/tools v0.23.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240903143218-8af14fe29dc1 // indirect
	google.golang.org/grpc v1.66.0 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

// Issue with redislock package
replace github.com/bsm/redislock => github.com/bsm/redislock v0.7.2

// Issue with using wrong version of Redigo
replace github.com/gomodule/redigo => github.com/gomodule/redigo v1.8.9

// Issue: go.mongodb.org/mongo-driver/x/bsonx: cannot find module providing package go.mongodb.org/mongo-driver/x/bsonx
replace go.mongodb.org/mongo-driver => go.mongodb.org/mongo-driver v1.11.7
