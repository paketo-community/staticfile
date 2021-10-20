package fakes

import (
	"sync"

	"github.com/paketo-community/staticfile"
)

type BpYMLParser struct {
	ParseCall struct {
		mutex     sync.Mutex
		CallCount int
		Receives  struct {
			Path string
		}
		Returns struct {
			Config staticfile.Config
			Err    error
		}
		Stub func(string) (staticfile.Config, error)
	}
	ValidConfigCall struct {
		mutex     sync.Mutex
		CallCount int
		Receives  struct {
			Path string
		}
		Returns struct {
			Valid bool
			Err   error
		}
		Stub func(string) (bool, error)
	}
}

func (f *BpYMLParser) Parse(param1 string) (staticfile.Config, error) {
	f.ParseCall.mutex.Lock()
	defer f.ParseCall.mutex.Unlock()
	f.ParseCall.CallCount++
	f.ParseCall.Receives.Path = param1
	if f.ParseCall.Stub != nil {
		return f.ParseCall.Stub(param1)
	}
	return f.ParseCall.Returns.Config, f.ParseCall.Returns.Err
}
func (f *BpYMLParser) ValidConfig(param1 string) (bool, error) {
	f.ValidConfigCall.mutex.Lock()
	defer f.ValidConfigCall.mutex.Unlock()
	f.ValidConfigCall.CallCount++
	f.ValidConfigCall.Receives.Path = param1
	if f.ValidConfigCall.Stub != nil {
		return f.ValidConfigCall.Stub(param1)
	}
	return f.ValidConfigCall.Returns.Valid, f.ValidConfigCall.Returns.Err
}
