package metrics

// SetXPubCount adds a value to the stats gauge with the label "xpub"
func (m *Metrics) SetXPubCount(value int64) {
	m.stats.WithLabelValues("xpub").Set(float64(value))
}

// SetUtxoCount adds a value to the stats gauge with the label "utxo"
func (m *Metrics) SetUtxoCount(value int64) {
	m.stats.WithLabelValues("utxo").Set(float64(value))
}

// SetPaymailCount adds a value to the stats gauge with the label "paymail"
func (m *Metrics) SetPaymailCount(value int64) {
	m.stats.WithLabelValues("paymail").Set(float64(value))
}

// SetDestinationCount adds a value to the stats gauge with the label "destination"
func (m *Metrics) SetDestinationCount(value int64) {
	m.stats.WithLabelValues("destination").Set(float64(value))
}

// SetAccessKeyCount adds a value to the stats gauge with the label "access_key
func (m *Metrics) SetAccessKeyCount(value int64) {
	m.stats.WithLabelValues("access_key").Set(float64(value))
}
