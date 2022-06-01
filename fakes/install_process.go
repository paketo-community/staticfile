package fakes

import (
	"sync"

	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-community/staticfile"
)

type InstallProcess struct {
	ExecuteCall struct {
		mutex     sync.Mutex
		CallCount int
		Receives  struct {
			Context     packit.BuildContext
			TemplConfig staticfile.Config
		}
		Returns struct {
			Error error
		}
		Stub func(packit.BuildContext, staticfile.Config) error
	}
}

func (f *InstallProcess) Execute(param1 packit.BuildContext, param2 staticfile.Config) error {
	f.ExecuteCall.mutex.Lock()
	defer f.ExecuteCall.mutex.Unlock()
	f.ExecuteCall.CallCount++
	f.ExecuteCall.Receives.Context = param1
	f.ExecuteCall.Receives.TemplConfig = param2
	if f.ExecuteCall.Stub != nil {
		return f.ExecuteCall.Stub(param1, param2)
	}
	return f.ExecuteCall.Returns.Error
}
