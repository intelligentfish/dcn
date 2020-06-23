package serviceSorter

import (
	"github.com/intelligentfish/dcn/types"
	"sort"
)

// SortKey sort key
type SortKey int

const (
	SortKeyStartupPriority  = SortKey(iota) // sorted by startup priority
	SortKeyShutdownPriority                 // sorted by shutdown priority
)

// ServiceSorter service sorter
type ServiceSorter struct {
	sortKey SortKey          // sort key
	list    []types.IService // services
}

// New factory method
func New(list []types.IService) ServiceSorter {
	return ServiceSorter{
		sortKey: SortKeyStartupPriority,
		list:    list,
	}
}

// Len length of list
func (object ServiceSorter) Len() int {
	return len(object.list)
}

// Swap swap service
func (object ServiceSorter) Swap(i, j int) {
	object.list[i], object.list[j] = object.list[j], object.list[i]
}

// Less compare function
func (object ServiceSorter) Less(i, j int) bool {
	switch object.sortKey {
	case SortKeyShutdownPriority:
		return object.list[i].GetShutdownPriority() < object.list[j].GetShutdownPriority()
	}
	return object.list[i].GetStartupPriority() < object.list[j].GetStartupPriority()
}

// SortByStartupPriority sort services by startup priority
func (object ServiceSorter) SortByStartupPriority() *ServiceSorter {
	object.sortKey = SortKeyStartupPriority
	sort.Sort(object)
	return &object
}

// SortByShutdownPriority sort services by shutdown priority
func (object ServiceSorter) SortByShutdownPriority() *ServiceSorter {
	object.sortKey = SortKeyShutdownPriority
	sort.Sort(object)
	return &object
}

// Foreach iterator services
func (object *ServiceSorter) Foreach(callback func(srv types.IService) (ok bool)) {
	for _, srv := range object.list {
		if !callback(srv) {
			break
		}
	}
}
