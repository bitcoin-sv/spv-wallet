package utils

import (
	"encoding/hex"
	"regexp"

	"github.com/bitcoinschema/go-bitcoin/v2"
	bscript2 "github.com/libsv/go-bt/v2/bscript"
)

const (
	// ScriptTypePubKey alias from bscript
	ScriptTypePubKey = bscript2.ScriptTypePubKey

	// ScriptTypePubKeyHash alias from bscript
	ScriptTypePubKeyHash = bscript2.ScriptTypePubKeyHash

	// ScriptTypeNullData alias from bscript
	ScriptTypeNullData = bscript2.ScriptTypeNullData

	// ScriptTypeMultiSig alias from bscript
	ScriptTypeMultiSig = bscript2.ScriptTypeMultiSig

	// ScriptTypeNonStandard alias from bscript
	ScriptTypeNonStandard = bscript2.ScriptTypeNonStandard

	// ScriptHashType is the type for the deprecated script hash
	ScriptHashType = "scripthash"

	// ScriptMetanet is the type for a metanet transaction
	ScriptMetanet = "metanet"

	// ScriptTypeTokenStas is the type for a STAS output
	ScriptTypeTokenStas = "token_stas"

	// ScriptTypeTokenSensible is the type for a Sensible output
	ScriptTypeTokenSensible = "token_sensible" // 73656e7369626c65
)

var (
	// HashPuzzleRegexpString regex string of a hash puzzle
	HashPuzzleRegexpString = `a914[\da-f]{40}8876a9[\da-f]{40}88ac`

	// HashPuzzleRegexp regex string of a hash puzzle
	HashPuzzleRegexp = `^a914[\da-f]{40}8876a9[\da-f]{40}88ac`

	// RunRegexpString regex string of run protocol in op_return
	RunRegexpString = `006a0372756e`

	// RunRegexp regex of run protocol in op_return
	RunRegexp = `^006a0372756e`

	// P2PKHRegexpString regex string of P2PKH
	P2PKHRegexpString = `76a914[\da-f]{40}88ac`

	// P2PKHRegexp OP_DUP OP_HASH160 [pubkey hash] OP_EQUALVERIFY OP_CHECKSIG
	P2PKHRegexp, _ = regexp.Compile(`^` + P2PKHRegexpString + `$`) // exact match

	// P2PKHSubstringRegexp substring of OP_DUP OP_HASH160 [pubkey hash] OP_EQUALVERIFY OP_CHECKSIG
	P2PKHSubstringRegexp, _ = regexp.Compile(P2PKHRegexpString)

	// P2SHRegexp OP_HASH160 [hash] OP_EQUAL
	P2SHRegexp, _ = regexp.Compile(`^a914[\da-f]{40}87$`)

	// P2SHSubstringRegexp substring of OP_HASH160 [hash] OP_EQUAL
	P2SHSubstringRegexp, _ = regexp.Compile(`a914[\da-f]{40}87`)

	// MetanetRegexp OP_FALSE OP_RETURN 1635018093
	MetanetRegexp, _ = regexp.Compile("^006a046d65746142")

	// MetanetSubstringRegexp substring of OP_FALSE OP_RETURN 1635018093
	MetanetSubstringRegexp, _ = regexp.Compile("006a046d65746142")

	// OpReturnRegexp OP_FALSE OP_RETURN
	OpReturnRegexp, _ = regexp.Compile("^006a")

	// OpReturnSubstringRegexp substring of OP_FALSE OP_RETURN
	OpReturnSubstringRegexp, _ = regexp.Compile("006a")

	// StasRegexp OP_DUP OP_HASH160 [pubkey hash] OP_EQUALVERIFY OP_CHECKSIG OP_VERIFY ...
	StasRegexp, _ = regexp.Compile("^76a914[0-9a-f]{40}88ac6976aa607f5f7f7c5e7f7c5d7f7c5c7f7c5b7f7c5a7f7c597f7c587f7c577f7c567f7c557f7c547f7c537f7c527f7c517f7c7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7c5f7f7c5e7f7c5d7f7c5c7f7c5b7f7c5a7f7c597f7c587f7c577f7c567f7c557f7c547f7c537f7c527f7c517f7c7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e01007e818b21414136d08c5ed2bf3ba048afe6dcaebafeffffffffffffffffffffffffffffff007d976e7c5296a06394677768827601249301307c7e23022079be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798027e7c7e7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e01417e21038ff83d8cf12121491609c4939dc11c4aa35503508fe432dc5a5c1905608b9218ad547f7701207f01207f7701247f517f7801007e8102fd00a063546752687f7801007e817f727e7b01177f777b557a766471567a577a786354807e7e676d68aa880067765158a569765187645294567a5379587a7e7e78637c8c7c53797e577a7e6878637c8c7c53797e577a7e6878637c8c7c53797e577a7e6878637c8c7c53797e577a7e6878637c8c7c53797e577a7e6867567a6876aa587a7d54807e577a597a5a7a786354807e6f7e7eaa727c7e676d6e7eaa7c687b7eaa587a7d877663516752687c72879b69537a647500687c7b547f77517f7853a0916901247f77517f7c01007e817602fc00a06302fd00a063546752687f7c01007e816854937f77788c6301247f77517f7c01007e817602fc00a06302fd00a063546752687f7c01007e816854937f777852946301247f77517f7c01007e817602fc00a06302fd00a063546752687f7c01007e816854937f77686877517f7c52797d8b9f7c53a09b91697c76638c7c587f77517f7c01007e817602fc00a06302fd00a063546752687f7c01007e81687f777c6876638c7c587f77517f7c01007e817602fc00a06302fd00a063546752687f7c01007e81687f777c6863587f77517f7c01007e817602fc00a06302fd00a063546752687f7c01007e81687f7768587f517f7801007e817602fc00a06302fd00a063546752687f7801007e81727e7b7b687f75537f7c0376a9148801147f775379645579887567726881766968789263556753687a76026c057f7701147f8263517f7c766301007e817f7c6775006877686b537992635379528763547a6b547a6b677c6b567a6b537a7c717c71716868547a587f7c81547a557964936755795187637c686b687c547f7701207f75748c7a7669765880748c7a76567a876457790376a9147e7c7e557967041976a9147c7e0288ac687e7e5579636c766976748c7a9d58807e6c0376a9147e748c7a7e6c7e7e676c766b8263828c007c80517e846864745aa0637c748c7a76697d937b7b58807e56790376a9147e748c7a7e55797e7e6868686c567a5187637500678263828c007c80517e846868647459a0637c748c7a76697d937b7b58807e55790376a9147e748c7a7e55797e7e687459a0637c748c7a76697d937b7b58807e55790376a9147e748c7a7e55797e7e68687c537a9d547963557958807e041976a91455797e0288ac7e7e68aa87726d77776a14[0-9a-fA-F]{40}(0100|0101).*$")

	// StasSubstringRegexp substring of OP_DUP OP_HASH160 [pubkey hash] OP_EQUALVERIFY OP_CHECKSIG OP_VERIFY ...
	StasSubstringRegexp, _ = regexp.Compile("76a914[0-9a-f]{40}88ac6976aa607f5f7f7c5e7f7c5d7f7c5c7f7c5b7f7c5a7f7c597f7c587f7c577f7c567f7c557f7c547f7c537f7c527f7c517f7c7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7c5f7f7c5e7f7c5d7f7c5c7f7c5b7f7c5a7f7c597f7c587f7c577f7c567f7c557f7c547f7c537f7c527f7c517f7c7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e01007e818b21414136d08c5ed2bf3ba048afe6dcaebafeffffffffffffffffffffffffffffff007d976e7c5296a06394677768827601249301307c7e23022079be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798027e7c7e7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e01417e21038ff83d8cf12121491609c4939dc11c4aa35503508fe432dc5a5c1905608b9218ad547f7701207f01207f7701247f517f7801007e8102fd00a063546752687f7801007e817f727e7b01177f777b557a766471567a577a786354807e7e676d68aa880067765158a569765187645294567a5379587a7e7e78637c8c7c53797e577a7e6878637c8c7c53797e577a7e6878637c8c7c53797e577a7e6878637c8c7c53797e577a7e6878637c8c7c53797e577a7e6867567a6876aa587a7d54807e577a597a5a7a786354807e6f7e7eaa727c7e676d6e7eaa7c687b7eaa587a7d877663516752687c72879b69537a647500687c7b547f77517f7853a0916901247f77517f7c01007e817602fc00a06302fd00a063546752687f7c01007e816854937f77788c6301247f77517f7c01007e817602fc00a06302fd00a063546752687f7c01007e816854937f777852946301247f77517f7c01007e817602fc00a06302fd00a063546752687f7c01007e816854937f77686877517f7c52797d8b9f7c53a09b91697c76638c7c587f77517f7c01007e817602fc00a06302fd00a063546752687f7c01007e81687f777c6876638c7c587f77517f7c01007e817602fc00a06302fd00a063546752687f7c01007e81687f777c6863587f77517f7c01007e817602fc00a06302fd00a063546752687f7c01007e81687f7768587f517f7801007e817602fc00a06302fd00a063546752687f7801007e81727e7b7b687f75537f7c0376a9148801147f775379645579887567726881766968789263556753687a76026c057f7701147f8263517f7c766301007e817f7c6775006877686b537992635379528763547a6b547a6b677c6b567a6b537a7c717c71716868547a587f7c81547a557964936755795187637c686b687c547f7701207f75748c7a7669765880748c7a76567a876457790376a9147e7c7e557967041976a9147c7e0288ac687e7e5579636c766976748c7a9d58807e6c0376a9147e748c7a7e6c7e7e676c766b8263828c007c80517e846864745aa0637c748c7a76697d937b7b58807e56790376a9147e748c7a7e55797e7e6868686c567a5187637500678263828c007c80517e846868647459a0637c748c7a76697d937b7b58807e55790376a9147e748c7a7e55797e7e687459a0637c748c7a76697d937b7b58807e55790376a9147e748c7a7e55797e7e68687c537a9d547963557958807e041976a91455797e0288ac7e7e68aa87726d77776a14[0-9a-fA-F]{40}(0100|0101).*")

	// SensibleRegexp OP_FALSE OP_RETURN 115 101 110 115 105 98 108 101
	// TODO this detects sensible tokens, but is very broad. Make more refined
	SensibleRegexp, _ = regexp.Compile("73656e7369626c65")

	// SensibleSubstringRegexp substring of OP_FALSE OP_RETURN 115 101 110 115 105 98 108 101
	SensibleSubstringRegexp, _ = regexp.Compile("73656e7369626c65")
)

// IsP2PK Check whether the given string is a p2pk output
// This is the original destination type that was used in the first blocks by Satoshi
func IsP2PK(lockingScript string) bool {
	script, err := bscript2.NewFromHexString(lockingScript)
	if err != nil {
		return false
	}
	return script.IsP2PK()
}

// IsP2PKH Check whether the given string is a p2pkh output
func IsP2PKH(lockingScript string) bool {
	return P2PKHRegexp.MatchString(lockingScript)
}

// IsP2SH Check whether the given string is a p2shHex output
func IsP2SH(lockingScript string) bool {
	return P2SHRegexp.MatchString(lockingScript)
}

// IsMetanet Check whether the given string is a metanet output
func IsMetanet(lockingScript string) bool {
	return MetanetRegexp.MatchString(lockingScript)
}

// IsOpReturn Check whether the given string is an op_return
func IsOpReturn(lockingScript string) bool {
	return OpReturnRegexp.MatchString(lockingScript)
}

// IsStas Check whether the given string is a STAS token output
func IsStas(lockingScript string) bool {
	return StasRegexp.MatchString(lockingScript)
}

// IsSensible Check whether the given string is a Sensible token output
func IsSensible(lockingScript string) bool {
	return SensibleRegexp.MatchString(lockingScript)
}

// IsRunJS Check whether the given string is a Run JS output
func IsRunJS(lockingScript string) bool {
	return StasRegexp.MatchString(lockingScript)
}

// IsMultiSig Check whether the given string is a multi-sig locking script
func IsMultiSig(lockingScript string) bool {
	script, err := bscript2.NewFromHexString(lockingScript)
	if err != nil {
		return false
	}
	return script.IsMultiSigOut()
}

// GetDestinationType Get the type of output script destination
func GetDestinationType(lockingScript string) string {
	if IsP2PKH(lockingScript) {
		return ScriptTypePubKeyHash
	} else if IsMetanet(lockingScript) {
		// metanet is a special op_return - needs to be checked first
		return ScriptMetanet
	} else if IsOpReturn(lockingScript) {
		return ScriptTypeNullData
	} else if IsP2SH(lockingScript) {
		return ScriptHashType
	} else if IsMultiSig(lockingScript) {
		return ScriptTypeMultiSig
	} else if IsStas(lockingScript) {
		return ScriptTypeTokenStas
	} else if IsSensible(lockingScript) {
		return ScriptTypeTokenSensible
	} else if IsP2PK(lockingScript) {
		return ScriptTypePubKey
	}

	return ScriptTypeNonStandard
}

// GetDestinationTypeRegex Get the regex for destination type
func GetDestinationTypeRegex(destType string) *regexp.Regexp {
	switch destType {
	case ScriptTypePubKeyHash:
		return P2PKHSubstringRegexp
	case ScriptMetanet:
		return MetanetSubstringRegexp
	case ScriptTypeTokenStas:
		return StasSubstringRegexp
	case ScriptTypeTokenSensible:
		return SensibleSubstringRegexp
	}
	return nil
}

// GetAddressFromScript gets the destination address from the given locking script
func GetAddressFromScript(lockingScript string) (address string) {
	scriptType := GetDestinationType(lockingScript)
	if scriptType == ScriptTypePubKeyHash {
		address, _ = bitcoin.GetAddressFromScript(lockingScript)
	} else if scriptType == ScriptTypePubKey {
		s, err := bscript2.NewFromHexString(lockingScript) // no need for error check, if type is set, it should be valid
		if err != nil {
			return ""
		}

		var parts [][]byte
		parts, err = bscript2.DecodeParts(*s)
		if err != nil {
			return ""
		}

		pubkeyHex := hex.EncodeToString(parts[0])

		var addressScript *bscript2.Address
		addressScript, err = bitcoin.GetAddressFromPubKeyString(pubkeyHex, len(parts[0]) <= 33)
		if err == nil {
			address = addressScript.AddressString
		}
	} else if scriptType == ScriptTypeTokenStas {
		// stas is just a normal PubKeyHash with more data appended
		address, _ = bitcoin.GetAddressFromScript(lockingScript[:50])
		// } else if scriptType == ScriptTypeTokenSensible {
		// sensible tokens do not have the receiving address in the token output, but in another output
		// sensible does not seem to be a utxo protocol, but an output protocol (all outputs of the tx matter)
		// TODO saving the utxo for the token in spv wallet engine
	}
	return address
}

// GetDestinationLockingScript gets the destination locking script from non-standard locking scripts that reference
// a destination (hashed pubkey) - this can be used to determine destination for tokens (like STAS, run etc.)
func GetDestinationLockingScript(lockingScript string) string {
	destinationType := GetDestinationType(lockingScript)

	switch destinationType {
	case ScriptTypeTokenStas:
		// STAS token, extract the destination from the token and return it's locking script
		// this will match the token to an existing destination and add the output to the correct xpub utxo
		tokenLockingScript, err := GetLockingScriptFromSTASLockingScript(lockingScript)
		if err != nil {
			return lockingScript
		}
		return tokenLockingScript
	}

	// just return the locking script, type is ScriptTypePubKeyHash or ScriptTypeNonStandard
	return lockingScript
}
