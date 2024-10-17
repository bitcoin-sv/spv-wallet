package config

import (
	"net"
	"net/url"
	"regexp"
	"strings"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
)

// explicitHTTPURLRegex is a regex pattern to check the callback URL (host)
var explicitHTTPURLRegex = regexp.MustCompile(`^https?://`)

// Validate checks the configuration for specific rules
func (n *ARCConfig) Validate() error {
	if n == nil {
		return spverrors.Newf("arc is not configured")
	}

	if n.URL == "" {
		return spverrors.Newf("arc url is not configured")
	}

	if !n.isValidCallbackURL() {
		return spverrors.Newf("invalid callback host: %s - must be a valid external url - not a localhost", n.Callback.Host)
	}

	return nil
}

func (n *ARCConfig) isValidCallbackURL() bool {
	if !n.Callback.Enabled {
		return true
	}

	callbackUrl := n.Callback.Host

	if !explicitHTTPURLRegex.MatchString(callbackUrl) {
		return false
	}
	u, err := url.Parse(callbackUrl)
	if err != nil {
		return false
	}

	hostname := u.Hostname()

	return !n.isLocalNetworkHost(hostname)
}

func (n *ARCConfig) isLocalNetworkHost(hostname string) bool {
	if strings.Contains(hostname, "localhost") {
		return true
	}

	ip := net.ParseIP(hostname)
	if ip != nil {
		_, private10, _ := net.ParseCIDR("10.0.0.0/8")
		_, private172, _ := net.ParseCIDR("172.16.0.0/12")
		_, private192, _ := net.ParseCIDR("192.168.0.0/16")
		_, loopback, _ := net.ParseCIDR("127.0.0.0/8")
		_, linkLocal, _ := net.ParseCIDR("169.254.0.0/16")

		return private10.Contains(ip) || private172.Contains(ip) || private192.Contains(ip) || loopback.Contains(ip) || linkLocal.Contains(ip)
	}

	return false
}
