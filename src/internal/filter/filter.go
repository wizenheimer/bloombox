package filter

// Filter interface for different filtering implementations
type Filter interface {
	Add(item string) error
	Contains(item string) bool
	Size() int
	Delete(item string) bool
}

// ExtendedFilter interface for filters that support additional operations
type ExtendedFilter interface {
	Filter
	Clear()
	GetAll() []string
}
