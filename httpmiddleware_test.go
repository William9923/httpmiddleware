package httpmiddleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/William9923/httpmiddleware"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

func Test_Use(t *testing.T) {

	dummyMiddleware := func(fn httprouter.Handle) httprouter.Handle {
		return fn
	}

	t.Run("Add single middleware", func(t *testing.T) {

		assert := assert.New(t)
		middlewareGroup := httpmiddleware.New()

		assert.Equal(middlewareGroup.Count(), 0, "Middleware count not same!")

		middlewareGroup.Use(dummyMiddleware)
		assert.Equal(middlewareGroup.Count(), 1, "Middleware count not same!")
	})

	t.Run("Add Multiple middleware simultaneously", func(t *testing.T) {
		assert := assert.New(t)
		middlewareGroup := httpmiddleware.New()

		assert.Equal(middlewareGroup.Count(), 0, "Middleware count not same!")

		middlewareGroup.Use(dummyMiddleware, dummyMiddleware, dummyMiddleware)
		assert.Equal(middlewareGroup.Count(), 3, "Middleware count not same!")

		middlewareGroup.Use(dummyMiddleware)
		middlewareGroup.Use(dummyMiddleware)
		assert.Equal(middlewareGroup.Count(), 5, "Middleware count not same!")
	})
}

func Test_Wrap(t *testing.T) {

	t.Run("Wrap middleware to http handler", func(t *testing.T) {

		var called bool
		middleware := func(fn httprouter.Handle) httprouter.Handle {
			called = true
			return fn
		}

		dummyHandler := func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {}

		middlewareGroup := httpmiddleware.New()
		middlewareGroup.Use(middleware)

		router := httprouter.New()
		router.GET("/foo", middlewareGroup.Wrap(dummyHandler))

		req, _ := http.NewRequest("GET", "/foo", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)
		if status := w.Code; status != http.StatusOK {
			t.Error("Wrong status")
		}

		if !called {
			t.Error("Middleware not called!")
		}
	})

	t.Run("Wrap up multiple middleware to http handler (with ordering)", func(t *testing.T) {
		callStack := []int{}
		expectedCallStack := []int{4, 3, 2, 1}

		assert := assert.New(t)

		firstMiddleware := func(fn httprouter.Handle) httprouter.Handle {
			callStack = append(callStack, 1)
			return fn
		}

		secondMiddleware := func(fn httprouter.Handle) httprouter.Handle {
			callStack = append(callStack, 2)
			return fn
		}

		thirdMiddleware := func(fn httprouter.Handle) httprouter.Handle {
			callStack = append(callStack, 3)
			return fn
		}

		forthMiddleware := func(fn httprouter.Handle) httprouter.Handle {
			callStack = append(callStack, 4)
			return fn
		}

		dummyHandler := func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {}

		middlewareGroup := httpmiddleware.New()
		middlewareGroup.Use(firstMiddleware, secondMiddleware, thirdMiddleware, forthMiddleware)

		router := httprouter.New()
		router.GET("/foo", middlewareGroup.Wrap(dummyHandler))

		req, _ := http.NewRequest("GET", "/foo", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)
		if status := w.Code; status != http.StatusOK {
			t.Error("Wrong status")
		}

		assert.EqualValues(expectedCallStack, callStack, "Different call stack order!")

	})

	t.Run("Wrap up empty middleware to http handler", func(t *testing.T) {
		var called bool
		dummyHandler := func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			called = true
		}

		middlewareGroup := httpmiddleware.New()

		router := httprouter.New()
		router.GET("/foo", middlewareGroup.Wrap(dummyHandler))

		req, _ := http.NewRequest("GET", "/foo", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)
		if status := w.Code; status != http.StatusOK {
			t.Error("Wrong status")
		}

		if !called {
			t.Error("Handler not called!")
		}
	})
}
