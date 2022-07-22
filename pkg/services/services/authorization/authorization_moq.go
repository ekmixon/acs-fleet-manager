// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package authorization

import (
	"sync"
)

// Ensure, that AuthorizationMock does implement Authorization.
// If this is not the case, regenerate this file with moq.
var _ Authorization = &AuthorizationMock{}

// AuthorizationMock is a mock implementation of Authorization.
//
// 	func TestSomethingThatUsesAuthorization(t *testing.T) {
//
// 		// make and configure a mocked Authorization
// 		mockedAuthorization := &AuthorizationMock{
// 			CheckUserValidFunc: func(username string, orgID string) (bool, error) {
// 				panic("mock out the CheckUserValid method")
// 			},
// 		}
//
// 		// use mockedAuthorization in code that requires Authorization
// 		// and then make assertions.
//
// 	}
type AuthorizationMock struct {
	// CheckUserValidFunc mocks the CheckUserValid method.
	CheckUserValidFunc func(username string, orgID string) (bool, error)

	// calls tracks calls to the methods.
	calls struct {
		// CheckUserValid holds details about calls to the CheckUserValid method.
		CheckUserValid []struct {
			// Username is the username argument value.
			Username string
			// OrgID is the orgID argument value.
			OrgID string
		}
	}
	lockCheckUserValid sync.RWMutex
}

// CheckUserValid calls CheckUserValidFunc.
func (mock *AuthorizationMock) CheckUserValid(username string, orgID string) (bool, error) {
	if mock.CheckUserValidFunc == nil {
		panic("AuthorizationMock.CheckUserValidFunc: method is nil but Authorization.CheckUserValid was just called")
	}
	callInfo := struct {
		Username string
		OrgID    string
	}{
		Username: username,
		OrgID:    orgID,
	}
	mock.lockCheckUserValid.Lock()
	mock.calls.CheckUserValid = append(mock.calls.CheckUserValid, callInfo)
	mock.lockCheckUserValid.Unlock()
	return mock.CheckUserValidFunc(username, orgID)
}

// CheckUserValidCalls gets all the calls that were made to CheckUserValid.
// Check the length with:
//     len(mockedAuthorization.CheckUserValidCalls())
func (mock *AuthorizationMock) CheckUserValidCalls() []struct {
	Username string
	OrgID    string
} {
	var calls []struct {
		Username string
		OrgID    string
	}
	mock.lockCheckUserValid.RLock()
	calls = mock.calls.CheckUserValid
	mock.lockCheckUserValid.RUnlock()
	return calls
}