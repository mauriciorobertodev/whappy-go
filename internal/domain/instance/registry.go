package instance

type InstanceRegistry interface {
	Add(instance *Instance)
	Get(id string) (*Instance, bool)
	Remove(id string)
}
