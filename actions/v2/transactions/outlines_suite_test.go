package transactions_test

import (
	"github.com/bitcoin-sv/spv-wallet/actions/v2/transactions/internal/testabilities"
	testengine "github.com/bitcoin-sv/spv-wallet/engine/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"github.com/stretchr/testify/suite"
)

type txOutlineTestcase struct {
	request          string
	responseTemplate string
	responseParams   map[string]any
	outValues        []bsv.Satoshis
}

type txOutlineSuite struct {
	suite.Suite
	txOutlineTestcase
	initialSatoshis bsv.Satoshis
	explicitFormat  string
}

func newTxOutlineSuite(testcase txOutlineTestcase, initialSatoshis bsv.Satoshis) *txOutlineSuite {
	return &txOutlineSuite{
		txOutlineTestcase: testcase,
		initialSatoshis:   initialSatoshis,
	}
}

func (s *txOutlineSuite) withExplicitFormat(format string) *txOutlineSuite {
	s.explicitFormat = format
	return s
}

func (s *txOutlineSuite) expectedFormat() string {
	if s.explicitFormat != "" {
		return s.explicitFormat
	}
	return "BEEF"
}

func (s *txOutlineSuite) Test() {
	// given:
	given, then := testabilities.New(s.T())
	cleanup := given.StartedSPVWalletWithConfiguration(testengine.WithV2())
	defer cleanup()

	// and:
	given.Faucet(fixtures.Sender).TopUp(s.initialSatoshis)

	// and:
	client := given.HttpClient().ForUser()

	// when:
	req := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(s.request)

	if s.explicitFormat != "" {
		req.SetQueryParam("format", s.explicitFormat)
	}

	res, _ := req.Post(transactionsOutlinesURL)

	// then:
	thenResponse := then.Response(res)

	thenResponse.IsOK().
		WithJSONMatching(s.responseTemplate, given.OutlineResponseContext(s.expectedFormat(), s.responseParams))

	thenResponse.ContainsValidTransaction(s.expectedFormat()).
		WithOutputValues(s.outValues...)
}
