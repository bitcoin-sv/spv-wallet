package filter

// SearchDestinations is a model for handling searching with filters and metadata
type SearchDestinations = SearchModel[DestinationFilter]

// CountDestinations is a model for handling counting filtered destinations
type CountDestinations = ConditionsModel[DestinationFilter]

// SearchContacts is a model for handling searching with filters and metadata
type SearchContacts = SearchModel[ContactFilter]

// AdminSearchPaymails is a model for handling searching with filters and metadata
type AdminSearchPaymails = SearchModel[AdminPaymailFilter]

// AdminCountPaymails is a model for handling counting filtered paymails
type AdminCountPaymails = ConditionsModel[AdminPaymailFilter]

// SearchTransactions is a model for handling searching with filters and metadata
type SearchTransactions = SearchModel[TransactionFilter]

// CountTransactions is a model for handling counting filtered transactions
type CountTransactions = ConditionsModel[TransactionFilter]

// SearchXpubs is a model for handling searching with filters and metadata
type SearchXpubs = SearchModel[XpubFilter]

// CountXpubs is a model for handling counting filtered xPubs
type CountXpubs = ConditionsModel[XpubFilter]

// AdminSearchUtxos is a model for handling searching with filters and metadata
type AdminSearchUtxos = SearchModel[AdminUtxoFilter]

// AdminCountUtxos is a model for handling counting filtered UTXOs
type AdminCountUtxos = ConditionsModel[AdminUtxoFilter]

// SearchUtxos is a model for handling searching with filters and metadata
type SearchUtxos = SearchModel[UtxoFilter]

// CountUtxos is a model for handling counting filtered UTXOs
type CountUtxos = ConditionsModel[UtxoFilter]

// AdminSearchAccessKeys is a model for handling searching with filters and metadata
type AdminSearchAccessKeys = SearchModel[AdminAccessKeyFilter]

// AdminCountAccessKeys is a model for handling counting filtered transactions
type AdminCountAccessKeys = ConditionsModel[AdminAccessKeyFilter]

// SearchAccessKeys is a model for handling searching with filters and metadata
type SearchAccessKeys = SearchModel[AccessKeyFilter]

// CountAccessKeys is a model for handling counting filtered transactions
type CountAccessKeys = ConditionsModel[AccessKeyFilter]
