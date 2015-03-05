package goth

import "fmt"

// Provider needs to be implemented for each 3rd party authentication provider
// e.g. Facebook, Twitter, etc...
type BaseProvider interface {
	Name() string
	Debug(bool)
}

type Provider interface {
	BaseProvider
	UnmarshalSession(string) (Session, error)
	FetchUser(Session) (User, error)
	BeginAuth(state string) (Session, error)
}

type Verifier interface {
	BaseProvider
	VerifyAuth(access_token string) (User, error)
}

// Providers is list of known/available providers.
type Providers map[string]BaseProvider

var providers = Providers{}

// UseProviders sets a list of available providers for use with Goth.
func UseProviders(viders ...BaseProvider) {
	for _, provider := range viders {
		providers[provider.Name()] = provider
	}
}

// GetProviders returns a list of all the providers currently in use.
func GetProviders() Providers {
	return providers
}

// GetProvider returns a previously created provider. If Goth has not
// been told to use the named provider it will return an error.
func GetProvider(name string) (BaseProvider, error) {
	provider := providers[name]
	if provider == nil {
		return nil, fmt.Errorf("no provider for %s exists", name)
	}
	return provider, nil
}

// ClearProviders will remove all providers currently in use.
// This is useful, mostly, for testing purposes.
func ClearProviders() {
	providers = Providers{}
}
