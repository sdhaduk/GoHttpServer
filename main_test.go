package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleRoot(t *testing.T) {
	w := httptest.NewRecorder()

	handleRoot(w, nil)

	desiredCode := http.StatusOK

	if w.Code != desiredCode {
		t.Errorf("bad response code, expected %v, but got %v\nbody: %v", desiredCode, w.Code, w.Body.String())
	}

	expectedMessage := []byte("Welcome to our Homepage!\n")

	if !bytes.Equal(expectedMessage, w.Body.Bytes()) {
		t.Errorf("bad return, got %v, expected %v", w.Body.String(), string(expectedMessage))
	}
}

func TestHandleGoodbye(t *testing.T) {
	w := httptest.NewRecorder()

	handleGoodbye(w, nil)

	desiredCode := http.StatusOK

	if w.Code != desiredCode {
		t.Errorf("bad response code, expected %v, but got %v\nbody: %v", desiredCode, w.Code, w.Body.String())
	}

	expectedMessage := []byte("Goodbye, World!\n")
	if !bytes.Equal(expectedMessage, w.Body.Bytes()) {
		t.Errorf("bad return, got %v, expected %v", w.Body.String(), string(expectedMessage))
	}
}

func TestHandleHelloParameterized(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/hello?user=TestMan", nil)

	w := httptest.NewRecorder()

	handleHelloParameterized(w, req)
	
	desiredCode := http.StatusOK

	if w.Code != desiredCode {
		t.Errorf("bad response code, expected %v, but got %v\nbody: %v", desiredCode, w.Code, w.Body.String())
	}
	
	expectedMessage := []byte("Hello, TestMan!\n")
	if !bytes.Equal(expectedMessage, w.Body.Bytes()) {
		t.Errorf("bad return, got %v, expected %v", w.Body.String(), string(expectedMessage))
	}
}

func TestHandleHelloParameterizedNoParam(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/hello/", nil)
	w := httptest.NewRecorder()

	handleHelloParameterized(w, req)
	desiredCode := http.StatusOK

	if w.Code != desiredCode {
		t.Errorf("bad response code, expected %v, but got %v\nbody: %v", desiredCode, w.Code, w.Body.String())
	}
	
	expectedMessage := []byte("Hello, User!\n")
	if !bytes.Equal(expectedMessage, w.Body.Bytes()) {
		t.Errorf("bad return, got %v, expected %v", w.Body.String(), string(expectedMessage))
	}
}
func TestHandleHelloParameterizedWrongParam(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/hello?foo=bar", nil)
	w := httptest.NewRecorder()

	handleHelloParameterized(w, req)
	desiredCode := http.StatusOK

	if w.Code != desiredCode {
		t.Errorf("bad response code, expected %v, but got %v\nbody: %v", desiredCode, w.Code, w.Body.String())
	}

	expectedMessage := []byte("Hello, User!\n")
	if !bytes.Equal(expectedMessage, w.Body.Bytes()) {
		t.Errorf("bad return, got %v, expected %v", w.Body.String(), string(expectedMessage))
	}
}