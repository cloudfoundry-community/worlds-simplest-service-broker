package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/cloudfoundry-community/go-cfenv"
	"github.com/cloudfoundry-community/types-cf"
	"github.com/go-martini/martini"
	"github.com/kr/pretty"
)

func init() {
	log.SetFlags(log.Ltime | log.Lshortfile)
}

type serviceInstanceResponse struct {
	DashboardURL string `json:"dashboard_url"`
}

type serviceBindingResponse struct {
	Credentials    map[string]interface{} `json:"credentials"`
	SyslogDrainURL string                 `json:"syslog_drain_url"`
}

var serviceName, servicePlan, baseGUID, tags, imageURL string
var serviceBinding serviceBindingResponse
var appURL string

func brokerCatalog() (int, []byte) {
	tagArray := []string{}
	if len(tags) > 0 {
		tagArray = strings.Split(tags, ",")
	}
	var requires []string
	if len(serviceBinding.SyslogDrainURL) > 0 {
		requires = []string{"syslog_drain"}
	}
	catalog := cf.Catalog{
		Services: []*cf.Service{
			{
				ID:          baseGUID + "-service-" + serviceName,
				Name:        serviceName,
				Description: "Shared service for " + serviceName,
				Bindable:    true,
				Tags:        tagArray,
				Requires:    requires,
				Metadata: &cf.ServiceMeta{
					DisplayName: serviceName,
					ImageURL:    imageURL,
				},
				Plans: []*cf.Plan{
					{
						ID:          baseGUID + "-plan-" + servicePlan,
						Name:        servicePlan,
						Description: "Shared service for " + serviceName,
						Free:        true,
					},
				},
			},
		},
	}
	json, err := json.Marshal(catalog)
	if err != nil {
		fmt.Println("Um, how did we fail to marshal this catalog:")
		fmt.Printf("%# v\n", pretty.Formatter(catalog))
		return 500, []byte{}
	}
	return 200, json
}

func createServiceInstance(params martini.Params) (int, []byte) {
	serviceID := params["service_id"]
	fmt.Printf("Creating service instance %s for service %s plan %s\n", serviceID, serviceName, servicePlan)

	instance := serviceInstanceResponse{DashboardURL: fmt.Sprintf("%s/dashboard", appURL)}
	json, err := json.Marshal(instance)
	if err != nil {
		fmt.Println("Um, how did we fail to marshal this service instance:")
		fmt.Printf("%# v\n", pretty.Formatter(instance))
		return 500, []byte{}
	}
	return 201, json
}

func deleteServiceInstance(params martini.Params) (int, string) {
	serviceID := params["service_id"]
	fmt.Printf("Deleting service instance %s for service %s plan %s\n", serviceID, serviceName, servicePlan)
	return 200, "{}"
}

func createServiceBinding(params martini.Params) (int, []byte) {
	serviceID := params["service_id"]
	serviceBindingID := params["binding_id"]
	fmt.Printf("Creating service binding %s for service %s plan %s instance %s\n",
		serviceBindingID, serviceName, servicePlan, serviceID)

	json, err := json.Marshal(serviceBinding)
	if err != nil {
		fmt.Println("Um, how did we fail to marshal this binding:")
		fmt.Printf("%# v\n", pretty.Formatter(serviceBinding))
		return 500, []byte{}
	}
	return 201, json
}

func deleteServiceBinding(params martini.Params) (int, string) {
	serviceID := params["service_id"]
	serviceBindingID := params["binding_id"]
	fmt.Printf("Delete service binding %s for service %s plan %s instance %s\n",
		serviceBindingID, serviceName, servicePlan, serviceID)
	return 200, "{}"
}

func showServiceInstanceDashboard(params martini.Params) (int, string) {
	fmt.Printf("Show dashboard for service %s plan %s\n", serviceName, servicePlan)
	return 200, "Dashboard"
}

func main() {
	m := martini.Classic()

	baseGUID = os.Getenv("BASE_GUID")
	if baseGUID == "" {
		baseGUID = "29140B3F-0E69-4C7E-8A35"
	}
	serviceName = os.Getenv("SERVICE_NAME")
	if serviceName == "" {
		serviceName = "some-service-name" // replace with cfenv.AppName
	}
	servicePlan = os.Getenv("SERVICE_PLAN")
	if servicePlan == "" {
		servicePlan = "shared"
	}

	serviceBinding.SyslogDrainURL = os.Getenv("SYSLOG_DRAIN_URL")

	credentials := os.Getenv("CREDENTIALS")
	if credentials == "" {
		credentials = "{\"port\": \"4000\"}"
	}
	tags = os.Getenv("TAGS")
	imageURL = os.Getenv("IMAGE_URL")

	json.Unmarshal([]byte(credentials), &serviceBinding.Credentials)
	fmt.Printf("%# v\n", pretty.Formatter(serviceBinding))

	appEnv, err := cfenv.Current()
	if err == nil {
		appURL = fmt.Sprintf("https://%s", appEnv.ApplicationURIs[0])
	} else {
		appURL = "http://localhost:5000"
	}
	fmt.Println("Running as", appURL)

	// Cloud Foundry Service API
	m.Get("/v2/catalog", brokerCatalog)
	m.Put("/v2/service_instances/:service_id", createServiceInstance)
	m.Delete("/v2/service_instances/:service_id", deleteServiceInstance)
	m.Put("/v2/service_instances/:service_id/service_bindings/:binding_id", createServiceBinding)
	m.Delete("/v2/service_instances/:service_id/service_bindings/:binding_id", deleteServiceBinding)

	// Service Instance Dashboard
	m.Get("/dashboard", showServiceInstanceDashboard)

	m.Run()
}
