package providers

import (
	"context"
	"errors"
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
		Host:   "lisatech.net",
		Path:   "/oauth/authorize",
	}
	// Default Redeem URL for Wordpress.
	wordpressDefaultRedeemURL = &url.URL{
		Scheme: "https",
		Host:   "lisatech.net",
		Path:   "/oauth/token",
	}

	// Default Profile URL for Wordpress.
	wordpressDefaultProfileURL = &url.URL{
		Scheme: "https",
		Host:   "lisatech.net",
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
		validateURL: wordpressDefaultProfileURL,
	})
	return &WordpressProvider{ProviderData: p}
}

// EnrichSession updates the User & Email after the initial Redeem
func (p *WordpressProvider) EnrichSession(ctx context.Context, s *sessions.SessionState) error {
	type result struct {
		UserEmail string `json:"user_email,omitempty"`
		UserLogin string `json:"user_login,omitempty"`
	}
	var r result
	err := requests.New(p.ProfileURL.String()).
		WithContext(ctx).
		WithHeaders(makeOIDCHeader(s.AccessToken)).
		Do().
		UnmarshalInto(&r)
	if err != nil {
		return err
	}
	if r.UserEmail == "" {
		return errors.New("no email")
	}
	s.Email = r.UserEmail
	if r.UserLogin == "" {
		return errors.New("no login name")
	}
	s.User = r.UserLogin
	return nil
}

// ValidateSession validates the AccessToken
func (p *WordpressProvider) ValidateSession(ctx context.Context, s *sessions.SessionState) bool {
	return validateToken(ctx, p, s.AccessToken, makeOIDCHeader(s.AccessToken))
}
