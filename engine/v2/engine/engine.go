package engine

import (
	"context"
	"errors"

	"github.com/bitcoin-sv/go-paymail/server"
	"github.com/bitcoin-sv/spv-wallet/config"
	"github.com/bitcoin-sv/spv-wallet/engine/chain"
	"github.com/bitcoin-sv/spv-wallet/engine/chain/models"
	"github.com/bitcoin-sv/spv-wallet/engine/paymail"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/utils"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/addresses"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/data"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/database/repository"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/engine/internal"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/fee"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/operations"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/paymails"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/paymailserver"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/beef"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/outlines"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/record"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/txsync"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/users"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/utils/must"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// V2 is the engine of the wallet, it is creating all needed services, and preparing database connection.
type V2 struct {
	cfg     *config.AppConfig
	storage *internal.Storage

	repositories *repository.All

	chainService chain.Service

	usersService               *users.Service
	paymailsService            *paymails.Service
	addressesService           *addresses.Service
	dataService                *data.Service
	operationsService          *operations.Service
	transactionsOutlineService outlines.Service
	transactionsRecordService  *record.Service
	txSyncService              *txsync.Service
	paymailServerConfig        *server.Configuration
}

// NewEngine creates a new engine.V2 instance.
func NewEngine(cfg *config.AppConfig, logger zerolog.Logger, overridesOpts ...InternalsOverride) *V2 {
	logger = logger.With().Int("v", 2).Str("service", "engine").Logger()

	overridesToApply := &overrides{}
	for _, opt := range overridesOpts {
		opt(overridesToApply)
	}

	// Database
	storage := internal.NewStorage(cfg, logger)
	err := storage.Start()
	must.HaveNoErrorf(err, "failed to start wallet storage")

	repos := storage.CreateRepositories()

	// Low level services
	var httpClient *resty.Client
	if overridesToApply.resty != nil {
		httpClient = overridesToApply.resty
	} else {
		httpClient = resty.New()
		if overridesToApply.transport != nil {
			httpClient.SetTransport(overridesToApply.transport)
		}
	}

	paymailClient := setupPaymailClient(overridesToApply, httpClient)

	cache := internal.NewCache(cfg, logger)

	chainService := chain.NewChainService(logger, httpClient, extractARCConfig(cfg), extractBHSConfig(cfg))
	feeService := fee.NewService(cfg, chainService, logger)

	feeUnit, err := feeService.GetFeeUnit(context.Background())
	must.HaveNoErrorf(err, "failed to setup fee unit")

	paymailServiceClient := paymail.NewServiceClient(cache, paymailClient, logger)

	utxoSelector := storage.CreateUTXOSelector(feeService)

	beefService := beef.NewService(repos.Transactions)

	userService := users.NewService(repos.Users, cfg)
	paymailService := paymails.NewService(repos.Paymails, userService, cfg)
	addressesService := addresses.NewService(repos.Addresses)
	dataService := data.NewService(repos.Data)
	operationsService := operations.NewService(repos.Operations)
	txSyncService := txsync.NewService(logger, repos.Transactions)

	transactionsOutlineService := outlines.NewService(
		paymailServiceClient,
		paymailService,
		beefService,
		utxoSelector,
		feeUnit,
		logger,
		userService,
	)

	transactionsRecordService := record.NewService(
		logger,
		addressesService,
		userService,
		repos.Outputs,
		repos.Operations,
		repos.Transactions,
		chainService,
		paymailServiceClient,
	)

	paymailServiceProvider := paymailserver.NewServiceProvider(
		logger,
		paymailService,
		userService,
		addressesService,
		chainService,
		transactionsRecordService,
	)

	paymailServerConfig := setupPaymailServer(cfg, logger, paymailServiceProvider)

	return &V2{
		cfg:                        cfg,
		storage:                    storage,
		repositories:               repos,
		chainService:               chainService,
		usersService:               userService,
		paymailsService:            paymailService,
		addressesService:           addressesService,
		dataService:                dataService,
		operationsService:          operationsService,
		transactionsOutlineService: transactionsOutlineService,
		transactionsRecordService:  transactionsRecordService,
		txSyncService:              txSyncService,
		paymailServerConfig:        paymailServerConfig,
	}
}

// Close closes the V2 and all its services
func (e *V2) Close(_ context.Context) error {
	var allErrors error
	err := e.storage.Close()
	if err != nil {
		allErrors = errors.Join(allErrors, spverrors.Wrapf(err, "couldn't close storage"))
	}
	return allErrors
}

// DB returns the database
// Deprecated: DB used as adapter for engine v1
func (e *V2) DB() *gorm.DB {
	return e.storage.DB()
}

// Repositories returns all repositories
func (e *V2) Repositories() *repository.All {
	return e.repositories
}

// Chain returns the chain service
func (e *V2) Chain() chain.Service {
	return e.chainService
}

// UsersService returns the users service
func (e *V2) UsersService() *users.Service {
	return e.usersService
}

// PaymailsService returns the paymails service
func (e *V2) PaymailsService() *paymails.Service {
	return e.paymailsService
}

// AddressesService returns the addresses service
func (e *V2) AddressesService() *addresses.Service {
	return e.addressesService
}

// DataService returns the data service
func (e *V2) DataService() *data.Service {
	return e.dataService
}

// OperationsService returns the operations service
func (e *V2) OperationsService() *operations.Service {
	return e.operationsService
}

// TransactionOutlinesService returns the transaction outlines service
func (e *V2) TransactionOutlinesService() outlines.Service {
	return e.transactionsOutlineService
}

// TransactionRecordService returns the transaction record service
func (e *V2) TransactionRecordService() *record.Service {
	return e.transactionsRecordService
}

// TxSyncService returns the tx sync service
func (e *V2) TxSyncService() *txsync.Service {
	return e.txSyncService
}

// PaymailServerConfiguration returns the paymail server configuration
func (e *V2) PaymailServerConfiguration() *server.Configuration {
	return e.paymailServerConfig
}

func extractBHSConfig(cfg *config.AppConfig) chainmodels.BHSConfig {
	return chainmodels.BHSConfig{
		URL:       cfg.BHS.URL,
		AuthToken: cfg.BHS.AuthToken,
	}
}

func extractARCConfig(cfg *config.AppConfig) chainmodels.ARCConfig {
	arcCfg := chainmodels.ARCConfig{
		URL:          cfg.ARC.URL,
		Token:        cfg.ARC.Token,
		DeploymentID: cfg.ARC.DeploymentID,
		WaitFor:      cfg.ARC.WaitForStatus,
	}

	if cfg.ARC.Callback.Enabled {
		var err error
		if cfg.ARC.Callback.Token == "" {
			// This also sets the token to the config reference and, it is used in the callbacktoken_middleware
			cfg.ARC.Callback.Token, err = utils.HashAdler32(config.DefaultAdminXpub)
			must.HaveNoErrorf(err, "error while generating callback token")
		}
		arcCfg.Callback = &chainmodels.ARCCallbackConfig{
			URL:   cfg.ARC.Callback.Host + config.BroadcastCallbackRoute,
			Token: cfg.ARC.Callback.Token,
		}
	}

	if cfg.ExperimentalFeatures != nil && cfg.ExperimentalFeatures.UseJunglebus {
		arcCfg.UseJunglebus = true
	}

	return arcCfg
}
