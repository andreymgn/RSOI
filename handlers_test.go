package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var s = newServer()

func TestHello(t *testing.T) {
	req, err := http.NewRequest("GET", "/hello/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(hello)
	handler.ServeHTTP(rr, req)
	resp := rr.Result()

	if status := resp.StatusCode; status != http.StatusOK {
		t.Errorf("handler returned unexpected status code: got %v want %v",
			status, http.StatusOK)
	}

	expectedCtype := "text/plain"
	if ctype := resp.Header.Get("Content-Type"); ctype != expectedCtype {
		t.Errorf("handler returned unexpected Content-Type: got %v want %v",
			ctype, expectedCtype)
	}

	expectedBody := []byte("Hello, world!")
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if bytes.Compare(body, expectedBody) != 0 {
		t.Errorf("handler returned unexpected body: got %v want %v", body, expectedBody)
	}
}

func TestAdd(t *testing.T) {
	json := []byte(`{"lhs":1,"rhs":2}`)
	req, err := http.NewRequest("GET", "/add/", bytes.NewBuffer(json))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(add)
	handler.ServeHTTP(rr, req)
	resp := rr.Result()

	if status := resp.StatusCode; status != http.StatusOK {
		t.Errorf("handler returned unexpected status code: got %v want %v",
			status, http.StatusOK)
	}

	expectedCtype := "application/json"
	if ctype := resp.Header.Get("Content-Type"); ctype != expectedCtype {
		t.Errorf("handler returned unexpected Content-Type: got %v want %v",
			ctype, expectedCtype)
	}

	expectedBody := []byte(`{"result":3}`)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if bytes.Compare(body, expectedBody) != 0 {
		t.Errorf("handler returned unexpected body: got %v want %v", string(body[:]), string(expectedBody[:]))
	}

}
