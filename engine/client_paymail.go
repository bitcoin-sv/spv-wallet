package engine

import (
	"github.com/bitcoin-sv/go-paymail"
	paymailclient "github.com/bitcoin-sv/spv-wallet/engine/paymail"
)

// PaymailClient will return the Paymail if it exists
func (c *Client) PaymailClient() paymail.ClientInterface {
	if c.options.paymail != nil && c.options.paymail.client != nil {
		return c.options.paymail.Client()
	}
	return nil
}

// PaymailService will return the Paymail Service if it exists
func (c *Client) PaymailService() paymailclient.ServiceClient {
	if c.options.paymail != nil && c.options.paymail.service != nil {
		return c.options.paymail.ServiceClient()
	}
	return nil
}

// GetPaymailConfig will return the Paymail server config if it exists
func (c *Client) GetPaymailConfig() *PaymailServerOptions {
	if c.options.paymail != nil && c.options.paymail.serverConfig != nil {
		return c.options.paymail.serverConfig
	}
	return nil
}

// Client will return the paymail client from the options struct
func (p *paymailOptions) Client() paymail.ClientInterface {
	return p.client
}

// ServiceClient will return the paymail service client from the options struct
func (p *paymailOptions) ServiceClient() paymailclient.ServiceClient {
	return p.service
}

// FromSender will return either the configuration value or the application default
func (p *paymailOptions) FromSender() string {
	if len(p.serverConfig.DefaultFromPaymail) > 0 {
		return p.serverConfig.DefaultFromPaymail
	}
	return defaultSenderPaymail
}

// ServerConfig will return the Paymail Server configuration from the options struct
func (p *paymailOptions) ServerConfig() *PaymailServerOptions {
	return p.serverConfig
}
