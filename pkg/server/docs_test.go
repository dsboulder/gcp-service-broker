// Copyright 2018 the Service Broker Project Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package server

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/providers/builtin"
	"github.com/gorilla/mux"
)

func TestAddDocsHandler(t *testing.T) {
	registry := builtin.BuiltinBrokerRegistry(broker.ServiceConfigMap{})
	router := mux.NewRouter()
	// Test that the handler sets the correct header and contains some imporant
	// strings that will indicate (but not prove!) that the rendering was correct.
	if err := AddDocsHandler(router, registry); err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest(http.MethodGet, "/docs", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, request)

	if w.Code != http.StatusOK {
		t.Errorf("Expected response code: %d got: %d", http.StatusOK, w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "text/html" {
		t.Errorf("Expected text/html content type got: %q", contentType)
	}

	importantStrings := []string{"<html", "bootstrap.min.css"}
	for _, svc := range registry.GetAllServices() {
		importantStrings = append(importantStrings, svc.Name)
	}

	body := w.Body.Bytes()

	for _, is := range importantStrings {
		if !bytes.Contains(body, []byte(is)) {
			t.Errorf("Expected body to contain the string %q", is)
		}
	}
}

func TestAddServiceConfigHandler(t *testing.T) {
	registry := builtin.BuiltinBrokerRegistry(broker.ServiceConfigMap{})
	router := mux.NewRouter()

	// Test that the handler sets the correct header and contains some imporant
	// strings that will indicate (but not prove!) that the rendering was correct.

	if err := AddServiceConfigHandler(router, registry); err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest(http.MethodGet, "/service-config", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, request)

	if w.Code != http.StatusOK {
		t.Errorf("Expected response code: %d got: %d", http.StatusOK, w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "text/html" {
		t.Errorf("Expected text/html content type got: %q", contentType)
	}

	importantStrings := []string{"<html", "bootstrap.min.css"}
	for _, svc := range registry.GetAllServices() {
		importantStrings = append(importantStrings, svc.Name)
		importantStrings = append(importantStrings, svc.Id)
	}

	body := w.Body.Bytes()
	for _, is := range importantStrings {
		if !bytes.Contains(body, []byte(is)) {
			t.Errorf("Expected body to contain the string %q", is)
		}
	}
}
