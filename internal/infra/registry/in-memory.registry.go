package registry

import (
	"sync"

	"github.com/mauriciorobertodev/whappy-go/internal/domain/instance"
)

type InMemoryInstanceRegistry struct {
	mu        sync.RWMutex
	instances map[string]*instance.Instance
}

func NewInMemoryInstanceRegistry() *InMemoryInstanceRegistry {
	return &InMemoryInstanceRegistry{
		instances: make(map[string]*instance.Instance),
	}
}

func (r *InMemoryInstanceRegistry) Add(i *instance.Instance) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.instances[i.ID] = i
}

func (r *InMemoryInstanceRegistry) Get(id string) (*instance.Instance, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	i, ok := r.instances[id]
	return i, ok
}

func (r *InMemoryInstanceRegistry) Remove(id string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.instances, id)
}
