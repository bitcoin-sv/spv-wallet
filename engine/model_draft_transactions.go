package engine

import (
	"context"
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/bitcoin-sv/go-paymail"
	compat "github.com/bitcoin-sv/go-sdk/compat/bip32"
	ec "github.com/bitcoin-sv/go-sdk/primitives/ec"
	"github.com/bitcoin-sv/go-sdk/script"
	trx "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/go-sdk/transaction/template/p2pkh"
	"github.com/bitcoin-sv/spv-wallet/conv"
	"github.com/bitcoin-sv/spv-wallet/engine/datastore"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/utils"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// DraftTransaction is an object representing the draft BitCoin transaction prior to the final transaction
//
// Gorm related models & indexes: https://gorm.io/docs/models.html - https://gorm.io/docs/indexes.html
type DraftTransaction struct {
	// Base model
	Model

	// Standard transaction model base fields
	TransactionBase

	// Model specific fields
	XpubID        string            `json:"xpub_id" toml:"xpub_id" yaml:"xpub_id" gorm:"<-:create;type:char(64);index;comment:This is the related xPub"`
	ExpiresAt     time.Time         `json:"expires_at" toml:"expires_at" yaml:"expires_at" gorm:"<-:create;comment:Time when the draft expires"`
	Configuration TransactionConfig `json:"configuration" toml:"configuration" yaml:"configuration" gorm:"<-;type:text;comment:This is the configuration struct in JSON"`
	Status        DraftStatus       `json:"status" toml:"status" yaml:"status" gorm:"<-;type:varchar(10);index;comment:This is the status of the draft"`
	FinalTxID     string            `json:"final_tx_id,omitempty" toml:"final_tx_id" yaml:"final_tx_id" gorm:"<-;type:char(64);index;comment:This is the final tx ID"`
}

// newDraftTransaction will start a new draft tx
func newDraftTransaction(rawXpubKey string, config *TransactionConfig, opts ...ModelOps) (*DraftTransaction, error) {
	// Random GUID
	id, _ := utils.RandomHex(32)

	// Set the expires time (default)
	expiresAt := time.Now().UTC().Add(defaultDraftTxExpiresIn)
	if config.ExpiresIn > 0 {
		expiresAt = time.Now().UTC().Add(config.ExpiresIn)
	}

	// Start the model
	draft := &DraftTransaction{
		Configuration:   *config,
		ExpiresAt:       expiresAt,
		Status:          DraftStatusDraft,
		TransactionBase: TransactionBase{ID: id},
		XpubID:          utils.Hash(rawXpubKey),
		Model: *NewBaseModel(
			ModelDraftTransaction,
			append(opts, WithXPub(rawXpubKey))...,
		),
	}

	if config.FeeUnit == nil {
		unit := draft.Client().FeeUnit()
		draft.Configuration.FeeUnit = &unit
	}

	err := draft.createTransactionHex(context.Background())
	if err != nil {
		return nil, err
	}
	return draft, nil
}

// getDraftTransactionID will get the draft transaction with the given conditions
func getDraftTransactionID(ctx context.Context, xPubID, id string,
	opts ...ModelOps,
) (*DraftTransaction, error) {
	// Get the record
	conditions := map[string]interface{}{
		xPubIDField: xPubID,
		idField:     id,
	}
	if len(xPubID) == 0 {
		conditions = map[string]interface{}{
			idField: id,
		}
	}

	draftTransaction := &DraftTransaction{Model: *NewBaseModel(
		ModelDraftTransaction,
		opts...,
	)}
	if err := Get(ctx, draftTransaction, conditions, false, defaultDatabaseReadTimeout, true); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return draftTransaction, nil
}

// GetModelName will get the name of the current model
func (m *DraftTransaction) GetModelName() string {
	return ModelDraftTransaction.String()
}

// GetModelTableName will get the db table name of the current model
func (m *DraftTransaction) GetModelTableName() string {
	return tableDraftTransactions
}

// Save will save the model into the Datastore
func (m *DraftTransaction) Save(ctx context.Context) (err error) {
	if err = Save(ctx, m); err != nil {

		m.Client().Logger().Error().
			Str("draftTxID", m.GetID()).
			Msgf("save tx error: %s", err.Error())

		// todo: run in a go routine?
		// un-reserve the utxos
		if utxoErr := unReserveUtxos(
			ctx, m.XpubID, m.ID, m.GetOptions(false)...,
		); utxoErr != nil {
			err = spverrors.Wrapf(err, utxoErr.Error())
		}
	}
	return
}

// GetID will get the model ID
func (m *DraftTransaction) GetID() string {
	return m.ID
}

// processConfigOutputs will process all the outputs,
// doing any lookups and creating locking scripts
func (m *DraftTransaction) processConfigOutputs(ctx context.Context) error {
	// Get the client
	c := m.Client()
	// Get sender's paymail
	paymailFrom := c.GetPaymailConfig().DefaultFromPaymail
	conditions := map[string]interface{}{
		xPubIDField: m.XpubID,
	}

	// Get the sender's paymail from the metadata, this help when sender has multiple paymails
	senderPaymail, ok := m.Metadata["sender"].(string)
	if ok {
		alias, _, address := paymail.SanitizePaymail(senderPaymail)
		if address != "" {
			conditions["alias"] = alias
		}
	}

	paymails, err := c.GetPaymailAddressesByXPubID(ctx, m.XpubID, nil, conditions, nil)
	if err == nil && len(paymails) != 0 {
		paymailFrom = fmt.Sprintf("%s@%s", paymails[0].Alias, paymails[0].Domain)
	}

	paymailService := c.PaymailService()

	// Special case where we are sending all funds to a single (address, paymail, handle)
	if m.Configuration.SendAllTo != nil {
		outputs := m.Configuration.Outputs

		m.Configuration.SendAllTo.UseForChange = true
		m.Configuration.SendAllTo.Satoshis = 0
		m.Configuration.Outputs = []*TransactionOutput{m.Configuration.SendAllTo}

		if err := m.Configuration.Outputs[0].processOutput(ctx, paymailService, paymailFrom, false); err != nil {
			return err
		}

		// re-add the other outputs we had before
		for _, output := range outputs {
			output.UseForChange = false // make sure we do not add change to this output
			if err := output.processOutput(ctx, paymailService, paymailFrom, true); err != nil {
				return err
			}
			m.Configuration.Outputs = append(m.Configuration.Outputs, output)
		}
	} else {
		// Loop all outputs and process
		for index := range m.Configuration.Outputs {

			// Start the output script slice
			if m.Configuration.Outputs[index].Scripts == nil {
				m.Configuration.Outputs[index].Scripts = make([]*ScriptOutput, 0)
			}

			// Process the outputs
			if err := m.Configuration.Outputs[index].processOutput(ctx, paymailService, paymailFrom, true); err != nil {
				return err
			}
		}
	}

	return nil
}

// createTransactionHex will create the transaction with the given inputs and outputs
func (m *DraftTransaction) createTransactionHex(ctx context.Context) (err error) {
	// Check that we have outputs
	if len(m.Configuration.Outputs) == 0 && m.Configuration.SendAllTo == nil {
		return spverrors.ErrMissingTransactionOutputs
	}

	// Get the total satoshis needed to make this transaction
	satoshisNeeded := m.getTotalSatoshis()

	// Set opts
	opts := m.GetOptions(false)

	// Process the outputs first
	// if an error occurs in processing the outputs, we have at least not made any reservations yet
	if err = m.processConfigOutputs(ctx); err != nil {
		return
	}

	inputUtxos, satoshisReserved, err := m.prepareUtxos(ctx, opts, satoshisNeeded)
	if err != nil {
		return
	}

	// Start a new transaction from the reservedUtxos
	tx := trx.NewTransaction()
	if err = tx.AddInputsFromUTXOs(inputUtxos...); err != nil {
		return
	}

	// Estimate the fee for the transaction
	err = m.calculateAndSetFee(ctx, satoshisReserved, satoshisNeeded)
	if err != nil {
		return
	}

	// Add the outputs to the bt transaction
	if err = m.addOutputsToTx(tx); err != nil {
		return
	}

	if err = validateOutputsInputs(m.Configuration.Inputs, m.Configuration.Outputs, m.Configuration.Fee); err != nil {
		return
	}

	// Create the final hex (without signatures)
	m.Hex = tx.String()

	return
}

func (m *DraftTransaction) calculateAndSetFee(ctx context.Context, satoshisReserved uint64, satoshisNeeded uint64) error {
	fee := m.estimateFee(m.Configuration.FeeUnit, 0)
	if m.Configuration.SendAllTo != nil {
		if m.Configuration.Outputs[0].Satoshis <= dustLimit {
			return spverrors.ErrOutputValueTooLow
		}

		m.Configuration.Fee = fee
		m.Configuration.Outputs[0].Satoshis -= fee

		// subtract all the satoshis sent in other outputs
		for _, output := range m.Configuration.Outputs {
			if !output.UseForChange { // only normal outputs
				m.Configuration.Outputs[0].Satoshis -= output.Satoshis
			}
		}

		m.Configuration.Outputs[0].Scripts[0].Satoshis = m.Configuration.Outputs[0].Satoshis
		return nil
	}

	if satoshisReserved < satoshisNeeded+fee {
		return spverrors.ErrNotEnoughUtxos
	}

	// if we have a remainder, add that to an output to our own wallet address
	satoshisChange := satoshisReserved - satoshisNeeded - fee
	m.Configuration.Fee = fee
	if satoshisChange > 0 {
		var newFee uint64
		newFee, err := m.setChangeDestination(
			ctx, satoshisChange, fee,
		)
		if err != nil {
			return err
		}
		m.Configuration.Fee = newFee
	}

	return nil
}

func (m *DraftTransaction) prepareUtxos(ctx context.Context, opts []ModelOps, satoshisNeeded uint64) ([]*trx.UTXO, uint64, error) {
	if m.Configuration.SendAllTo != nil {
		inputUtxos, satoshisReserved, err := m.prepareSendAllToUtxos(ctx, opts)
		if err != nil {
			return nil, 0, err
		}
		return inputUtxos, satoshisReserved, nil
	}

	inputUtxos, satoshisReserved, err := m.prepareSeparateUtxos(ctx, opts, satoshisNeeded)
	if err != nil {
		return nil, 0, err
	}
	return inputUtxos, satoshisReserved, nil
}

// prepareSeparateUtxos will get user's utxos which will have the required amount of satoshi and then reserve and process them.
func (m *DraftTransaction) prepareSeparateUtxos(ctx context.Context, opts []ModelOps, satoshisNeeded uint64) ([]*trx.UTXO, uint64, error) {
	// we can only include separate utxos (like tokens) when not using SendAllTo
	var includeUtxoSatoshis uint64
	var err error
	if m.Configuration.IncludeUtxos != nil {
		includeUtxoSatoshis, err = m.addIncludeUtxos(ctx)
		if err != nil {
			return nil, 0, err
		}
	}

	// Reserve and Get utxos for the transaction
	var reservedUtxos []*Utxo
	//TODO: Fixme in new transaction-flow
	feePerByte := float64(m.Configuration.FeeUnit.Satoshis) / float64(m.Configuration.FeeUnit.Bytes)

	reserveSatoshis := satoshisNeeded + m.estimateFee(m.Configuration.FeeUnit, 0)
	if reserveSatoshis <= dustLimit && !m.containsOpReturn() {
		m.client.Logger().Error().
			Str("txID", m.GetID()).
			Msg("amount of satoshis to send less than the dust limit")
		return nil, 0, err
	}
	if reservedUtxos, err = reserveUtxos(
		ctx, m.XpubID, m.ID, reserveSatoshis, feePerByte, m.Configuration.FromUtxos, opts...,
	); err != nil {
		return nil, 0, err
	}

	// Get the inputUtxos (in bt.UTXO format) and the total amount of satoshis from the utxos
	inputUtxos, satoshisReserved, err := m.getInputsFromUtxos(reservedUtxos)
	if err != nil {
		return nil, 0, err
	}

	// add the satoshis from the utxos we forcibly included to the total input sats
	satoshisReserved += includeUtxoSatoshis

	// Reserve the utxos
	if err = m.processUtxos(
		ctx, reservedUtxos,
	); err != nil {
		return nil, 0, err
	}
	return inputUtxos, satoshisReserved, nil
}

// prepareSendAllToUtxos will reserve and process all the user's utxos which will be sent to one address
func (m *DraftTransaction) prepareSendAllToUtxos(ctx context.Context, opts []ModelOps) ([]*trx.UTXO, uint64, error) {
	// todo should all utxos be sent to the SendAllTo address, not only the p2pkhs?
	spendableUtxos, err := getSpendableUtxos(
		ctx, m.XpubID, utils.ScriptTypePubKeyHash, nil, m.Configuration.FromUtxos, opts...,
	)
	if err != nil {
		return nil, 0, err
	}
	for _, utxo := range spendableUtxos {
		// Reserve the utxos
		utxo.DraftID.Valid = true
		utxo.DraftID.String = m.ID
		utxo.ReservedAt.Valid = true
		utxo.ReservedAt.Time = time.Now().UTC()

		// Save the UTXO
		if err = utxo.Save(ctx); err != nil {
			return nil, 0, err
		}

		m.Configuration.Outputs[0].Satoshis += utxo.Satoshis
	}

	// Get the inputUtxos (in bt.UTXO format) and the total amount of satoshis from the utxos
	inputUtxos, satoshisReserved, err := m.getInputsFromUtxos(
		spendableUtxos,
	)
	if err != nil {
		return nil, 0, err
	}

	if err = m.processUtxos(
		ctx, spendableUtxos,
	); err != nil {
		return nil, 0, err
	}
	return inputUtxos, satoshisReserved, nil
}

func validateOutputsInputs(inputs []*TransactionInput, outputs []*TransactionOutput, fee uint64) error {
	usedUtxos := make([]string, 0)
	inputValue := uint64(0)
	outputValue := uint64(0)

	for _, input := range inputs {
		if utils.StringInSlice(input.Utxo.ID, usedUtxos) {
			return spverrors.ErrDuplicateUTXOs
		}
		usedUtxos = append(usedUtxos, input.Utxo.ID)
		inputValue += input.Satoshis
	}

	for _, output := range outputs {
		outputValue += output.Satoshis
	}

	if inputValue < outputValue {
		return spverrors.ErrOutputValueTooHigh
	}

	if inputValue-outputValue != fee {
		return spverrors.ErrTransactionFeeInvalid
	}
	return nil
}

// addIncludeUtxos will add the included utxos
func (m *DraftTransaction) addIncludeUtxos(ctx context.Context) (uint64, error) {
	// Whatever utxos are selected, the IncludeUtxos should be added to the transaction
	// This can be used to add for instance tokens where fees need to be paid from other utxos
	// The satoshis of these inputs are not added to the reserved satoshis. If these inputs contain satoshis
	// that will be added to the total inputs and handled with the change addresses.
	includeUtxos := make([]*Utxo, 0)
	opts := m.GetOptions(false)
	var includeUtxoSatoshis uint64
	for _, utxo := range m.Configuration.IncludeUtxos {
		utxoModel, err := getUtxo(ctx, utxo.TransactionID, utxo.OutputIndex, opts...)
		if err != nil {
			return 0, err
		} else if utxoModel == nil {
			return 0, spverrors.ErrCouldNotFindUtxo
		}
		includeUtxos = append(includeUtxos, utxoModel)
		includeUtxoSatoshis += utxoModel.Satoshis
	}
	return includeUtxoSatoshis, m.processUtxos(ctx, includeUtxos)
}

// processUtxos will process the utxos
func (m *DraftTransaction) processUtxos(ctx context.Context, utxos []*Utxo) error {
	// Get destinations
	opts := m.GetOptions(false)
	for _, utxo := range utxos {
		lockingScript := utils.GetDestinationLockingScript(utxo.ScriptPubKey)
		destination, err := getDestinationWithCache(
			ctx, m.Client(), "", "", lockingScript, opts...,
		)
		if err != nil {
			return err
		}
		if destination == nil {
			return spverrors.ErrCouldNotFindDestination
		}
		m.Configuration.Inputs = append(
			m.Configuration.Inputs, &TransactionInput{
				Utxo:        *utxo,
				Destination: *destination,
			})
	}

	return nil
}

// estimateSize will loop the inputs and outputs and estimate the size of the transaction
func (m *DraftTransaction) estimateSize() uint64 {
	size := defaultOverheadSize // version + nLockTime

	inputSize := trx.VarInt(len(m.Configuration.Inputs))

	value, err := conv.IntToUint64(inputSize.Length())
	if err != nil {
		m.client.Logger().Error().Msg(err.Error())
		return 0
	}
	size += value

	for _, input := range m.Configuration.Inputs {
		size += utils.GetInputSizeForType(input.Type)
	}

	outputSize := trx.VarInt(len(m.Configuration.Outputs))
	value, err = conv.IntToUint64(outputSize.Length())
	if err != nil {
		m.client.Logger().Error().Msg(err.Error())
		return 0
	}
	size += value
	for _, output := range m.Configuration.Outputs {
		for _, s := range output.Scripts {
			size += utils.GetOutputSize(s.Script)
		}
	}

	return size
}

// estimateFee will loop the inputs and outputs and estimate the required fee
func (m *DraftTransaction) estimateFee(unit *bsv.FeeUnit, addToSize uint64) uint64 {
	size := m.estimateSize() + addToSize
	feeEstimate := float64(size) * (float64(unit.Satoshis) / float64(unit.Bytes))
	return uint64(math.Ceil(feeEstimate))
}

// addOutputs will add the given outputs to the SDK Transaction
func (m *DraftTransaction) addOutputsToTx(tx *trx.Transaction) (err error) {
	var s *script.Script
	for _, output := range m.Configuration.Outputs {
		for _, sc := range output.Scripts {
			if s, err = script.NewFromHex(
				sc.Script,
			); err != nil {
				return
			}

			scriptType := sc.ScriptType
			if scriptType == "" {
				scriptType = utils.GetDestinationType(sc.Script)
			}

			if scriptType == utils.ScriptTypeNullData {
				// op_return output - only one allowed to have 0 satoshi value ???
				if sc.Satoshis > 0 {
					return spverrors.ErrInvalidOpReturnOutput
				}

				tx.AddOutput(&trx.TransactionOutput{
					LockingScript: s,
					Satoshis:      0,
				})
			} else if scriptType == utils.ScriptTypePubKeyHash {
				// sending to a p2pkh
				if sc.Satoshis == 0 {
					return spverrors.ErrOutputValueTooLow
				}

				tx.AddOutput(
					&trx.TransactionOutput{
						LockingScript: s,
						Satoshis:      sc.Satoshis,
					})
			} else {
				// add non-standard output script
				tx.AddOutput(&trx.TransactionOutput{
					LockingScript: s,
					Satoshis:      sc.Satoshis,
				})
			}
		}
	}
	return
}

// setChangeDestination will make a new change destination
func (m *DraftTransaction) setChangeDestination(ctx context.Context, satoshisChange uint64, fee uint64) (uint64, error) {
	m.Configuration.ChangeSatoshis = satoshisChange

	useExistingOutputsForChange := make([]int, 0)
	for index := range m.Configuration.Outputs {
		if m.Configuration.Outputs[index].UseForChange {
			useExistingOutputsForChange = append(useExistingOutputsForChange, index)
		}
	}

	newFee := fee
	if len(useExistingOutputsForChange) > 0 {
		// reset destinations if set
		m.Configuration.ChangeDestinationsStrategy = ChangeStrategyDefault
		m.Configuration.ChangeDestinations = nil

		numberOfExistingOutputs := uint64(len(useExistingOutputsForChange))
		changePerOutput := uint64(float64(satoshisChange) / float64(numberOfExistingOutputs))
		remainderOutput := satoshisChange - (changePerOutput * numberOfExistingOutputs)
		for _, outputIndex := range useExistingOutputsForChange {
			m.Configuration.Outputs[outputIndex].Satoshis += changePerOutput + remainderOutput
			remainderOutput = 0 // reset remainder to 0 for other outputs
		}
	} else {
		numberOfDestinations := m.Configuration.ChangeNumberOfDestinations
		if numberOfDestinations <= 0 {
			numberOfDestinations = 1 // todo get from config
		}
		minimumSatoshis := m.Configuration.ChangeMinimumSatoshis
		if minimumSatoshis <= 0 { // todo: protect against un-spendable amount? less than fee to miner for min tx?
			minimumSatoshis = 1250 // todo get from config
		}

		if float64(satoshisChange)/float64(numberOfDestinations) < float64(minimumSatoshis) {
			// we cannot split our change to the number of destinations given, re-calc
			numberOfDestinations = 1
		}

		// Check if numberOfDestinations is negative or too large before conversion
		if numberOfDestinations < 0 {
			return fee, fmt.Errorf("invalid number of destinations: %d", numberOfDestinations)
		}

		newFee = m.estimateFee(m.Configuration.FeeUnit, uint64(numberOfDestinations)*changeOutputSize)
		satoshisChange -= newFee - fee
		m.Configuration.ChangeSatoshis = satoshisChange

		if m.Configuration.ChangeDestinations == nil {
			if err := m.setChangeDestinations(
				ctx, numberOfDestinations,
			); err != nil {
				return fee, err
			}
		}

		changeSatoshis, err := m.getChangeSatoshis(satoshisChange)
		if err != nil {
			return fee, err
		}

		for _, destination := range m.Configuration.ChangeDestinations {
			m.Configuration.Outputs = append(m.Configuration.Outputs, &TransactionOutput{
				To: destination.Address,
				Scripts: []*ScriptOutput{{
					Address:    destination.Address,
					Satoshis:   changeSatoshis[destination.LockingScript],
					Script:     destination.LockingScript,
					ScriptType: utils.ScriptTypePubKeyHash,
				}},
				Satoshis: changeSatoshis[destination.LockingScript],
			})
		}
	}

	return newFee, nil
}

// split the change satoshis amongst the change destinations according to the strategy given in config
func (m *DraftTransaction) getChangeSatoshis(satoshisChange uint64) (changeSatoshis map[string]uint64, err error) {
	changeSatoshis = make(map[string]uint64)
	var lastDestination string
	changeUsed := uint64(0)

	if m.Configuration.ChangeDestinationsStrategy == ChangeStrategyNominations {
		return nil, spverrors.ErrChangeStrategyNotImplemented
	} else if m.Configuration.ChangeDestinationsStrategy == ChangeStrategyRandom {
		nDestinations := float64(len(m.Configuration.ChangeDestinations))
		var a *big.Int
		for _, destination := range m.Configuration.ChangeDestinations {
			if a, err = rand.Int(
				rand.Reader, big.NewInt(math.MaxInt64),
			); err != nil {
				return
			}
			randomChange := (((float64(a.Int64()) / (1 << 63)) * 50) + 75) / 100
			changeForDestination := uint64(randomChange * float64(satoshisChange) / nDestinations)

			changeSatoshis[destination.LockingScript] = changeForDestination
			lastDestination = destination.LockingScript
			changeUsed += changeForDestination
		}
	} else {
		// default
		changePerDestination := uint64(float64(satoshisChange) / float64(len(m.Configuration.ChangeDestinations)))
		for _, destination := range m.Configuration.ChangeDestinations {
			changeSatoshis[destination.LockingScript] = changePerDestination
			lastDestination = destination.LockingScript
			changeUsed += changePerDestination
		}
	}

	// handle remainder
	changeSatoshis[lastDestination] += satoshisChange - changeUsed

	return
}

// setChangeDestinations will set the change destinations based on the number
func (m *DraftTransaction) setChangeDestinations(ctx context.Context, numberOfDestinations int) error {
	// Set the options
	opts := m.GetOptions(false)
	optsNew := append(opts, New())
	c := m.Client()

	var err error
	var xPub *Xpub
	var num uint32

	// Loop for each destination
	for i := 0; i < numberOfDestinations; i++ {
		if xPub, err = getXpubWithCache(
			ctx, c, m.rawXpubKey, "", opts...,
		); err != nil {
			return err
		} else if xPub == nil {
			return spverrors.ErrMissingFieldXpub
		}

		if num, err = xPub.incrementNextNum(
			ctx, utils.ChainInternal,
		); err != nil {
			return err
		}

		var destination *Destination
		if destination, err = newAddress(
			m.rawXpubKey, utils.ChainInternal, num, optsNew...,
		); err != nil {
			return err
		}

		destination.DraftID = m.ID
		if err = destination.Save(ctx); err != nil {
			return err
		}

		m.Configuration.ChangeDestinations = append(m.Configuration.ChangeDestinations, destination)
	}

	return nil
}

// getInputsFromUtxos this function transforms SPV Wallet utxos to SDK UTXOs
func (m *DraftTransaction) getInputsFromUtxos(reservedUtxos []*Utxo) ([]*trx.UTXO, uint64, error) {
	// transform to bt.utxo and check if we have enough
	inputUtxos := make([]*trx.UTXO, 0)
	satoshisReserved := uint64(0)
	var lockingScript *script.Script
	var err error
	for _, utxo := range reservedUtxos {

		if lockingScript, err = script.NewFromHex(
			utxo.ScriptPubKey,
		); err != nil {
			return nil, 0, spverrors.ErrInvalidLockingScript
		}

		utxo, err := trx.NewUTXO(
			utxo.TransactionID,
			utxo.OutputIndex,
			lockingScript.String(),
			utxo.Satoshis,
		)
		if err != nil {
			return nil, 0, spverrors.ErrFailedToCreateUTXO.Wrap(err)
		}

		inputUtxos = append(inputUtxos, utxo)
		satoshisReserved += utxo.Satoshis
	}

	return inputUtxos, satoshisReserved, nil
}

// getTotalSatoshis calculate the total satoshis of all outputs
func (m *DraftTransaction) getTotalSatoshis() (satoshis uint64) {
	for _, output := range m.Configuration.Outputs {
		satoshis += output.Satoshis
	}
	return
}

// BeforeCreating will fire before the model is being inserted into the Datastore
func (m *DraftTransaction) BeforeCreating(_ context.Context) (err error) {
	m.Client().Logger().Debug().
		Str("draftTxID", m.GetID()).
		Msgf("starting: %s BeforeCreating hook...", m.Name())

	m.Client().Logger().Debug().
		Str("draftTxID", m.GetID()).
		Msgf("end: %s BeforeCreating hook", m.Name())
	return
}

// AfterUpdated will fire after a successful update into the Datastore
func (m *DraftTransaction) AfterUpdated(ctx context.Context) error {
	m.Client().Logger().Debug().
		Str("draftTxID", m.GetID()).
		Msgf("starting: %s AfterUpdated hook...", m.Name())

	// todo: run these in go routines?

	// remove reservation from all utxos related to this draft transaction
	if m.Status == DraftStatusCanceled || m.Status == DraftStatusExpired {
		utxos, err := getUtxosByDraftID(
			ctx, m.ID,
			nil,
			m.GetOptions(false)...,
		)
		if err != nil {
			return err
		}
		for index := range utxos {
			utxos[index].DraftID.String = ""
			utxos[index].DraftID.Valid = false
			utxos[index].ReservedAt.Time = time.Time{}
			utxos[index].ReservedAt.Valid = false
			if err = utxos[index].Save(ctx); err != nil {
				return err
			}
		}
	}

	m.Client().Logger().Debug().
		Str("draftTxID", m.GetID()).
		Msgf("end: %s AfterUpdated hook", m.Name())
	return nil
}

// PostMigrate is called after the model is migrated
func (m *DraftTransaction) PostMigrate(client datastore.ClientInterface) error {
	err := client.IndexMetadata(client.GetTableName(tableDraftTransactions), metadataField)
	return spverrors.Wrapf(err, "failed to index metadata column on model %s", m.GetModelName())
}

// SignInputsWithKey will sign all the inputs using a key (string) (helper method)
func (m *DraftTransaction) SignInputsWithKey(xPrivKey string) (signedHex string, err error) {
	// Decode the xPriv using the key
	var xPriv *compat.ExtendedKey
	if xPriv, err = compat.NewKeyFromString(xPrivKey); err != nil {
		return
	}

	return m.SignInputs(xPriv)
}

// SignInputs will sign all the inputs using the given xPriv key
func (m *DraftTransaction) SignInputs(xPriv *compat.ExtendedKey) (signedHex string, err error) {
	// Start a bt draft transaction
	var txDraft *trx.Transaction
	if txDraft, err = trx.NewTransactionFromHex(m.Hex); err != nil {
		return
	}

	// Sign the inputs
	for index, input := range m.Configuration.Inputs {

		// Get the locking script
		var ls *script.Script
		if ls, err = script.NewFromHex(
			input.Destination.LockingScript,
		); err != nil {
			return
		}
		txDraft.Inputs[index].SetSourceTxOutput(&trx.TransactionOutput{
			Satoshis:      input.Satoshis,
			LockingScript: ls,
		})

		// Derive the child key (chain)
		var chainKey *compat.ExtendedKey
		if chainKey, err = xPriv.Child(
			input.Destination.Chain,
		); err != nil {
			return
		}

		// Derive the child key (num)
		var numKey *compat.ExtendedKey
		if numKey, err = chainKey.Child(
			input.Destination.Num,
		); err != nil {
			return
		}

		// Get the private key
		var privateKey *ec.PrivateKey
		if privateKey, err = compat.GetPrivateKeyFromHDKey(
			numKey,
		); err != nil {
			return
		}

		idx32, conversionError := conv.IntToUint32(index)
		if err != nil {
			return "", spverrors.Wrapf(conversionError, "failed to convert index %d to uint32", index)
		}
		var s *p2pkh.P2PKH
		if s, err = utils.GetUnlockingScript(
			txDraft, idx32, privateKey,
		); err != nil {
			return
		}

		// Insert the locking script
		if txDraft.Inputs[index] == nil {
			return "", spverrors.Newf("input with index %d not found in transaction draft %v", index, txDraft)
		}
		txDraft.Inputs[index].UnlockingScriptTemplate = s
	}

	// Return the signed hex
	err = txDraft.Sign()
	if err != nil {
		return "", spverrors.Wrapf(err, "failed to sign inputs on model %s", m.GetModelName())
	}

	signedHex = txDraft.String()
	return
}

func (m *DraftTransaction) containsOpReturn() bool {
	for _, output := range m.Configuration.Outputs {
		if output.OpReturn != nil {
			return true
		}
	}
	return false
}
