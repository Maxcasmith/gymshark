package cors

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCors(t *testing.T) {
	handlerCalled := false
	handler := Cors([]string{"http://localhost:5173"}, http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		handlerCalled = true
	}))

	testCases := []struct {
		name                string
		method              string
		origin              string
		preflight           bool
		expectedAllowOrigin string
		expectedStatus      int
		expectHandlerCalled bool
	}{
		{"allowed origin", http.MethodGet, "http://localhost:5173", false, "http://localhost:5173", http.StatusOK, true},
		{"unknown origin", http.MethodGet, "http://evil.example", false, "", http.StatusOK, true},
		{"preflight", http.MethodOptions, "http://localhost:5173", true, "http://localhost:5173", http.StatusNoContent, false},
	}
	for _, testCase := range testCases {
		handlerCalled = false
		request := httptest.NewRequest(testCase.method, "/api/orders", nil)
		request.Header.Set("Origin", testCase.origin)
		if testCase.preflight {
			request.Header.Set("Access-Control-Request-Method", http.MethodPost)
		}
		recorder := httptest.NewRecorder()

		handler.ServeHTTP(recorder, request)

		if allowOrigin := recorder.Header().Get("Access-Control-Allow-Origin"); allowOrigin != testCase.expectedAllowOrigin {
			t.Errorf("%s: Access-Control-Allow-Origin = %q, expected %q", testCase.name, allowOrigin, testCase.expectedAllowOrigin)
		}
		if recorder.Code != testCase.expectedStatus {
			t.Errorf("%s: status = %d, expected %d", testCase.name, recorder.Code, testCase.expectedStatus)
		}
		if handlerCalled != testCase.expectHandlerCalled {
			t.Errorf("%s: handler called = %v, expected %v", testCase.name, handlerCalled, testCase.expectHandlerCalled)
		}
	}
}
