package dexcomp

import (
	"context"
	"flag"
	"fmt"
	"github.com/coreos/go-oidc"
	sctx "github.com/phathdt/service-context"
	"golang.org/x/oauth2"
	"strings"
)

type DexComponent interface {
	GetOauthConfig() (*oauth2.Config, error)
	GetProvider() (*oidc.Provider, error)
	GetIdTokenVerifier() (*oidc.IDTokenVerifier, error)
	GetClientEndpoint() string
	GetClientErrEndpoint() string
	GetRedirect() bool
}

type dexcomp struct {
	id                string
	clientId          string
	clientSecret      string
	issuer            string
	atenEndpoint      string
	clientEndpoint    string
	clientErrEndpoint string
	scopes            string
	redirect          bool
}

func NewDexcomp(id string) *dexcomp {
	return &dexcomp{id: id}
}

func (d *dexcomp) ID() string {
	return d.id
}

func (d *dexcomp) InitFlags() {
	flag.StringVar(&d.clientId, "dex_client_id", "client_id", "dex client id")
	flag.StringVar(&d.clientSecret, "dex_client_secret", "client_secret", "dex client secret")
	flag.StringVar(&d.issuer, "dex_issuer", "http://127.0.0.1:5556", "dex issuer")
	flag.StringVar(&d.atenEndpoint, "dex_aten_endpoint", "http://localhost:4000", "dex aten endpoint")
	flag.StringVar(&d.clientEndpoint, "dex_client_endpoint", "http://localhost:3000/oauth/callback", "dex client endpoint")
	flag.StringVar(&d.clientErrEndpoint, "dex_client_err_endpoint", "http://localhost:3000/oauth/error", "dex client error endpoint")
	flag.StringVar(&d.scopes, "dex_scopes", "profile,email,groups,federated:id", "dex scopes ")
	flag.BoolVar(&d.redirect, "dex_redirect", true, "dex redirect or return json")
}

func (d *dexcomp) Activate(context sctx.ServiceContext) error {
	return nil
}

func (d *dexcomp) Stop() error {
	return nil
}

func (d *dexcomp) GetOauthConfig() (*oauth2.Config, error) {
	provider, err := d.GetProvider()
	if err != nil {
		return nil, err
	}

	scopes := append(strings.Split(d.scopes, ","), oidc.ScopeOpenID)

	return &oauth2.Config{
		// client_id and client_secret of the client.
		ClientID:     d.clientId,
		ClientSecret: d.clientSecret,

		// The redirectURL.
		RedirectURL: fmt.Sprintf("%s/auth/callback", d.atenEndpoint),

		// Discovery returns the OAuth2 endpoints.
		Endpoint: provider.Endpoint(),

		// "openid" is a required scope for OpenID Connect flows.
		//
		// Other scopes, such as "groups" can be requested.
		Scopes: scopes,
	}, nil
}

func (d *dexcomp) GetProvider() (*oidc.Provider, error) {
	return oidc.NewProvider(context.Background(), d.issuer)
}

func (d *dexcomp) GetIdTokenVerifier() (*oidc.IDTokenVerifier, error) {
	provider, err := d.GetProvider()
	if err != nil {
		return nil, err
	}

	idTokenVerifier := provider.Verifier(&oidc.Config{ClientID: d.clientId})

	return idTokenVerifier, nil
}

func (d *dexcomp) GetClientEndpoint() string {
	return d.clientEndpoint
}

func (d *dexcomp) GetClientErrEndpoint() string {
	return d.clientErrEndpoint
}

func (d *dexcomp) GetRedirect() bool {
	return d.redirect
}
