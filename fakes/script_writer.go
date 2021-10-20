package fakes

import "sync"

type ScriptWriter struct {
	WriteInitScriptCall struct {
		mutex     sync.Mutex
		CallCount int
		Receives  struct {
			ProfileDPath string
		}
		Returns struct {
			Error error
		}
		Stub func(string) error
	}
	WriteStartLoggingScriptCall struct {
		mutex     sync.Mutex
		CallCount int
		Receives  struct {
			ProfileDPath string
		}
		Returns struct {
			Error error
		}
		Stub func(string) error
	}
}

func (f *ScriptWriter) WriteInitScript(param1 string) error {
	f.WriteInitScriptCall.mutex.Lock()
	defer f.WriteInitScriptCall.mutex.Unlock()
	f.WriteInitScriptCall.CallCount++
	f.WriteInitScriptCall.Receives.ProfileDPath = param1
	if f.WriteInitScriptCall.Stub != nil {
		return f.WriteInitScriptCall.Stub(param1)
	}
	return f.WriteInitScriptCall.Returns.Error
}
func (f *ScriptWriter) WriteStartLoggingScript(param1 string) error {
	f.WriteStartLoggingScriptCall.mutex.Lock()
	defer f.WriteStartLoggingScriptCall.mutex.Unlock()
	f.WriteStartLoggingScriptCall.CallCount++
	f.WriteStartLoggingScriptCall.Receives.ProfileDPath = param1
	if f.WriteStartLoggingScriptCall.Stub != nil {
		return f.WriteStartLoggingScriptCall.Stub(param1)
	}
	return f.WriteStartLoggingScriptCall.Returns.Error
}
