package optionaltls

// DynamicOption controls optional TLS in runtime.
// It can be used to turn on/off the feature based on some condition
// For example, an atomic boolean flag
type DynamicOption interface {
	IsActive() bool
}

type DynamicOptionFunc func() bool

func (f DynamicOptionFunc) IsActive() bool {
	return f()
}
