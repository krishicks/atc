// This file was generated by counterfeiter
package buildsfakes

import (
	"sync"

	"github.com/concourse/atc/builds"
)

type FakeATCListener struct {
	ListenStub        func(channel string) (chan bool, error)
	listenMutex       sync.RWMutex
	listenArgsForCall []struct {
		channel string
	}
	listenReturns struct {
		result1 chan bool
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeATCListener) Listen(channel string) (chan bool, error) {
	fake.listenMutex.Lock()
	fake.listenArgsForCall = append(fake.listenArgsForCall, struct {
		channel string
	}{channel})
	fake.recordInvocation("Listen", []interface{}{channel})
	fake.listenMutex.Unlock()
	if fake.ListenStub != nil {
		return fake.ListenStub(channel)
	} else {
		return fake.listenReturns.result1, fake.listenReturns.result2
	}
}

func (fake *FakeATCListener) ListenCallCount() int {
	fake.listenMutex.RLock()
	defer fake.listenMutex.RUnlock()
	return len(fake.listenArgsForCall)
}

func (fake *FakeATCListener) ListenArgsForCall(i int) string {
	fake.listenMutex.RLock()
	defer fake.listenMutex.RUnlock()
	return fake.listenArgsForCall[i].channel
}

func (fake *FakeATCListener) ListenReturns(result1 chan bool, result2 error) {
	fake.ListenStub = nil
	fake.listenReturns = struct {
		result1 chan bool
		result2 error
	}{result1, result2}
}

func (fake *FakeATCListener) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.listenMutex.RLock()
	defer fake.listenMutex.RUnlock()
	return fake.invocations
}

func (fake *FakeATCListener) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ builds.ATCListener = new(FakeATCListener)
