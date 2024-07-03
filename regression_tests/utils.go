package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	walletclient "github.com/bitcoin-sv/spv-wallet-go-client"
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
	return nil, nil
}

func UpdateConfigWithUserKeys(config *Config, user *User) {
	config.ClientOneLeaderXPriv = user.XPriv
	config.ClientTwoLeaderXPriv = user.XPriv
}

func UseUserFromEnv(config *Config, paymailAlias string) (*User, error) {
	return nil, nil
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
