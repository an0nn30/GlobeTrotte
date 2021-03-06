package router

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	structs "github.com/binhonglee/GlobeTrotte/src/turbine/structs"
)

func TestPasswd(t *testing.T) {
	request, _ := http.NewRequest("GET", API_PREFIX+"/passwd", nil)

	response := httptest.NewRecorder()
	NewRouter().ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Errorf("/passwd, expected code 200 but received " + strconv.Itoa(response.Code) + " instead.")
		return
	}
}

func addTest(
	path string,
	t *testing.T,
	toAdd structs.IStructs,
	objType interface{},
	expectedCode int,
) {
	res := httptest.NewRecorder()
	data, err := json.Marshal(toAdd)
	req, _ := http.NewRequest(
		"POST", API_PREFIX+path, bytes.NewBuffer(data))
	NewRouter().ServeHTTP(res, req)
	if res.Code != expectedCode {
		t.Errorf("addTest(), expected code " +
			strconv.Itoa(expectedCode) +
			" but received " + strconv.Itoa(res.Code) +
			" instead.")
	}

	if expectedCode == http.StatusCreated {
		data, err = ioutil.ReadAll(res.Body)
		if err != nil {
			t.Errorf("Adding failed.")
		}
		json.Unmarshal(data, objType)
	}
}

func getTest(path string, t *testing.T, objType interface{}, expectedCode int) {
	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", API_PREFIX+path, nil)
	NewRouter().ServeHTTP(res, req)
	if res.Code != expectedCode {
		t.Errorf("getTest(), expected code " +
			strconv.Itoa(expectedCode) + " but received " +
			strconv.Itoa(res.Code) + " instead.")
		return
	}

	if expectedCode == http.StatusOK {
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Errorf("Adding failed.")
		}
		json.Unmarshal(data, objType)
	}
}

func updateTest(
	path string,
	t *testing.T,
	toAdd structs.IStructs,
	expectedCode int,
) {
	res := httptest.NewRecorder()
	data, _ := json.Marshal(toAdd)
	req, _ := http.NewRequest(
		"POST",
		API_PREFIX+path+strconv.Itoa(toAdd.GetID()),
		bytes.NewBuffer(data),
	)
	NewRouter().ServeHTTP(res, req)
	if res.Code != expectedCode {
		t.Errorf("updateTest(), expected code " +
			strconv.Itoa(expectedCode) + " but received " +
			strconv.Itoa(res.Code) + " instead.")
		return
	}
}

func deleteTest(path string, t *testing.T, expectedCode int) {
	res := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", API_PREFIX+path, nil)
	NewRouter().ServeHTTP(res, req)
	if res.Code != expectedCode {
		t.Errorf("deleteTest(), expected code " + strconv.Itoa(expectedCode) + " but received " + strconv.Itoa(res.Code) + " instead.")
	}
}
