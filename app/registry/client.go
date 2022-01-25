package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
)

// RegisterService registers the service with the registry service - with the
// given `Registration` by making a RPC.
func RegisterService(r Registration) error {
	serviceUpdateURL, err := url.Parse(r.ServiceUpdateUrl)
	if err != nil {
		return err
	}
	http.Handle(serviceUpdateURL.Path, &serviceUpdateHandler{})

	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	err = enc.Encode(r)
	if err != nil {
		return err
	}
	res, err := http.Post(ServicesURL, "application/json", buf)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to register service. Registry service responded with code %v", res.StatusCode)
	}
	return nil
}

type serviceUpdateHandler struct{}

func (sh serviceUpdateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	dec := json.NewDecoder(r.Body)
	var p patch
	err := dec.Decode(&p)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	prov.Update(p)
}

// ShutdownService deregisters a service with the given url from the registry.
// Note that we're not taking the service name here and instead using the
// service url because there can be possibly more than one service running with
// the same service name (they're the same service type).
func ShutdownService(serviceURL string) error {
	req, err := http.NewRequest(
		http.MethodDelete,
		ServicesURL,
		bytes.NewBuffer([]byte(serviceURL)),
	)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "text/plain")

	res, err := http.DefaultClient.Do(req)
	if err != nil || res.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to deregister service. Registry service responded with code %v", res.StatusCode)
	}
	return nil
}

type providers struct {
	// services is a mapping between between `ServiceName` and their URLs
	services map[ServiceName][]string
	mu       *sync.RWMutex
}

// Update takes a `patch` and stores the newly added `patchEntry`-ies and
// removes the one that are marked for removal. There can be multiple services
// with the same name running on different URLs i.e multiple grading service or
// multiple logging service.
func (p *providers) Update(pat patch) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Add the URLs for the added services to the providers list.
	for _, patchEntry := range pat.Added {
		if _, ok := p.services[patchEntry.Name]; !ok {
			p.services[patchEntry.Name] = make([]string, 0)
		}
		p.services[patchEntry.Name] = append(p.services[patchEntry.Name], patchEntry.URL)
	}

	// Remove the URLs for the removed services from the providers list.
	for _, patchEntry := range pat.Removed {
		providerURLs, ok := p.services[patchEntry.Name]
		if !ok {
			continue
		}
		for i := range providerURLs {
			if providerURLs[i] == patchEntry.URL {
				p.services[patchEntry.Name] = append(
					providerURLs[:i],
					providerURLs[:i+1]...,
				)
			}
		}
	}
}

func (p providers) get(name ServiceName) (string, error) {
	providers, ok := p.services[name]
	if !ok {
		return "", fmt.Errorf("No providers available for service %v", name)
	}
	idx := int(rand.Float32() * float32(len(providers)))
	return providers[idx], nil
}

func GetProvider(name ServiceName) (string, error) {
	return prov.get(name)
}

var prov = providers{
	services: make(map[ServiceName][]string),
	mu:       new(sync.RWMutex),
}
