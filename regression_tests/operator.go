package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	walletclient "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-go-client/xpriv"
)

const (
	domainLocalHost     = "http://localhost:3003"
	adminXPriv          = "xprv9s21ZrQH143K3CbJXirfrtpLvhT3Vgusdo8coBritQ3rcS7Jy7sxWhatuxG5h2y1Cqj8FKmPp69536gmjYRpfga2MJdsGyBsnB12E19CESK"
	adminXPub           = "xpub661MyMwAqRbcFgfmdkPgE2m5UjHXu9dj124DbaGLSjaqVESTWfCD4VuNmEbVPkbYLCkykwVZvmA8Pbf8884TQr1FgdG2nPoHR8aB36YdDQh"
	leaderPaymailAlias  = "leader"
	minimalBalance      = 100
	defaultGoClientPath = "../../spv-wallet-go-client/regression_tests"
	defaultJSClientPath = "../../spv-wallet-js-client/regression_tests"
)

func main() {
	loadConfigFlag := flag.Bool("l", false, "Load configuration from .env.config file")
	flag.Parse()

	var config *regressionTestConfig
	var err error
	user := &regressionTestUser{}

	if *loadConfigFlag {
		config, err = loadConfig()
		if err != nil {
			fmt.Println("error loading config:", err)
			return
		}
		user.XPriv = config.ClientOneLeaderXPriv
	} else {
		config = &regressionTestConfig{}
	}

	if !isSPVWalletRunning(domainLocalHost) {
		fmt.Println("spv-wallet is not running. Run spv-wallet and try again")
		return
	}

	sharedConfig, err := getSharedConfig(adminXPub)
	if err != nil {
		fmt.Println("error getting shared config:", err)
		return
	}

	setConfigClientsUrls(config, sharedConfig.PaymailDomains[0])

	if !isSPVWalletRunning(config.ClientOneURL) {
		fmt.Println("can't connect to spv-wallet at", config.ClientOneURL, "\nEnable tunneling from localhost")
		return
	}

	user, err = handleUserCreation(leaderPaymailAlias, config)
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

// handleUserCreation handles the creation of a user.
func handleUserCreation(paymailAlias string, config *regressionTestConfig) (*regressionTestUser, error) {
	if config.ClientOneLeaderXPriv != "" {
		if promptUserAndCheck("Would you like to use user from env? (y/yes or n/no): ") == 1 {
			return useUserFromEnv(config, paymailAlias)
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
	if err.Error() == "paymail already exists" {
		return handleExistingPaymail(paymailAlias, config)
	} else {
		return nil, fmt.Errorf("error creating user: %v", err)
	}
}

// handleExistingPaymail handles the case when the paymail already exists.
func handleExistingPaymail(paymailAlias string, config *regressionTestConfig) (*regressionTestUser, error) {
	if promptUserAndCheck("Paymail already exists. Would you like to use it (you need to have xpriv)? (y/yes or n/no): ") == 1 {
		return useExistingPaymail(paymailAlias, config)
	}
	if promptUserAndCheck("Would you like to recreate user? (y/yes or n/no):") == 1 {
		return recreateUser(paymailAlias, config)
	}
	return nil, fmt.Errorf("user should be recreated or xpriv should be provided")
}

// recreateUser deletes and recreates the user.
func recreateUser(paymailAlias string, config *regressionTestConfig) (*regressionTestUser, error) {
	err := deleteUser(paymailAlias, config)
	if err != nil {
		return nil, fmt.Errorf("error deleting user: %v", err)
	}
	user, err := createUser(paymailAlias, config)
	if err != nil {
		return nil, fmt.Errorf("error creating user after deletion: %v", err)
	}
	return user, nil
}

// useExistingPaymail uses an existing paymail address.
func useExistingPaymail(paymailAlias string, config *regressionTestConfig) (*regressionTestUser, error) {
	validatedXPriv := getValidXPriv()
	keys, err := xpriv.FromString(validatedXPriv)
	if err != nil {
		return nil, fmt.Errorf("error parsing xpriv: %v", err)
	}
	return &regressionTestUser{
		XPriv:   keys.XPriv(),
		XPub:    keys.XPub().String(),
		Paymail: preparePaymail(paymailAlias, config.ClientOneURL),
	}, nil
}

// handleFundsTransfer handles the transfer of funds to a user.
func handleFundsTransfer(user *regressionTestUser, config *regressionTestConfig) error {
	response := promptUserAndCheck("Do you have xpriv and master instance URL? (y/yes or n/no): ")
	if response == 0 {
		fmt.Printf("Please send %d Sato for full regression tests:\n%s\n", minimalBalance, user.Paymail)
		isSent := 0
		for isSent < 1 {
			isSent = promptUserAndCheck("Did you make the transaction? (y/yes or n/no): ")
		}
	} else {
		if err := takeMasterUrlAndXPriv(user); err != nil {
			return fmt.Errorf("error sending funds: %v", err)
		}
	}

	leaderBalance := 0
	for leaderBalance == 0 {
		fmt.Print("Waiting for funds")
		for i := 0; i < 3; i++ {
			fmt.Print(".")
			time.Sleep(1 * time.Second)
		}
		leaderBalance = checkBalance(config.ClientOneURL, config.ClientOneLeaderXPriv)
		fmt.Println()
	}
	return nil
}

// takeMasterUrlAndXPriv takes the master URL and xpriv for transferring funds.
func takeMasterUrlAndXPriv(leaderPaymail *regressionTestUser) error {
	url := getValidURL()
	xprivMaster := getValidXPriv()

	err := sendFundsWithGoClient(url, xprivMaster, leaderPaymail.Paymail)
	if err != nil {
		fmt.Println("error sending funds:", err)
		return err
	}
	return nil
}

// sendFundsWithGoClient sends funds using the Go client.
func sendFundsWithGoClient(instanceUrl string, istanceXPriv string, receiverPaymail string) error {
	client := walletclient.NewWithXPriv(addPrefixIfNeeded(instanceUrl), istanceXPriv)
	ctx := context.Background()

	balance := checkBalance(instanceUrl, istanceXPriv)
	if balance <= minimalBalance {
		return fmt.Errorf("balance too low: %d", balance)
	}
	recipient := walletclient.Recipients{To: receiverPaymail, Satoshis: uint64(minimalBalance)}
	recipients := []*walletclient.Recipients{&recipient}

	_, err := client.SendToRecipients(ctx, recipients, map[string]any{"message": "regression test funds"})
	if err != nil {
		return fmt.Errorf("error sending to recipients: %v", err)
	}
	return nil
}

// runTests runs the regression tests, asks for type of client and path to it and executes command.
func runTests(clientType string, defaultPath string) error {
	// TODO: adjust command and path when regression tests are implemented
	var command string
	if clientType == "go" {
		command = "go test ./..."
	}
	if clientType == "js" {
		command = "yarn install && yarn test"
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
		return fmt.Errorf("error running tests: %v", err)
	}

	return nil
}
