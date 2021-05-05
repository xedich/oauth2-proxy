package providers

import (
	"context"
	"fmt"
	"net/url"

	"github.com/oauth2-proxy/oauth2-proxy/v7/pkg/apis/sessions"
	"github.com/oauth2-proxy/oauth2-proxy/v7/pkg/requests"
)

// WordpressProvider represents an Wordpress based Identity Provider
type WordpressProvider struct {
	*ProviderData
}

var (
	// Default Login URL for Wordpress.
	wordpressDefaultLoginURL = &url.URL{
		Scheme: "https",
		Host:   "loole.net",
		Path:   "/oauth/authorize",
	}
	// Default Redeem URL for Wordpress.
	wordpressDefaultRedeemURL = &url.URL{
		Scheme: "https",
		Host:   "loole.net",
		Path:   "/oauth/token",
	}

	// Default Profile URL for Wordpress.
	wordpressDefaultProfileURL = &url.URL{
		Scheme: "https",
		Host:   "loole.net",
		Path:   "/oauth/me",
	}
)

var _ Provider = (*WordpressProvider)(nil)

const wordpressProviderName = "Wordpress"

// NewWordpressProvider initiates a new WordpressProvider
func NewWordpressProvider(p *ProviderData) *WordpressProvider {
	p.setProviderDefaults(providerDefaults{
		name:        wordpressProviderName,
		loginURL:    wordpressDefaultLoginURL,
		redeemURL:   wordpressDefaultRedeemURL,
		profileURL:  wordpressDefaultProfileURL,
	})
	return &WordpressProvider{ProviderData: p}
}

// GetEmailAddress returns the Account email address
func (p *WordpressProvider) GetEmailAddress(ctx context.Context, s *sessions.SessionState) (string, error) {
	json, err := requests.New(p.ProfileURL.String()).
		WithContext(ctx).
		WithHeaders(makeOIDCHeader(s.AccessToken)).
		Do().
		UnmarshalJSON()
	if err != nil {
		return "", fmt.Errorf("error making request: %v", err)
	}
	fmt.Printf("%s", json)
	email, err := json.Get("user_email").String()
	return email, err
}
