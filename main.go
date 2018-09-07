package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/cloudfoundry-community/go-cfenv"
	"github.com/cloudfoundry-community/types-cf"
	"github.com/go-martini/martini"
	"github.com/kr/pretty"
	"github.com/martini-contrib/auth"
)

var appURL, dashboardURL, syslogDrainUrl, credentials string
var serviceName, servicePlan, baseGUID, authUser, authPassword, tags, serviceDescription string
var metadataDisplayName, metadataLongDescription, metadataImageURL, metadataProviderDisplayName, metadataDocumentationUrl, metadataSupportUrl string
var fakeAsync bool

type lastOperationResponse struct {
	State       string `json:"state"`
	Description string `json:"description,omitempty"`
}

func init() {
	log.SetFlags(log.Ltime | log.Lshortfile)
}

func brokerCatalog() (int, []byte) {
	tagArray := []string{}
	if len(tags) > 0 {
		tagArray = strings.Split(tags, ",")
	}
	var requires []string
	if syslogDrainUrl != "" {
		requires = []string{"syslog_drain"}
	}
	catalog := cf.Catalog{
		Services: []*cf.Service{
			{
				ID:          baseGUID + "-service-" + serviceName,
				Name:        serviceName,
				Description: serviceDescription,
				Bindable:    true,
				Tags:        tagArray,
				Requires:    requires,
				Metadata: &cf.ServiceMeta{
					DisplayName:         metadataDisplayName,
					ImageURL:            metadataImageURL,
					Description:         metadataLongDescription,
					ProviderDisplayName: metadataProviderDisplayName,
					DocURL:              metadataDocumentationUrl,
					SupportURL:          metadataSupportUrl,
				},
				Plans: []*cf.Plan{
					{
						ID:          baseGUID + "-plan-" + servicePlan,
						Name:        servicePlan,
						Description: serviceDescription,
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
		return http.StatusInternalServerError, []byte{}
	}
	return http.StatusOK, json
}

func createServiceInstance(params martini.Params) (int, []byte) {
	serviceID := params["service_id"]
	fmt.Printf("Creating service instance %s for service %s plan %s\n", serviceID, serviceName, servicePlan)
	instance := cf.ServiceCreationResponse{
		DashboardURL: dashboardURL,
	}
	json, err := json.Marshal(instance)
	if err != nil {
		fmt.Println("Um, how did we fail to marshal this service instance:")
		fmt.Printf("%# v\n", pretty.Formatter(instance))
		return http.StatusInternalServerError, []byte{}
	}
	if fakeAsync {
		return http.StatusAccepted, json
	}
	return http.StatusCreated, json
}

func deleteServiceInstance(params martini.Params) (int, string) {
	serviceID := params["service_id"]
	fmt.Printf("Deleting service instance %s for service %s plan %s\n", serviceID, serviceName, servicePlan)
	if fakeAsync {
		return http.StatusAccepted, "{}"
	}
	return http.StatusOK, "{}"
}

func lastOperation(params martini.Params) (int, []byte) {
	lastOp := lastOperationResponse{
		State:       "succeeded",
		Description: "async in action",
	}
	json, err := json.Marshal(lastOp)
	if err != nil {
		fmt.Println("Um, how did we fail to marshal this service instance:")
		fmt.Printf("%# v\n", pretty.Formatter(lastOp))
		return http.StatusInternalServerError, []byte{}
	}
	return http.StatusOK, json
}

func createServiceBinding(params martini.Params) (int, []byte) {
	type serviceCredentials map[string]string
	serviceID := params["service_id"]
	serviceBindingID := params["binding_id"]
	fmt.Printf("Creating service binding %s for service %s plan %s instance %s\n", serviceBindingID, serviceName, servicePlan, serviceID)

	c := make(serviceCredentials)
	e := json.Unmarshal([]byte(credentials), &c)
	serviceBinding := cf.ServiceBindingResponse{
		Credentials: c,
	}
	if e != nil {
		fmt.Printf("Failed to load credentials: %s", credentials)
		return http.StatusInternalServerError, []byte{}
	}

	if syslogDrainUrl != "" {
		serviceBinding.SyslogDrainURL = syslogDrainUrl
	}
	json, err := json.Marshal(serviceBinding)
	if err != nil {
		fmt.Println("Um, how did we fail to marshal this binding:")
		fmt.Printf("%# v\n", pretty.Formatter(serviceBinding))
		return http.StatusInternalServerError, []byte{}
	}
	return http.StatusCreated, json
}

func deleteServiceBinding(params martini.Params) (int, string) {
	serviceID := params["service_id"]
	serviceBindingID := params["binding_id"]
	fmt.Printf("Delete service binding %s for service %s plan %s instance %s\n", serviceBindingID, serviceName, servicePlan, serviceID)
	return http.StatusOK, "{}"
}

func showServiceInstanceDashboard(params martini.Params) (int, string) {
	fmt.Printf("Show dashboard for service %s plan %s\n", serviceName, servicePlan)
	return http.StatusOK, "Dashboard"
}

func getEnvVar(v string, def string) string {
	r := os.Getenv(v)
	if r == "" {
		return def
	}
	return r
}

func main() {
	m := martini.Classic()

	appPort := getEnvVar("PORT", "3000") // default of martini
	appName := "some-service"
	appEnv, err := cfenv.Current()
	if err == nil {
		appURL = fmt.Sprintf("https://%s", appEnv.ApplicationURIs[0])
		appName = appEnv.Name
	} else {
		appURL = "http://localhost:" + appPort
	}
	baseGUID = getEnvVar("SERVICE_BASE_GUID", "29140B3F-0E69-4C7E-8A35")
	serviceName = getEnvVar("SERVICE_NAME", appName)
	servicePlan = getEnvVar("SERVICE_PLAN", "shared")
	serviceDescription = getEnvVar("SERVICE_DESCRIPTION", "Shared service for "+serviceName)
	authUser = getEnvVar("SERVICE_AUTH_USER", "")
	authPassword = getEnvVar("SERVICE_AUTH_PASSWORD", "")
	if (authUser != "") && (authPassword != "") {
		// secure service broker with basic auth if both env variables are set
		m.Use(auth.Basic(authUser, authPassword))
	}
	syslogDrainUrl = getEnvVar("SYSLOG_DRAIN_URL", "")
	tags = getEnvVar("SERVICE_TAGS", "")
	dashboardURL = getEnvVar("SERVICE_DASHBOARD_URL", fmt.Sprintf("%s/dashboard", appURL))
	metadataDisplayName = getEnvVar("SERVICE_METADATA_DISPLAYNAME", serviceName)
	metadataLongDescription = getEnvVar("SERVICE_METADATA_LONGDESC", serviceDescription)
	metadataImageURL = getEnvVar("SERVICE_METADATA_IMAGEURL", "")
	metadataProviderDisplayName = getEnvVar("SERVICE_METADATA_PROVIDERDISPLAYNAME", "")
	metadataDocumentationUrl = getEnvVar("SERVICE_METADATA_DOCURL", "")
	metadataSupportUrl = getEnvVar("SERVICE_METADATA_SUPPORTURL", "")
	credentials = getEnvVar("SERVICE_CREDENTIALS", "{}")
	// Each provision/deprovision request will support an async GET /last_operation request
	fakeAsync = getEnvVar("SERVICE_FAKE_ASYNC", "") == "true"

	// Cloud Foundry Service API
	m.Get("/v2/catalog", brokerCatalog)
	m.Get("/v2/service_instances/:service_id/last_operation", lastOperation)
	m.Put("/v2/service_instances/:service_id", createServiceInstance)
	m.Delete("/v2/service_instances/:service_id", deleteServiceInstance)
	m.Put("/v2/service_instances/:service_id/service_bindings/:binding_id", createServiceBinding)
	m.Delete("/v2/service_instances/:service_id/service_bindings/:binding_id", deleteServiceBinding)
	// Service Instance Dashboard
	m.Get("/dashboard", showServiceInstanceDashboard)

	fmt.Println("Running as", appURL)

	m.Run()
}
