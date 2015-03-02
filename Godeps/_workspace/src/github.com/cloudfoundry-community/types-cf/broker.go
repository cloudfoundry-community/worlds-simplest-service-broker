package cf

// ServiceCreationRequest describes Cloud Foundry service provisioning request
type ServiceCreationRequest struct {
	InstanceID       string `json:"-"`
	ServiceID        string `json:"service_id"`
	PlanID           string `json:"plan_id"`
	OrganizationGUID string `json:"organization_guid"`
	SpaceGUID        string `json:"space_guid"`
}

// ServiceCreationResponse describes Cloud Foundry service provisioning response
type ServiceCreationResponse struct {
	DashboardURL string `json:"dashboard_url"`
}

// ServiceBindingRequest describes Cloud Foundry service binding request
type ServiceBindingRequest struct {
	InstanceID string `json:"-"`
	BindingID  string `json:"-"`
	ServiceID  string `json:"service_id"`
	PlanID     string `json:"plan_id"`
	AppGUID    string `json:"app_guid"`
}

// ServiceBindingResponse describes Cloud Foundry service binding response
type ServiceBindingResponse struct {
	Credentials    map[string]string `json:"credentials"`
	SyslogDrainURL string            `json:"syslog_drain_url"`
}

// BrokerError describes Cloud Foundry broker error
type BrokerError struct {
	Description string `json:"description"`
}
