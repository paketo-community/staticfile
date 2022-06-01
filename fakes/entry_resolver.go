package fakes

import (
	"sync"

	"github.com/paketo-buildpacks/packit/v2"
)

type EntryResolver struct {
	ResolveCall struct {
		mutex     sync.Mutex
		CallCount int
		Receives  struct {
			Name       string
			Entries    []packit.BuildpackPlanEntry
			Priorities []interface {
			}
		}
		Returns struct {
			BuildpackPlanEntry      packit.BuildpackPlanEntry
			BuildpackPlanEntrySlice []packit.BuildpackPlanEntry
		}
		Stub func(string, []packit.BuildpackPlanEntry, []interface {
		}) (packit.BuildpackPlanEntry, []packit.BuildpackPlanEntry)
	}
}

func (f *EntryResolver) Resolve(param1 string, param2 []packit.BuildpackPlanEntry, param3 []interface {
}) (packit.BuildpackPlanEntry, []packit.BuildpackPlanEntry) {
	f.ResolveCall.mutex.Lock()
	defer f.ResolveCall.mutex.Unlock()
	f.ResolveCall.CallCount++
	f.ResolveCall.Receives.Name = param1
	f.ResolveCall.Receives.Entries = param2
	f.ResolveCall.Receives.Priorities = param3
	if f.ResolveCall.Stub != nil {
		return f.ResolveCall.Stub(param1, param2, param3)
	}
	return f.ResolveCall.Returns.BuildpackPlanEntry, f.ResolveCall.Returns.BuildpackPlanEntrySlice
}
