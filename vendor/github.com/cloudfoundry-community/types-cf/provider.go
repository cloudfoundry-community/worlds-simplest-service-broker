package cf

// ServiceProvider defines the required provider functionality
type ServiceProvider interface {

	// GetCatalog returns the catalog of services managed by this broker
	GetCatalog() (*Catalog, *ServiceProviderError)

	// CreateService creates a service instance for specific plan
	CreateService(r *ServiceCreationRequest) (*ServiceCreationResponse, *ServiceProviderError)

	// DeleteService deletes previously created service instance
	DeleteService(instanceID string) *ServiceProviderError

	// BindService binds to specified service instance and
	// Returns credentials necessary to establish connection to that service
	BindService(r *ServiceBindingRequest) (*ServiceBindingResponse, *ServiceProviderError)

	// UnbindService removes previously created binding
	UnbindService(instanceID, bindingID string) *ServiceProviderError

	LastOperation(instanceID string) (*ServiceLastOperationResponse, *ServiceProviderError)
}
