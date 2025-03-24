//nolint:nolintlint,revive
package manualtests

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/api/manualtests/client"
	"github.com/go-viper/mapstructure/v2"
	"github.com/joomcode/errorx"
	"github.com/samber/lo"
	"github.com/spf13/viper"
)

const stateFileName = "state.yaml"
const notConfiguredPaymailDomain = "replace_me"
const notConfiguredFaucetURL = "https://replace_me.localhost"
const notConfiguredXprv = "xprvreplaceme"
const notConfiguredRecipientID = "replace_me"
const notConfiguredExternalPaymail = "replace.me@locahost"
const notConfiguredRegressionPaymail = "regression_tests_funds@replace.me.locahost"

var StateError = errorx.NewType(errorx.CommonErrors, "state_error")
var NotFound = errorx.NewType(errorx.CommonErrors, "not_found")

type State struct {
	Domain                string `mapstructure:"domain"     yaml:"domain"`
	ServerURL             string `mapstructure:"server_url" yaml:"server_url"`
	Faucet                Faucet `mapstructure:"faucet"     yaml:"faucet"`
	AdminXpub             string `mapstructure:"admin"      yaml:"admin"`
	User                  User
	Payment               Payment `mapstructure:"payment" yaml:"payment"`
	OldUsers              []*User `mapstructure:"zzz_old_users" yaml:"zzz_old_users"`
	configFilePath        string
	configFileJustCreated bool
}

func NewState() *State {
	state := &State{
		Domain:    notConfiguredPaymailDomain,
		ServerURL: "http://localhost:3003",
		AdminXpub: "xpub661MyMwAqRbcFgfmdkPgE2m5UjHXu9dj124DbaGLSjaqVESTWfCD4VuNmEbVPkbYLCkykwVZvmA8Pbf8884TQr1FgdG2nPoHR8aB36YdDQh",
	}

	user := User{
		state: state,
	}

	faucet := Faucet{
		URL:                notConfiguredFaucetURL,
		Xpriv:              notConfiguredXprv,
		DefaultTopUpAmount: 11,
		state:              state,
	}

	payment := Payment{
		Amount:                 10,
		RecipientID:            notConfiguredRecipientID,
		ExternalPaymail:        notConfiguredExternalPaymail,
		RegressionTestsPaymail: notConfiguredRegressionPaymail,
		state:                  state,
	}

	state.User = user
	state.Faucet = faucet
	state.Payment = payment

	return state
}

func (s *State) Load() error {
	s.envConfig()
	err := s.setDefaultViperState()
	if err != nil {
		return err
	}

	err = s.loadFromFile()
	if err != nil {
		return err
	}

	err = s.unmarshal()
	if err != nil {
		return err
	}

	err = s.validate()
	if err != nil {
		return err
	}

	err = s.init()
	if err != nil {
		return err
	}
	return nil
}

func (s *State) Save() error {
	err := s.updateViperState()
	if err != nil {
		return err
	}

	err = viper.WriteConfig()
	if err != nil {
		return StateError.Wrap(err, "could not write state file")
	}

	return nil
}

func (s *State) SaveOnSuccess(res Result) error {
	if res.StatusCode() < 200 && res.StatusCode() >= 300 {
		return nil
	}

	return s.Save()
}

func (s *State) NewUser(xpriv string, xpub string) (*User, error) {
	if !s.User.IsEmpty() {
		oldUser := s.User
		s.OldUsers = append(s.OldUsers, &oldUser)
	}
	err := s.User.new(xpriv, xpub)
	if err != nil {
		return nil, StateError.Wrap(err, "could not create new user")
	}

	return &s.User, nil
}

func (s *State) envConfig() {
	viper.SetEnvPrefix("MAN_TEST")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
}

func (s *State) setDefaultViperState() error {
	defaultsMap := make(map[string]interface{})
	if err := mapstructure.Decode(s, &defaultsMap); err != nil {
		return StateError.Wrap(err, "error occurred while setting defaults")
	}

	for key, value := range defaultsMap {
		viper.SetDefault(key, value)
	}

	return nil
}

func (s *State) updateViperState() error {
	defaultsMap := make(map[string]interface{})
	if err := mapstructure.Decode(s, &defaultsMap); err != nil {
		return StateError.Wrap(err, "error occurred while setting defaults")
	}

	for key, value := range defaultsMap {
		viper.Set(key, value)
	}

	return nil
}

func (s *State) loadFromFile() error {
	configFilePath, err := s.prepareConfigFile()
	if err != nil {
		return StateError.Wrap(err, "could not prepare state file")
	}

	s.configFilePath = configFilePath
	viper.SetConfigFile(configFilePath)

	if err = viper.ReadInConfig(); err != nil {
		return StateError.Wrap(err, "could not read state file")
	}

	if s.configFileJustCreated {
		err = viper.WriteConfig()
		if err != nil {
			return StateError.Wrap(err, "could not initialise state file")
		}
	}
	return nil
}

func (s *State) prepareConfigFile() (string, error) {
	// Get the absolute path of the current source file
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		return "", StateError.New("failed to get current file path")
	}

	dir := filepath.Dir(currentFile)

	stateFilePath := filepath.Join(dir, stateFileName)

	// Create the file if doesn't exist
	if _, err := os.Stat(stateFilePath); os.IsNotExist(err) {
		file, err := os.Create(stateFilePath) //nolint:gosec // this is only for testing purpose
		if err != nil {
			return "", StateError.Wrap(err, "could not create state file")
		}
		defer func(file *os.File) {
			_ = file.Close()
		}(file)
		s.configFileJustCreated = true
	}

	return stateFilePath, nil
}

func (s *State) unmarshal() error {
	if err := viper.Unmarshal(s); err != nil {
		return StateError.Wrap(err, "error when unmarshalling state to App Config")
	}
	return nil
}

func (s *State) init() error {
	err := s.User.init()
	if err != nil {
		return err
	}

	return nil
}

func (s *State) validate() error {
	if s.Domain == notConfiguredPaymailDomain {
		return StateError.New("Please configure (adjust) the domain in file://%s", s.configFilePath)
	}
	return nil
}

func (s *State) CurrentUser() *User {
	return &s.User
}

func (s *State) LatestDataID() string {
	dataOutpoints := s.User.DataOutpoints
	if len(dataOutpoints) == 0 {
		return ""
	}
	return dataOutpoints[len(dataOutpoints)-1]
}

func (s *State) UseUserWithID(userID string) error {
	if s.User.ID == userID {
		return nil
	}
	founded, err := s.GetOldUserById(userID)
	if err != nil {
		return errorx.Decorate(err, "failed to get user to switch to")
	}

	err = founded.init()
	if err != nil {
		return err
	}

	oldUser := s.User

	s.User = *founded

	s.OldUsers = lo.Filter(s.OldUsers, func(user *User, _ int) bool {
		return user.ID != founded.ID
	})
	s.OldUsers = append(s.OldUsers, &oldUser)

	return nil
}

func (s *State) GetUserById(userID string) (*User, error) {
	if s.User.ID == userID {
		return &s.User, nil
	}

	return s.GetOldUserById(userID)
}

func (s *State) GetOldUserById(userID string) (*User, error) {
	if userID == "" {
		return nil, StateError.New("You must provide ID to search for")
	}

	if len(s.OldUsers) == 0 {
		return nil, StateError.New("no old users to search for user")
	}

	founded, success := lo.Find(s.OldUsers, func(user *User) bool {
		return user.ID == userID
	})

	if !success {
		return nil, NotFound.New("user with ID %s not found", userID)
	}

	return founded, nil

}

func (s *State) UnlockOutlineHex(outline *client.ResponsesCreateTransactionOutlineSuccess) (string, error) {
	hex, err := NewTxSigner(s).UsingAnnotations(outline.Annotations.Inputs).UnlockToHex(string(outline.Format), outline.Hex)
	if err != nil {
		return "", errorx.Decorate(err, "failed to unlock outline hex")
	}
	return hex, nil
}

func (s *State) UnlockOutline(outline *client.ResponsesCreateTransactionOutlineSuccess) (*sdk.Transaction, error) {
	tx, err := NewTxSigner(s).UsingAnnotations(outline.Annotations.Inputs).Unlock(string(outline.Format), outline.Hex)
	if err != nil {
		return nil, errorx.Decorate(err, "failed to unlock outline hex")
	}
	return tx, nil
}

func (s *State) SaveUserDataOutpoint(txID string, vout int) error {
	s.CurrentUser().AddDataOutpoint(txID, vout)
	err := s.Save()
	if err != nil {
		return err
	}
	return nil
}

func (s *State) CleanupOldUsersByTag(tag string) {
	s.OldUsers = lo.Filter(s.OldUsers, func(user *User, _ int) bool {
		return !lo.Contains(user.Tags, tag)
	})
}

func (s *State) GetLastOldUserFromState() (User, error) {
	if len(s.OldUsers) == 0 {
		return User{}, NotFound.New("no old users")
	}
	return *s.OldUsers[0], nil
}
