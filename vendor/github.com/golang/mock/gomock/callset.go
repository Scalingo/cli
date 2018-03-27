// Copyright 2011 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gomock

<<<<<<< HEAD
import (
	"bytes"
	"fmt"
)

// callSet represents a set of expected calls, indexed by receiver and method
// name.
type callSet struct {
	// Calls that are still expected.
	expected map[callSetKey][]*Call
	// Calls that have been exhausted.
	exhausted map[callSetKey][]*Call
}

// callSetKey is the key in the maps in callSet
type callSetKey struct {
	receiver interface{}
	fname    string
}

func newCallSet() *callSet {
	return &callSet{make(map[callSetKey][]*Call), make(map[callSetKey][]*Call)}
}

// Add adds a new expected call.
func (cs callSet) Add(call *Call) {
	key := callSetKey{call.receiver, call.method}
	m := cs.expected
	if call.exhausted() {
		m = cs.exhausted
	}
	m[key] = append(m[key], call)
=======
// callSet represents a set of expected calls, indexed by receiver and method
// name.
type callSet map[interface{}]map[string][]*Call

// Add adds a new expected call.
func (cs callSet) Add(call *Call) {
	methodMap, ok := cs[call.receiver]
	if !ok {
		methodMap = make(map[string][]*Call)
		cs[call.receiver] = methodMap
	}
	methodMap[call.method] = append(methodMap[call.method], call)
>>>>>>> Update dependencies
}

// Remove removes an expected call.
func (cs callSet) Remove(call *Call) {
<<<<<<< HEAD
	key := callSetKey{call.receiver, call.method}
	calls := cs.expected[key]
	for i, c := range calls {
		if c == call {
			// maintain order for remaining calls
			cs.expected[key] = append(calls[:i], calls[i+1:]...)
			cs.exhausted[key] = append(cs.exhausted[key], call)
=======
	methodMap, ok := cs[call.receiver]
	if !ok {
		return
	}
	sl := methodMap[call.method]
	for i, c := range sl {
		if c == call {
			// quick removal; we don't need to maintain call order
			if len(sl) > 1 {
				sl[i] = sl[len(sl)-1]
			}
			methodMap[call.method] = sl[:len(sl)-1]
>>>>>>> Update dependencies
			break
		}
	}
}

<<<<<<< HEAD
// FindMatch searches for a matching call. Returns error with explanation message if no call matched.
func (cs callSet) FindMatch(receiver interface{}, method string, args []interface{}) (*Call, error) {
	key := callSetKey{receiver, method}

	// Search through the expected calls.
	expected := cs.expected[key]
	var callsErrors bytes.Buffer
	for _, call := range expected {
		err := call.matches(args)
		if err != nil {
			fmt.Fprintf(&callsErrors, "\n%v", err)
		} else {
			return call, nil
		}
	}

	// If we haven't found a match then search through the exhausted calls so we
	// get useful error messages.
	exhausted := cs.exhausted[key]
	for _, call := range exhausted {
		if err := call.matches(args); err != nil {
			fmt.Fprintf(&callsErrors, "\n%v", err)
		}
	}

	if len(expected)+len(exhausted) == 0 {
		fmt.Fprintf(&callsErrors, "there are no expected calls of the method %q for that receiver", method)
	}

	return nil, fmt.Errorf(callsErrors.String())
}

// Failures returns the calls that are not satisfied.
func (cs callSet) Failures() []*Call {
	failures := make([]*Call, 0, len(cs.expected))
	for _, calls := range cs.expected {
		for _, call := range calls {
			if !call.satisfied() {
				failures = append(failures, call)
			}
		}
	}
	return failures
=======
// FindMatch searches for a matching call. Returns nil if no call matched.
func (cs callSet) FindMatch(receiver interface{}, method string, args []interface{}) *Call {
	methodMap, ok := cs[receiver]
	if !ok {
		return nil
	}
	calls, ok := methodMap[method]
	if !ok {
		return nil
	}

	// Search through the unordered set of calls expected on a method on a
	// receiver.
	for _, call := range calls {
		// A call should not normally still be here if exhausted,
		// but it can happen if, for instance, .Times(0) was used.
		// Pretend the call doesn't match.
		if call.exhausted() {
			continue
		}
		if call.matches(args) {
			return call
		}
	}

	return nil
>>>>>>> Update dependencies
}
