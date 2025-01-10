package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	walletclient "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
	walletclientcfg "github.com/bitcoin-sv/spv-wallet-go-client/config"
	"github.com/bitcoin-sv/spv-wallet-go-client/walletkeys"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/models/response"
	"github.com/joho/godotenv"
)

const (
	atSign                 = "@"
	domainPrefix           = "https://"
	spvWalletIndexResponse = "Welcome to the SPV Wallet ✌(◕‿-)✌"

	ClientOneURLEnvVar         = "CLIENT_ONE_URL"
	ClientTwoURLEnvVar         = "CLIENT_TWO_URL"
	ClientOneLeaderXPrivEnvVar = "CLIENT_ONE_LEADER_XPRIV"
	ClientTwoLeaderXPrivEnvVar = "CLIENT_TWO_LEADER_XPRIV"

	timeoutDuration = 120 * time.Second

	yes        = 1
	no         = 0
	wrongInput = -1
)

var explicitHTTPURLRegex = regexp.MustCompile(`^https?://`)

type regressionTestUser struct {
	XPriv   string `json:"xpriv"`
	XPub    string `json:"xpub"`
	Paymail string `json:"paymail"`
}

type regressionTestConfig struct {
	ClientOneURL         string
	ClientTwoURL         string
	ClientOneLeaderXPriv string
	ClientTwoLeaderXPriv string
}

type WalletResponse struct {
	Message string `json:"message"`
}

// saveConfig saves the configuration to a .env.config file.
func saveConfig(config *regressionTestConfig) error {
	envMap := map[string]string{
		ClientOneURLEnvVar:         config.ClientOneURL,
		ClientTwoURLEnvVar:         config.ClientTwoURL,
		ClientOneLeaderXPrivEnvVar: config.ClientOneLeaderXPriv,
		ClientTwoLeaderXPrivEnvVar: config.ClientTwoLeaderXPriv,
	}

	err := godotenv.Write(envMap, ".env.config")
	if err != nil {
		return fmt.Errorf("error saving .env.config file: %w", err)
	}
	return nil
}

// getEnvVariables loads the configuration from a .env.config file.
func getEnvVariables() (*regressionTestConfig, error) {
	if err := godotenv.Load(".env.config"); err != nil {
		return nil, fmt.Errorf("error loading .env.config file: %w", err)
	}

	return &regressionTestConfig{
		ClientOneURL:         os.Getenv(ClientOneURLEnvVar),
		ClientTwoURL:         os.Getenv(ClientTwoURLEnvVar),
		ClientOneLeaderXPriv: os.Getenv(ClientOneLeaderXPrivEnvVar),
		ClientTwoLeaderXPriv: os.Getenv(ClientTwoLeaderXPrivEnvVar),
	}, nil
}

// isSPVWalletRunning checks if the SPV wallet is running and prints the specific message at the specified URL.
func isSPVWalletRunning(url string) bool {
	client := http.Client{
		Timeout: timeoutDuration,
	}
	url = addPrefixIfNeeded(url)
	resp, err := client.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		return false
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error reading response body:", err)
		return false
	}

	var walletResp WalletResponse
	if err := json.Unmarshal(body, &walletResp); err != nil {
		fmt.Println("error parsing response JSON:", err)
		return false
	}
	return walletResp.Message == spvWalletIndexResponse
}

// addPrefixIfNeeded adds the HTTPS prefix to the URL if it is missing.
func addPrefixIfNeeded(url string) string {
	if !isValidURL(url) {
		return domainPrefix + url
	}
	return url
}

// getSharedConfig retrieves the shared configuration from the SPV Wallet.
func getSharedConfig(ctx context.Context) (*response.SharedConfig, error) {
	wc, err := walletclient.NewAdminAPIWithXPriv(walletclientcfg.New(walletclientcfg.WithAddr(domainLocalHost)), adminXPriv)
	if err != nil {
		return nil, err
	}

	sharedConfig, err := wc.SharedConfig(ctx)
	if err != nil {
		return nil, err
	}
	if len(sharedConfig.PaymailDomains) != 1 {
		return nil, fmt.Errorf("expected 1 paymail domain, got %d", len(sharedConfig.PaymailDomains))
	}
	return sharedConfig, nil
}

// promptUserAndCheck prompts the user with a question and validates the response.
func promptUserAndCheck(question string) (int, error) {
	var checkResult int
	input := make(chan string)
	defer close(input)

	for {
		fmt.Println(question)
		input := make(chan string, 1)
		go getInput(input)

		select {
		case response := <-input:
			checkResult = checkResponse(response)
			if checkResult != wrongInput {
				return checkResult, nil
			}
			fmt.Println("Invalid response. Please answer y/yes or n/no.")
		case <-time.After(timeoutDuration):
			os.Exit(1)
		}
	}
}

func getInput(input chan string) {
	in := bufio.NewReader(os.Stdin)
	result, _ := in.ReadString('\n')
	input <- result
}

// checkResponse checks the response and returns an integer indicating the result.
func checkResponse(response string) int {
	response = strings.ToLower(strings.TrimSpace(response))
	switch response {
	case "yes", "y":
		return yes
	case "no", "n":
		return no
	default:
		return wrongInput
	}
}

// preparePaymail constructs a paymail address from the alias and domain.
func preparePaymail(paymailAlias string, domain string) string {
	return paymailAlias + atSign + domain
}

func createUser(paymail string, config *regressionTestConfig) (*regressionTestUser, error) {
	keys, err := walletkeys.RandomKeys()
	if err != nil {
		return nil, err
	}

	user := &regressionTestUser{
		XPriv:   keys.XPriv(),
		XPub:    keys.XPub(),
		Paymail: preparePaymail(leaderPaymailAlias, paymail),
	}

	adminClient, err := walletclient.NewAdminAPIWithXPriv(walletclientcfg.New(walletclientcfg.WithAddr(config.ClientOneURL)), adminXPriv)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()

	_, err = adminClient.CreateXPub(ctx, &commands.CreateUserXpub{
		Metadata: map[string]any{"purpose": "regression-tests"},
		XPub:     user.XPub,
	})
	if err != nil {
		fmt.Println("adminNewXpub failed with status code:", err)
		return nil, err
	}

	createPaymailRes, err := adminClient.CreatePaymail(ctx, &commands.CreatePaymail{
		Key:        user.XPub,
		Address:    user.Paymail,
		PublicName: "Regression tests",
	})

	if err != nil {
		if err.Error() == spverrors.ErrPaymailAlreadyExists.Error() {
			return user, err
		}
		return nil, err
	}

	fmt.Println(keys.XPriv())
	user.Paymail = preparePaymail(createPaymailRes.Alias, createPaymailRes.Domain)
	return user, nil
}

// useUserFromEnv fills missing user data using xpriv from the environment variables.
func useUserFromEnv(paymailDomain string, config *regressionTestConfig) (*regressionTestUser, error) {
	xPriv, err := walletkeys.XPrivFromString(config.ClientOneLeaderXPriv)
	if err != nil {
		return nil, fmt.Errorf("failed to generate an extended private key (xPriv) from a string: %w", err)
	}
	xPub, err := walletkeys.XPubFromXPriv(xPriv.String())
	if err != nil {
		return nil, fmt.Errorf("failed to generate an extended public key (xPub) from a xPriv string: %w", err)
	}

	return &regressionTestUser{
		XPriv:   xPriv.String(),
		XPub:    xPub,
		Paymail: preparePaymail(leaderPaymailAlias, paymailDomain),
	}, nil
}

// deleteUser soft deletes paymail address from the SPV Wallet.
func deleteUser(paymail string, config *regressionTestConfig) error {
	paymail = preparePaymail(leaderPaymailAlias, paymail)
	adminClient, err := walletclient.NewAdminAPIWithXPriv(walletclientcfg.New(walletclientcfg.WithAddr(config.ClientOneURL)), adminXPriv)
	if err != nil {
		return err
	}

	ctx := context.Background()
	err = adminClient.DeletePaymail(ctx, paymail)
	if err != nil {
		return err
	}

	return nil
}

// getValidXPriv prompts the user for a valid xpriv and returns it.
func getValidXPriv() string {
	for {
		xpriv := promptUser("Enter xpriv: ")
		if strings.HasPrefix(xpriv, "xprv") {
			return xpriv
		}
		fmt.Println("Invalid xpriv. Please enter a valid xpriv")
	}
}

// promptUser prompts the user with a question and returns the response.
func promptUser(question string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println(question)
	response, _ := reader.ReadString('\n')
	return strings.TrimSpace(response)
}

// getValidURL prompts the user for a valid URL and returns it.
func getValidURL() string {
	for {
		url := promptUser("Enter master instance URL with prefix: ")
		if isValidURL(url) {
			return url
		}
		fmt.Println("Invalid URL. Please enter a valid URL with http/https prefix")
	}
}

// isValidURL validates the URL if it has http or https prefix.
func isValidURL(rawURL string) bool {
	return explicitHTTPURLRegex.MatchString(rawURL)
}

// checkBalance checks the balance of the specified xpriv at the given domain with given xpriv.
func checkBalance(domain, xpriv string) (int, error) {
	client, err := walletclient.NewUserAPIWithXPriv(walletclientcfg.New(walletclientcfg.WithAddr(domain)), xpriv)
	if err != nil {
		return wrongInput, err
	}
	ctx := context.Background()

	xpubInfo, err := client.XPub(ctx)
	if err != nil {
		return wrongInput, fmt.Errorf("error getting xpub info: %w", err)
	}
	return int(xpubInfo.CurrentBalance), nil
}

// setConfigClientsUrls sets the environment domains ulrs variables in the config.
func setConfigClientsUrls(config *regressionTestConfig, domain string) {
	config.ClientOneURL = domain
	config.ClientTwoURL = domain
}

// setConfigLeaderXPriv sets the environment xprivs variables in the config.
func setConfigLeaderXPriv(config *regressionTestConfig, xPriv string) {
	config.ClientOneLeaderXPriv = xPriv
	config.ClientTwoLeaderXPriv = xPriv
}

// setEnvVariables sets the environment variables.
func setEnvVariables(config *regressionTestConfig) {
	os.Setenv(ClientOneURLEnvVar, config.ClientOneURL)
	os.Setenv(ClientTwoURLEnvVar, config.ClientTwoURL)
	os.Setenv(ClientOneLeaderXPrivEnvVar, config.ClientOneLeaderXPriv)
	os.Setenv(ClientTwoLeaderXPrivEnvVar, config.ClientTwoLeaderXPriv)
}
