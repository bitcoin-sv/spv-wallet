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

type User struct {
	XPriv   string `json:"xpriv"`
	XPub    string `json:"xpub"`
	Paymail string `json:"paymail"`
}

type Config struct {
	ClientOneURL         string
	ClientTwoURL         string
	ClientOneLeaderXPriv string
	ClientTwoLeaderXPriv string
}

type WalletResponse struct {
	Message string `json:"message"`
}

func SaveConfig(config *Config) error {
	envMap := map[string]string{
		ClientOneURLEnvVar:         config.ClientOneURL,
		ClientTwoURLEnvVar:         config.ClientTwoURL,
		ClientOneLeaderXPrivEnvVar: config.ClientOneLeaderXPriv,
		ClientTwoLeaderXPrivEnvVar: config.ClientTwoLeaderXPriv,
	}

	err := godotenv.Write(envMap, ".env")
	if err != nil {
		return fmt.Errorf("error saving .env file: %w", err)
	}

	return nil
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(".env"); err != nil {
		return nil, fmt.Errorf("error loading .env file: %v", err)
	}

	return &Config{
		ClientOneURL:         os.Getenv(ClientOneURLEnvVar),
		ClientTwoURL:         os.Getenv(ClientTwoURLEnvVar),
		ClientOneLeaderXPriv: os.Getenv(ClientOneLeaderXPrivEnvVar),
		ClientTwoLeaderXPriv: os.Getenv(ClientTwoLeaderXPrivEnvVar),
	}, nil
}

func IsSPVWalletRunning(url string) bool {
	url = AddPrefixIfNeeded(url)
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		return false
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return false
	}

	var walletResp WalletResponse
	if err := json.Unmarshal(body, &walletResp); err != nil {
		fmt.Println("Error parsing response JSON:", err)
		return false
	}

	return walletResp.Message == spvWalletIndexResponse
}

func AddPrefixIfNeeded(url string) string {
	if !strings.HasPrefix(url, domainPrefix) {
		return domainPrefix + url
	}
	return url
}

func GetSharedConfig(xpub string) (*models.SharedConfig, error) {
	req, err := http.NewRequest(http.MethodGet, AddPrefixIfNeeded(domainLocalHost)+domainSuffixSharedConfig, nil)
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

func PromptUserAndCheck(question string) int {
	reader := bufio.NewReader(os.Stdin)
	var response string
	var checkResult int

	for {
		fmt.Println(question)
		response, _ = reader.ReadString('\n')
		response = strings.TrimSpace(response)

		checkResult = CheckResponse(response)
		if checkResult != -1 {
			break
		}
		fmt.Println("Invalid response. Please answer y/yes or n/no.")
	}

	return checkResult
}

func CheckResponse(response string) int {
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

func PreparePaymail(paymailAlias string, domain string) string {
	return paymailAlias + atSign + domain
}

func CreateUser(paymail string, config *Config) (*User, error) {
	keys, err := xpriv.Generate()
	if err != nil {
		return nil, err
	}

	user := &User{
		XPriv:   keys.XPriv(),
		XPub:    keys.XPub().String(),
		Paymail: PreparePaymail(paymail, config.ClientOneURL),
	}

	adminClient := walletclient.NewWithAdminKey(AddPrefixIfNeeded(domainLocalHost), adminXPriv)
	ctx := context.Background()

	if err := adminClient.AdminNewXpub(ctx, user.XPub, map[string]any{"some_metadata": "remove"}); err != nil {
		fmt.Println("AdminNewXpub failed with status code:", err)
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
	user.Paymail = PreparePaymail(createPaymailRes.Alias, createPaymailRes.Domain)
	return user, nil
}

func UseUserFromEnv(config *Config, paymailAlias string) (*User, error) {
	keys, err := xpriv.FromString(config.ClientOneLeaderXPriv)
	if err != nil {
		return nil, fmt.Errorf("error parsing xpriv: %v", err)
	}
	return &User{
		XPriv:   keys.XPriv(),
		XPub:    keys.XPub().String(),
		Paymail: PreparePaymail(paymailAlias, config.ClientOneURL),
	}, nil
}

func DeleteUser(paymail string, config *Config) error {
	paymail = PreparePaymail(paymail, config.ClientOneURL)
	adminClient := walletclient.NewWithAdminKey(AddPrefixIfNeeded(domainLocalHost), adminXPriv)
	ctx := context.Background()
	err := adminClient.AdminDeletePaymail(ctx, paymail)
	if err != nil {
		return err
	}
	return nil
}

func GetValidXPriv() string {
	for {
		xpriv := PromptUser("Enter xpriv: ")
		if strings.HasPrefix(xpriv, "xprv") {
			return xpriv
		}
		fmt.Println("Invalid xpriv. Please enter a valid xpriv")
	}
}

func PromptUser(question string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println(question)
	response, _ := reader.ReadString('\n')
	return strings.TrimSpace(response)
}

func GetValidURL() string {
	for {
		url := PromptUser("Enter master instance URL with prefix: ")
		if IsValidURL(url) {
			return url
		}
		fmt.Println("Invalid URL. Please enter a valid URL with http/https prefix")
	}
}

func IsValidURL(rawURL string) bool {
	return explicitHTTPURLRegex.MatchString(rawURL)
}

func CheckBalance(domain, xpriv string) int {
	client := walletclient.NewWithXPriv(AddPrefixIfNeeded(domain), xpriv)
	ctx := context.Background()

	xpubInfo, err := client.GetXPub(ctx)
	if err != nil {
		fmt.Println("Error getting xpub info:", err)
		os.Exit(1)
	}
	return int(xpubInfo.CurrentBalance)
}
