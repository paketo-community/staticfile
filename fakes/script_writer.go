package fakes

import "sync"

type ScriptWriter struct {
	WriteInitScriptCall struct {
		sync.Mutex
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
		sync.Mutex
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
	f.WriteInitScriptCall.Lock()
	defer f.WriteInitScriptCall.Unlock()
	f.WriteInitScriptCall.CallCount++
	f.WriteInitScriptCall.Receives.ProfileDPath = param1
	if f.WriteInitScriptCall.Stub != nil {
		return f.WriteInitScriptCall.Stub(param1)
	}
	return f.WriteInitScriptCall.Returns.Error
}
func (f *ScriptWriter) WriteStartLoggingScript(param1 string) error {
	f.WriteStartLoggingScriptCall.Lock()
	defer f.WriteStartLoggingScriptCall.Unlock()
	f.WriteStartLoggingScriptCall.CallCount++
	f.WriteStartLoggingScriptCall.Receives.ProfileDPath = param1
	if f.WriteStartLoggingScriptCall.Stub != nil {
		return f.WriteStartLoggingScriptCall.Stub(param1)
	}
	return f.WriteStartLoggingScriptCall.Returns.Error
}
