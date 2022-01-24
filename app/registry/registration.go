package registry

// Registration represents a service registration - to be used for registering
// and dealing with services.
type Registration struct {
	ServiceName ServiceName
	ServiceURL  string
}

// ServiceName as a separate type to avoid accidental use of strings.
type ServiceName string

// Name of the services.
const (
	LogService = ServiceName("LogService")
)
