module github.com/bsv-blockchain/spv-wallet/api/manualtests

go 1.26.0

replace github.com/bsv-blockchain/spv-wallet/models => ../../models //nolint:gomoddirectives // local development

// Issue with godotenv v1.6.0 pre-release
replace github.com/joho/godotenv => github.com/joho/godotenv v1.5.1 //nolint:gomoddirectives // version override

require (
	github.com/bsv-blockchain/go-sdk v1.4.1
	github.com/bsv-blockchain/spv-wallet-go-client v1.3.0
	github.com/go-viper/mapstructure/v2 v2.5.0
	github.com/joomcode/errorx v1.2.0
	github.com/oapi-codegen/runtime v1.7.0
	github.com/rs/zerolog v1.35.1
	github.com/samber/lo v1.53.0
	github.com/spf13/viper v1.21.0
	github.com/stretchr/testify v1.12.1
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/apapsch/go-jsonmerge/v2 v2.0.0 // indirect
	github.com/boombuler/barcode v1.1.0 // indirect
	github.com/bsv-blockchain/spv-wallet/models v1.0.1 // indirect
	github.com/fsnotify/fsnotify v1.10.1 // indirect
	github.com/go-resty/resty/v2 v2.17.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/mattn/go-colorable v0.1.15 // indirect
	github.com/mattn/go-isatty v0.0.24 // indirect
	github.com/pelletier/go-toml/v2 v2.4.3 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pquerna/otp v1.5.0 // indirect
	github.com/sagikazarmark/locafero v0.12.0 // indirect
	github.com/spf13/afero v1.15.0 // indirect
	github.com/spf13/cast v1.10.0 // indirect
	github.com/spf13/pflag v1.0.10 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	golang.org/x/crypto v0.57.0 // indirect
	golang.org/x/net v0.59.0 // indirect
	golang.org/x/sys v0.48.0 // indirect
	golang.org/x/text v0.42.0 // indirect
)
