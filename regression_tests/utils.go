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

	walletclient "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-go-client/xpriv"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/joho/godotenv"
)

const (
	atSign                   = "@"
	domainPrefix             = "http://"
	domainSuffixSharedConfig = "/v1/shared-config"
	spvWalletIndexResponse   = "Welcome to the SPV Wallet ✌(◕‿-)✌"

	ClientOneURLEnvVar         = "CLIENT_ONE_URL"
	ClientTwoURLEnvVar         = "CLIENT_TWO_URL"
	ClientOneLeaderXPrivEnvVar = "CLIENT_ONE_LEADER_XPRIV"
	ClientTwoLeaderXPrivEnvVar = "CLIENT_TWO_LEADER_XPRIV"
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

// loadConfig loads the configuration from a .env.config file.
func loadConfig() (*regressionTestConfig, error) {
	if err := godotenv.Load(".env.config"); err != nil {
		return nil, fmt.Errorf("error loading .env.config file: %v", err)
	}

	return &regressionTestConfig{
		ClientOneURL:         os.Getenv(ClientOneURLEnvVar),
		ClientTwoURL:         os.Getenv(ClientTwoURLEnvVar),
		ClientOneLeaderXPriv: os.Getenv(ClientOneLeaderXPrivEnvVar),
		ClientTwoLeaderXPriv: os.Getenv(ClientTwoLeaderXPrivEnvVar),
	}, nil
}

// isSPVWalletRunning checks if the SPV wallet is running at the specified URL.
func isSPVWalletRunning(url string) bool {
	url = addPrefixIfNeeded(url)
	resp, err := http.Get(url)
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

// addPrefixIfNeeded adds the HTTP prefix to the URL if it is missing.
func addPrefixIfNeeded(url string) string {
	if !strings.HasPrefix(url, domainPrefix) {
		return domainPrefix + url
	}
	return url
}

// getSharedConfig retrieves the shared configuration from the SPV Wallet.
func getSharedConfig(xpub string) (*models.SharedConfig, error) {
	req, err := http.NewRequest(http.MethodGet, addPrefixIfNeeded(domainLocalHost)+domainSuffixSharedConfig, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set(models.AuthHeader, xpub)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get shared config: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var configResponse models.SharedConfig
	if err := json.Unmarshal(body, &configResponse); err != nil {
		return nil, err
	}

	if len(configResponse.PaymailDomains) != 1 {
		return nil, fmt.Errorf("expected 1 paymail domain, got %d", len(configResponse.PaymailDomains))
	}
	return &configResponse, nil
}

// promptUserAndCheck prompts the user with a question and validates the response.
func promptUserAndCheck(question string) int {
	reader := bufio.NewReader(os.Stdin)
	var response string
	var checkResult int

	for {
		fmt.Println(question)
		response, _ = reader.ReadString('\n')
		response = strings.TrimSpace(response)

		checkResult = checkResponse(response)
		if checkResult != -1 {
			break
		}
		fmt.Println("Invalid response. Please answer y/yes or n/no.")
	}

	return checkResult
}

// checkResponse checks the response and returns an integer indicating the result.
func checkResponse(response string) int {
	response = strings.ToLower(strings.TrimSpace(response))
	switch response {
	case "yes", "y":
		return 1
	case "no", "n":
		return 0
	default:
		return -1
	}
}

// preparePaymail constructs a paymail address from the alias and domain.
func preparePaymail(paymailAlias string, domain string) string {
	return paymailAlias + atSign + domain
}

// createUser creates a new user with the specified paymail.
func createUser(paymail string, config *regressionTestConfig) (*regressionTestUser, error) {
	keys, err := xpriv.Generate()
	if err != nil {
		return nil, err
	}

	user := &regressionTestUser{
		XPriv:   keys.XPriv(),
		XPub:    keys.XPub().String(),
		Paymail: preparePaymail(paymail, config.ClientOneURL),
	}

	adminClient := walletclient.NewWithAdminKey(addPrefixIfNeeded(domainLocalHost), adminXPriv)
	ctx := context.Background()

	if err := adminClient.AdminNewXpub(ctx, user.XPub, map[string]any{"some_metadata": "remove"}); err != nil {
		fmt.Println("adminNewXpub failed with status code:", err)
		return nil, err
	}

	createPaymailRes, err := adminClient.AdminCreatePaymail(ctx, user.XPub, user.Paymail, "Test test", "")
	if err != nil {
		if err.Error() == "paymail address already exists" {
			return user, fmt.Errorf("paymail address already exists")
		}
		return nil, err
	}

	fmt.Println(keys.XPriv())
	user.Paymail = preparePaymail(createPaymailRes.Alias, createPaymailRes.Domain)
	return user, nil
}

// useUserFromEnv uses the user from the environment variables.
func useUserFromEnv(config *regressionTestConfig, paymailAlias string) (*regressionTestUser, error) {
	keys, err := xpriv.FromString(config.ClientOneLeaderXPriv)
	if err != nil {
		return nil, fmt.Errorf("error parsing xpriv: %v", err)
	}
	return &regressionTestUser{
		XPriv:   keys.XPriv(),
		XPub:    keys.XPub().String(),
		Paymail: preparePaymail(paymailAlias, config.ClientOneURL),
	}, nil
}

// deleteUser deletes the user with the specified paymail address from the SPV Wallet.
func deleteUser(paymail string, config *regressionTestConfig) error {
	paymail = preparePaymail(paymail, config.ClientOneURL)
	adminClient := walletclient.NewWithAdminKey(addPrefixIfNeeded(domainLocalHost), adminXPriv)
	ctx := context.Background()
	err := adminClient.AdminDeletePaymail(ctx, paymail)
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

// isValidURL validates the URL.
func isValidURL(rawURL string) bool {
	return explicitHTTPURLRegex.MatchString(rawURL)
}

// checkBalance checks the balance of the specified xpriv at the given domain.
func checkBalance(domain, xpriv string) int {
	client := walletclient.NewWithXPriv(addPrefixIfNeeded(domain), xpriv)
	ctx := context.Background()

	xpubInfo, err := client.GetXPub(ctx)
	if err != nil {
		fmt.Println("error getting xpub info:", err)
		os.Exit(1)
	}
	return int(xpubInfo.CurrentBalance)
}

// setConfigClientsUrls sets the environment domains variables in the config.
func setConfigClientsUrls(config *regressionTestConfig, domain string) {
	config.ClientOneURL = domain
	config.ClientTwoURL = domain
}

// setConfigLeaderXPriv sets the environment xprivs variables in the config.
func setConfigLeaderXPriv(config *regressionTestConfig, xPriv string) {
	config.ClientOneLeaderXPriv = xPriv
	config.ClientTwoLeaderXPriv = xPriv
}
