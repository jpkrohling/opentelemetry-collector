package dependency

import "context"

// WantsAuthenticator defines an interface that a component has to implement in order to get an authenticator set
type WantsAuthenticator interface {
	SetAuthenticator(Authenticator)
}

// Authenticator defines the contract for authenticators to fullfil in order to be used by other components
type Authenticator func(ctx context.Context, headers map[string][]string) (context.Context, error)
