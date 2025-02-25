package errorcases

import "github.com/bitcoin-sv/spv-wallet/api/manualtests"

var SomeHopefullyNotRegisteredUserTemplate = manualtests.User{
	Alias:     "badrequest",
	Xpriv:     "xprv9s21ZrQH143K3TzwskUB1NW5iHx4EH7cxquXfwpFR5HrWe6HQrYvECJsj3sg1DJhWhwjtw5WdXwje8pkyuuvzJUingwo4f5BkD5dNubfNUn",
	Xpub:      "xpub661MyMwAqRbcFx5Qyn1BNWSpGKnYdjqUL4q8ULDryQpqPSRRxPsAmzdMaLnoyUzLAQ5ukXgMZjYs5LfNfsPFwBoSwxChePB1DxKvyFz6F67",
	PublicKey: "035dcb59eb7b5c5982ba6fbbccffbb2460f8daa07c5f9a21f2c2cf0845dcc6dda9",
	ID:        "1QCibNnc8CK7bzTdupM91f4PKixhMqqJQw",
}

func UserDefinitionForMakingBadRequests(domain string) manualtests.User {
	user := SomeHopefullyNotRegisteredUserTemplate
	user.Domain = domain
	return user
}
