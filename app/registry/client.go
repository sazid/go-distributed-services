package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// RegisterService registers the service with the registry service with the
// given `Registration` by making a RPC.
func RegisterService(r Registration) error {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	err := enc.Encode(r)
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
