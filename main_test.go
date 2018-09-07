package main

import (
	"encoding/json"
	"github.com/cloudfoundry-community/types-cf"
	"github.com/go-martini/martini"
	"io/ioutil"
	"reflect"
	"testing"
)

type response struct {
	header int
	body   []byte
}

func TestGetCatalog(t *testing.T) {
	h, b := brokerCatalog()
	raw, _ := ioutil.ReadFile("test/catalog.json")
	e := response{
		200,
		raw,
	}

	testCases := []struct {
		header   int
		body     []byte
		name     string
		expected response
	}{
		{
			header:   h,
			body:     b,
			name:     "Should return empty catalog for empty env vars",
			expected: e,
		},
	}

	for _, tc := range testCases {
		body := &cf.Catalog{}
		expectedBody := &cf.Catalog{}
		assertBodyAgainstExpectations(tc.body, body, tc.expected, expectedBody, tc.header, tc.name, t, true)
	}
}

func TestGetLastOperation(t *testing.T) {
	p := make(martini.Params)
	p["service_id"] = "some service id"
	h, b := lastOperation(p)

	lastOp := lastOperationResponse{
		State:       "succeeded",
		Description: "async in action",
	}
	l, _ := json.Marshal(lastOp)
	expectResp := response{
		200,
		l,
	}

	testCases := []struct {
		header   int
		body     []byte
		name     string
		expected response
	}{
		{
			header:   h,
			body:     b,
			name:     "Should return last operation",
			expected: expectResp,
		},
	}

	for _, tc := range testCases {
		body := &lastOperationResponse{}
		expectedBody := &lastOperationResponse{}
		assertBodyAgainstExpectations(tc.body, body, tc.expected, expectedBody, tc.header, tc.name, t, true)
	}
}

func TestCreateServiceInstance(t *testing.T) {
	dashboardURL = "amazing.dashboard.com"
	p := make(martini.Params)
	p["service_id"] = "some service id"
	h, b := createServiceInstance(p)

	creationResp := cf.ServiceCreationResponse{
		DashboardURL: "amazing.dashboard.com",
	}
	cr, _ := json.Marshal(creationResp)
	expectResp := response{
		201,
		cr,
	}

	testCases := []struct {
		header   int
		body     []byte
		name     string
		expected response
	}{
		{
			header:   h,
			body:     b,
			name:     "Should return a successful creation response",
			expected: expectResp,
		},
	}

	for _, tc := range testCases {
		body := &cf.ServiceCreationResponse{}
		expectedBody := &cf.ServiceCreationResponse{}
		assertBodyAgainstExpectations(tc.body, body, tc.expected, expectedBody, tc.header, tc.name, t, true)
	}
}

func TestDeleteServiceInstance(t *testing.T) {
	p := make(martini.Params)
	p["service_id"] = "some service id"
	h, b := deleteServiceInstance(p)

	expectResp := response{
		200,
		[]byte("{}"),
	}

	testCases := []struct {
		header   int
		body     []byte
		name     string
		expected response
	}{
		{
			header:   h,
			body:     []byte(b),
			name:     "Should return a successful deletion response",
			expected: expectResp,
		},
	}

	for _, tc := range testCases {
		assertBodyAgainstExpectations(tc.body, tc.body, tc.expected, tc.expected.body, tc.header, tc.name, t, false)
	}
}

func TestCreateServiceBinding(t *testing.T) {
	credentials = `{"port": "5514", "host": "syslog-app.snpaas.eu"}`
	p := make(martini.Params)
	p["service_id"] = "some service id"
	h, b := createServiceBinding(p)

	credentials = `{"port": 5514, "host": "syslog-app.snpaas.eu"}`
	he, be:=createServiceBinding(p)

	cred := make(map[string]string)
	cred["port"] = "5514"
	cred["host"] = "syslog-app.snpaas.eu"

	creationResp := cf.ServiceBindingResponse{
		Credentials: cred,
	}
	cr, _ := json.Marshal(creationResp)
	expectResp := response{
		201,
		cr,
	}

	expectErrResp := response{
		500,
		[]byte{},
	}

	testCases := []struct {
		header   int
		body     []byte
		name     string
		expected response
	}{
		{
			header:   h,
			body:     b,
			name:     "Should return a successful binding response",
			expected: expectResp,
		},
		{
			header:   he,
			body:     be,
			name:     "Should return an error for credentials in the wrong format",
			expected: expectErrResp,
		},
	}

	for _, tc := range testCases {
		body := &cf.ServiceBindingResponse{}
		expectedBody := &cf.ServiceBindingResponse{}
		assertBodyAgainstExpectations(tc.body, body, tc.expected, expectedBody, tc.header, tc.name, t, true)
	}
}

func TestDeleteServiceBinding(t *testing.T) {
	p := make(martini.Params)
	p["service_id"] = "some service id"
	h, b := deleteServiceBinding(p)

	expectResp := response{
		200,
		[]byte("{}"),
	}

	testCases := []struct {
		header   int
		body     []byte
		name     string
		expected response
	}{
		{
			header:   h,
			body:     []byte(b),
			name:     "Should return a successful deletion response",
			expected: expectResp,
		},
	}

	for _, tc := range testCases {
		assertBodyAgainstExpectations(tc.body, tc.body, tc.expected, tc.expected.body, tc.header, tc.name, t, false)
	}
}

func TestShowServiceDashboard(t *testing.T) {
	p := make(martini.Params)
	h, b := showServiceInstanceDashboard(p)

	expectResp := response{
		200,
		[]byte("Dashboard"),
	}

	testCases := []struct {
		header   int
		body     []byte
		name     string
		expected response
	}{
		{
			header:   h,
			body:     []byte(b),
			name:     "Should return a successful deletion response",
			expected: expectResp,
		},
	}

	for _, tc := range testCases {
		assertBodyAgainstExpectations(tc.body, tc.body, tc.expected, tc.expected.body, tc.header, tc.name, t, false)
	}
}

func assertBodyAgainstExpectations(
	tcBody []byte,
	body interface{},
	tcExpected response,
	expectedBody interface{},
	header int,
	tcName string,
	t *testing.T,
	handleJson bool,
) {
	if handleJson {
		json.Unmarshal(tcBody, body)
		json.Unmarshal(tcExpected.body, expectedBody)
	}
	if header != tcExpected.header || !reflect.DeepEqual(body, expectedBody) {
		j, _ := json.Marshal(body)
		je, _ := json.Marshal(expectedBody)
		t.Errorf(
			"Test %s should return header |%d| and body |%s| but returned header |%d| and body |%s|",
			tcName,
			tcExpected.header,
			string(je),
			header,
			string(j),
		)
	}
}
