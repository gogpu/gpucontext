package gpucontext

import "sync"

// Registry provides thread-safe registration and lookup of named factories.
// It supports priority-based selection when multiple implementations exist.
//
// Type parameter T is the type returned by factories.
//
// Example:
//
//	var backends = gpucontext.NewRegistry[Backend](
//	    gpucontext.WithPriority("vulkan", "dx12", "metal", "gles", "software"),
//	)
//
//	backends.Register("vulkan", func() Backend { return NewVulkanBackend() })
//	backends.Register("software", func() Backend { return NewSoftwareBackend() })
//
//	best := backends.Best() // Returns vulkan if available, otherwise software
type Registry[T any] struct {
	mu        sync.RWMutex
	factories map[string]func() T
	priority  []string
}

// RegistryOption configures a Registry.
type RegistryOption func(*registryConfig)

type registryConfig struct {
	priority []string
}

// WithPriority sets the priority order for backend selection.
// Backends listed first are preferred over backends listed later.
// Backends not in the list have lowest priority (in registration order).
func WithPriority(names ...string) RegistryOption {
	return func(c *registryConfig) {
		c.priority = names
	}
}

// NewRegistry creates a new Registry with optional configuration.
func NewRegistry[T any](opts ...RegistryOption) *Registry[T] {
	cfg := &registryConfig{}
	for _, opt := range opts {
		opt(cfg)
	}

	return &Registry[T]{
		factories: make(map[string]func() T),
		priority:  cfg.priority,
	}
}

// Register adds a factory for the given name.
// If a factory with the same name already exists, it is replaced.
// Thread-safe.
func (r *Registry[T]) Register(name string, factory func() T) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.factories[name] = factory
}

// Unregister removes the factory with the given name.
// Thread-safe.
func (r *Registry[T]) Unregister(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.factories, name)
}

// Get returns the factory output for the given name.
// Returns the zero value of T if not found.
// Thread-safe.
func (r *Registry[T]) Get(name string) T {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if factory, ok := r.factories[name]; ok {
		return factory()
	}

	var zero T
	return zero
}

// Has returns true if a factory with the given name is registered.
// Thread-safe.
func (r *Registry[T]) Has(name string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, ok := r.factories[name]
	return ok
}

// Best returns the highest-priority registered implementation.
// Returns the zero value of T if no implementations are registered.
// Thread-safe.
func (r *Registry[T]) Best() T {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Try priority list first
	for _, name := range r.priority {
		if factory, ok := r.factories[name]; ok {
			return factory()
		}
	}

	// Fall back to any registered factory
	for _, factory := range r.factories {
		return factory()
	}

	var zero T
	return zero
}

// BestName returns the name of the highest-priority registered implementation.
// Returns empty string if no implementations are registered.
// Thread-safe.
func (r *Registry[T]) BestName() string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Try priority list first
	for _, name := range r.priority {
		if _, ok := r.factories[name]; ok {
			return name
		}
	}

	// Fall back to any registered factory
	for name := range r.factories {
		return name
	}

	return ""
}

// Available returns all registered names.
// Thread-safe.
func (r *Registry[T]) Available() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.factories))
	for name := range r.factories {
		names = append(names, name)
	}
	return names
}

// Count returns the number of registered factories.
// Thread-safe.
func (r *Registry[T]) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.factories)
}
