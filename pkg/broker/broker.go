package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"code.cloudfoundry.org/lager"
	"github.com/pivotal-cf/brokerapi"
)

type BrokerImpl struct {
	Logger lager.Logger
	Config Config
}

type Config struct {
	ServiceName    string
	ServicePlan    string
	BaseGUID       string
	Credentials    interface{}
	Tags           string
	ImageURL       string
	SysLogDrainURL string
	FakeAsync      bool
	Free           bool
}

func NewBrokerImpl(logger lager.Logger) (bkr *BrokerImpl) {
	var credentials interface{}
	json.Unmarshal([]byte(os.Getenv("CREDENTIALS")), &credentials)
	fmt.Printf("Credentials: %v\n", credentials)

	return &BrokerImpl{
		Logger: logger,
		Config: Config{
			BaseGUID:    os.Getenv("BASE_GUID"),
			ServiceName: os.Getenv("SERVICE_NAME"),
			ServicePlan: os.Getenv("SERVICE_PLAN_NAME"),
			Credentials: credentials,
			Tags:        os.Getenv("TAGS"),
			ImageURL:    os.Getenv("IMAGE_URL"),
			Free:        true,

			FakeAsync: os.Getenv("FAKE_ASYNC") == "true",
		},
	}
}

func (bkr *BrokerImpl) Services(ctx context.Context) ([]brokerapi.Service, error) {
	return []brokerapi.Service{
		brokerapi.Service{
			ID:          bkr.Config.BaseGUID + "-service-" + bkr.Config.ServiceName,
			Name:        bkr.Config.ServiceName,
			Description: "Shared service for " + bkr.Config.ServiceName,
			Bindable:    true,
			Metadata: &brokerapi.ServiceMetadata{
				DisplayName: bkr.Config.ServiceName,
				ImageUrl:    bkr.Config.ImageURL,
			},
			Plans: []brokerapi.ServicePlan{
				brokerapi.ServicePlan{
					ID:          bkr.Config.BaseGUID + "-plan-" + bkr.Config.ServicePlan,
					Name:        bkr.Config.ServicePlan,
					Description: "Shared service for " + bkr.Config.ServiceName,
					Free:        &bkr.Config.Free,
				},
			},
		},
	}, nil
}

func (bkr *BrokerImpl) Provision(ctx context.Context, instanceID string, details brokerapi.ProvisionDetails, asyncAllowed bool) (brokerapi.ProvisionedServiceSpec, error) {
	return brokerapi.ProvisionedServiceSpec{
		IsAsync: bkr.Config.FakeAsync,
	}, nil
}

func (bkr *BrokerImpl) Deprovision(ctx context.Context, instanceID string, details brokerapi.DeprovisionDetails, asyncAllowed bool) (brokerapi.DeprovisionServiceSpec, error) {
	return brokerapi.DeprovisionServiceSpec{
		IsAsync: bkr.Config.FakeAsync,
	}, nil
}

func (bkr *BrokerImpl) GetInstance(ctx context.Context, instanceID string) (brokerapi.GetInstanceDetailsSpec, error) {
	panic("not implemented")
}

func (bkr *BrokerImpl) Bind(ctx context.Context, instanceID string, bindingID string, details brokerapi.BindDetails, asyncAllowed bool) (brokerapi.Binding, error) {
	return brokerapi.Binding{
		Credentials: bkr.Config.Credentials,
	}, nil
}

func (bkr *BrokerImpl) Unbind(ctx context.Context, instanceID string, bindingID string, details brokerapi.UnbindDetails, asyncAllowed bool) (brokerapi.UnbindSpec, error) {
	return brokerapi.UnbindSpec{}, nil
}

func (bkr *BrokerImpl) GetBinding(ctx context.Context, instanceID string, bindingID string) (brokerapi.GetBindingSpec, error) {
	panic("not implemented")
}

func (bkr *BrokerImpl) Update(ctx context.Context, instanceID string, details brokerapi.UpdateDetails, asyncAllowed bool) (brokerapi.UpdateServiceSpec, error) {
	return brokerapi.UpdateServiceSpec{
		IsAsync: bkr.Config.FakeAsync,
	}, nil
}

func (bkr *BrokerImpl) LastOperation(ctx context.Context, instanceID string, details brokerapi.PollDetails) (brokerapi.LastOperation, error) {
	return brokerapi.LastOperation{
		State: brokerapi.Succeeded,
	}, nil
}

func (bkr *BrokerImpl) LastBindingOperation(ctx context.Context, instanceID string, bindingID string, details brokerapi.PollDetails) (brokerapi.LastOperation, error) {
	panic("not implemented")
}
