package gophercloud

// ProviderClient stores details that are required to interact with any services within a specific provider's API.
//
// Generally, you acquire a ProviderClient by calling the `NewClient()` method in the appropriate provider's child package,
// providing whatever authentication credentials are required.
type ProviderClient struct {

	// IdentityBase is the front door to an openstack provider.
	// Generally this will be populated when you authenticate.
	// It should be the *root* resource of the identity service, not of a specific identity version.
	IdentityBase string

	// IdentityEndpoint is the originally requested identity endpoint.
	// This may be a specific version of the identity service, in which case that endpoint is used rather than querying the
	// version-negotiation endpoint.
	IdentityEndpoint string

	// TokenID is the most recently valid token issued.
	TokenID string

	// EndpointLocator describes how this provider discovers the endpoints for its constituent services.
	EndpointLocator EndpointLocator
}

// AuthenticatedHeaders returns a map of HTTP headers that are common for all authenticated service requests.
func (client *ProviderClient) AuthenticatedHeaders() map[string]string {
	return map[string]string{"X-Auth-Token": client.TokenID}
}
