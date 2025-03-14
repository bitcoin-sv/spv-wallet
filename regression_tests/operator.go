package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	walletclient "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
	"github.com/bitcoin-sv/spv-wallet-go-client/config"
	"github.com/bitcoin-sv/spv-wallet-go-client/walletkeys"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
)

const (
	domainLocalHost     = "http://localhost:3003"
	adminXPriv          = "xprv9s21ZrQH143K3CbJXirfrtpLvhT3Vgusdo8coBritQ3rcS7Jy7sxWhatuxG5h2y1Cqj8FKmPp69536gmjYRpfga2MJdsGyBsnB12E19CESK"
	leaderPaymailAlias  = "leader"
	minimalBalance      = 14 // 7 satoshi per test
	defaultGoClientPath = "../../spv-wallet-go-client/regression_tests"
	defaultJSClientPath = "../../spv-wallet-js-client/src/regression_tests"
)

var (
	ErrTimeout = errors.New("timeout reached")
)

func main() {
	loadConfigFlag := flag.Bool("l", false, "Load configuration from .env.config file")
	flag.Parse()
	ctx := context.Background()

	config := &regressionTestConfig{}
	user := &regressionTestUser{}
	var err error

	config, err = getConfig(*loadConfigFlag, config)
	if err != nil {
		fmt.Println("error getting config:", err)
		return
	}

	if !isSPVWalletRunning(domainLocalHost) {
		fmt.Println("spv-wallet is not running. Run spv-wallet and try again")
		return
	}

	sharedConfig, err := getSharedConfig(ctx)
	if err != nil {
		fmt.Println("error getting shared config:", err)
		return
	}

	setConfigClientsUrls(config, domainLocalHost)

	paymail := sharedConfig.PaymailDomains[0]
	if !isSPVWalletRunning(paymail) {
		fmt.Println("can't connect to spv-wallet at", paymail, "\nEnable tunneling from localhost")
		return
	}

	user, err = handleUserCreation(paymail, config)
	if err != nil {
		fmt.Println("error handling user creation:", err)
		return
	}

	setConfigLeaderXPriv(config, user.XPriv)

	if err := handleFundsTransfer(user, config); err != nil {
		fmt.Println("error handling transactions:", err)
		return
	}

	clientType := ""
	for clientType != "go" && clientType != "js" {
		clientType = promptUser("Do you want to run tests from go-client or js-client? (enter 'go' or 'js'): ")
		clientType = strings.ToLower(clientType)
	}

	var defaultPath string
	if clientType == "go" {
		defaultPath = defaultGoClientPath
	} else {
		defaultPath = defaultJSClientPath
	}
	err = saveConfig(config)
	if err != nil {
		fmt.Println("error saving config:", err)
		return
	}
	setEnvVariables(config)

	if err := runTests(clientType, defaultPath); err != nil {
		fmt.Println("error running tests:", err)
	}
}

// getConfig retrieves the configuration from the environment or from the user.
func getConfig(loadConfigFlag bool, config *regressionTestConfig) (*regressionTestConfig, error) {
	var err error
	if loadConfigFlag {
		config, err = getEnvVariables()
		if err != nil {
			return nil, err
		}
	}
	return config, nil
}

// handleUserCreation handles the creation of a user.
func handleUserCreation(paymailAlias string, config *regressionTestConfig) (*regressionTestUser, error) {
	if config.ClientOneLeaderXPriv != "" {
		answer, err := promptUserAndCheck("Would you like to use user from env? (y/yes or n/no): ")
		if err != nil {
			return nil, fmt.Errorf("failed to prompt user: %w", err)
		}
		if answer == yes {
			return useUserFromEnv(paymailAlias, config)
		}
	}

	user, err := createUser(paymailAlias, config)
	if err != nil {
		return handleCreateUserError(err, paymailAlias, config)
	}
	return user, nil
}

// handleCreateUserError handles the error when creating a user.
func handleCreateUserError(err error, paymailAlias string, config *regressionTestConfig) (*regressionTestUser, error) {
	if errors.Is(err, spverrors.ErrPaymailAlreadyExists) {
		return handleExistingPaymail(paymailAlias, config)
	} else {
		return nil, fmt.Errorf("error creating user: %w", err)
	}
}

// handleExistingPaymail handles the case when the paymail already exists.
func handleExistingPaymail(paymailAlias string, config *regressionTestConfig) (*regressionTestUser, error) {

	answer, err := promptUserAndCheck("Paymail already exists. Would you like to use it (you need to have xpriv)? (y/yes or n/no): ")
	if err != nil {
		return nil, fmt.Errorf("failed to prompt user: %w", err)
	}
	if answer == yes {
		return useUserFromXPriv(paymailAlias)
	}

	answer, err = promptUserAndCheck("Would you like to recreate user? (y/yes or n/no):")
	if err != nil {
		return nil, fmt.Errorf("failed to prompt user: %w", err)
	}
	if answer == yes {
		return recreateUser(paymailAlias, config)
	}
	return nil, fmt.Errorf("user should be recreated or xpriv should be provided")
}

// recreateUser deletes paymail and recreates it with new set of keys.
func recreateUser(paymailAlias string, config *regressionTestConfig) (*regressionTestUser, error) {
	err := deleteUser(paymailAlias, config)
	if err != nil {
		return nil, fmt.Errorf("error deleting user: %w", err)
	}
	user, err := createUser(paymailAlias, config)
	if err != nil {
		return nil, fmt.Errorf("error creating user after deletion: %w", err)
	}
	return user, nil
}

// useUserFromXPriv fills missing user data using provided xpriv.
func useUserFromXPriv(paymailAlias string) (*regressionTestUser, error) {
	validatedXPriv := getValidXPriv()
	xPriv, err := walletkeys.XPrivFromString(validatedXPriv)
	if err != nil {
		return nil, fmt.Errorf("failed to generate an extended private key (xPriv) from the given string: %w", err)
	}

	xPub, err := walletkeys.XPubFromXPriv(validatedXPriv)
	if err != nil {
		return nil, fmt.Errorf("failed to derive an extended public key (xPub) from the given string: %w", err)
	}

	return &regressionTestUser{
		XPriv:   xPriv.String(),
		XPub:    xPub,
		Paymail: preparePaymail(leaderPaymailAlias, paymailAlias),
	}, nil
}

// handleFundsTransfer handles the transfer of funds to a user.
func handleFundsTransfer(user *regressionTestUser, config *regressionTestConfig) (err error) {
	response, err := promptUserAndCheck("Do you have xpriv and master instance URL? (y/yes or n/no): ")
	if err != nil {
		return fmt.Errorf("failed to prompt user: %w", err)
	}
	if response == no {
		fmt.Printf("Please send %d Sato for full regression tests:\n%s\n", minimalBalance, user.Paymail)
		isSent := wrongInput
		for isSent < no {
			isSent, err = promptUserAndCheck("Did you make the transaction? (y/yes or n/no): ")
			if err != nil {
				return fmt.Errorf("failed to prompt user: %w", err)
			}
			if isSent == no {
				fmt.Println("Checking balance")
			}
		}
	} else if err := takeMasterUrlAndXPriv(user); err != nil {
		return fmt.Errorf("error handling funds transfer: %w", err)
	}

	leaderBalance, err := checkBalance(config.ClientOneURL, config.ClientOneLeaderXPriv)
	if err != nil {
		return fmt.Errorf("error checking balance: %w", err)
	}
	fmt.Println("Leader balance:", leaderBalance)
	timeout := time.After(2 * timeoutDuration)
	for leaderBalance < minimalBalance {
		select {
		case <-timeout:
			return ErrTimeout
		default:
			fmt.Print("Waiting for funds")
			for i := 0; i < 3; i++ {
				fmt.Print(".")
				time.Sleep(1 * time.Second)
			}
			fmt.Println()
			leaderBalance, err = checkBalance(config.ClientOneURL, config.ClientOneLeaderXPriv)
			if err != nil {
				return fmt.Errorf("error checking balance: %w", err)
			}
		}
	}
	return nil
}

// takeMasterUrlAndXPriv takes the master URL and xpriv for transferring funds.
func takeMasterUrlAndXPriv(leaderPaymail *regressionTestUser) error {
	url := getValidURL()
	xprivMaster := getValidXPriv()

	err := sendFundsWithGoClient(url, xprivMaster, leaderPaymail.Paymail)
	if err != nil {
		return fmt.Errorf("error sending funds with go client: %w", err)
	}
	return nil
}

// sendFundsWithGoClient sends funds using the Go client.
func sendFundsWithGoClient(instanceUrl string, istanceXPriv string, receiverPaymail string) error {
	client, err := walletclient.NewUserAPIWithXPriv(config.New(
		config.WithAddr(instanceUrl),
	), istanceXPriv)

	if err != nil {
		return fmt.Errorf("error initalizing client: %w", err)
	}
	ctx := context.Background()

	balance, err := checkBalance(instanceUrl, istanceXPriv)
	if err != nil {
		return fmt.Errorf("error checking balance: %w", err)
	}
	if balance <= minimalBalance {
		return fmt.Errorf("balance too low: %d", balance)
	}

	_, err = client.SendToRecipients(ctx, &commands.SendToRecipients{
		Recipients: []*commands.Recipients{
			{
				Satoshis: uint64(minimalBalance),
				To:       receiverPaymail,
			},
		},
		Metadata: map[string]any{"message": "regression test funds"},
	})
	if err != nil {
		return fmt.Errorf("error sending to recipients: %w", err)
	}
	return nil
}

// runTests runs the regression tests, asks for type of client and path to it and executes command.
func runTests(clientType string, defaultPath string) error {
	var command string
	if clientType == "go" {
		command = "go test -tags=regression ./... -count=1"
	}
	if clientType == "js" {
		command = "yarn install && yarn test:regression"
	}
	var path string
	var cmd *exec.Cmd
	for {
		path = promptUser(fmt.Sprintf("Enter relative path to the %s client (default: %s): ", clientType, defaultPath))
		if path == "" {
			path = defaultPath
		}

		if _, err := os.Stat(path); os.IsNotExist(err) {
			fmt.Printf("Path %s does not exist. Please enter a valid path.\n", path)
		} else {
			break
		}
	}

	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", fmt.Sprintf("cd %s && %s", path, command))
	} else {
		cmd = exec.Command("sh", "-c", fmt.Sprintf("cd %s && %s", path, command))
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error running tests: %w", err)
	}
	return nil
}
