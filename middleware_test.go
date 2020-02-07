package middleware

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type (
	middleA struct{}
	middleB struct{}
	middleC struct{}
)

func (a middleA) Handle(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("middleA "))
		h.ServeHTTP(w, r)
	})

}
func (a middleB) Handle(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("middleB "))
		h.ServeHTTP(w, r)
	})
}
func (a middleC) Handle(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("middleC "))
		h.ServeHTTP(w, r)
	})
}

func TestMiddlewareNew(t *testing.T) {
	m := New()
	if len(m.middlewares) != 0 {
		t.Errorf("unexpected result: [got] %v [want] %v", len(m.middlewares), 0)
	}

	m = New(middleA{}, middleB{}, middleC{})
	if len(m.middlewares) != 3 {
		t.Errorf("unexpected result: [got] %v [want] %v", len(m.middlewares), 3)
	}
}

func TestMiddlewareAppend(t *testing.T) {
	m := New(middleA{})
	m = m.Append(middleB{})
	if len(m.middlewares) != 2 {
		t.Errorf("unexpected result: [got] %v [want] %v", len(m.middlewares), 2)
	}
	m = m.Append(middleC{})
	if len(m.middlewares) != 3 {
		t.Errorf("unexpected result: [got] %v [want] %v", len(m.middlewares), 3)
	}
}

func TestMiddlewareApplyEmpty(t *testing.T) {
	handler := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("test"))
		})
	recorder := httptest.NewRecorder()
	m := New()
	m.Apply(handler).ServeHTTP(recorder, nil)

	resp := recorder.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("unexpected result: [got] %v [want] %v", resp.StatusCode, http.StatusOK)
	}

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("unexpected result: [got] %v [want] %v", err, nil)
	}
	defer resp.Body.Close()

	body := string(buf)
	expected := "test"
	if body != expected {
		t.Errorf("unexpected result: [got] %v [want] %v", body, expected)
	}
}

func TestMiddlewareApply(t *testing.T) {
	handler := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("test"))
		})
	recorder := httptest.NewRecorder()
	m := New(middleA{}, middleB{}, middleC{})
	m.Apply(handler).ServeHTTP(recorder, nil)

	resp := recorder.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("unexpected result: [got] %v [want] %v", resp.StatusCode, http.StatusOK)
	}

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("unexpected result: [got] %v [want] %v", err, nil)
	}
	defer resp.Body.Close()

	body := string(buf)
	expected := "middleA middleB middleC test"
	if body != expected {
		t.Errorf("unexpected result: [got] %v [want] %v", body, expected)
	}
}
