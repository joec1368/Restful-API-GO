package test

import (
	"Dcard/controller"
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"unicode"
)

func makeToList(body []byte, t *testing.T) []map[string]string {
	var data []interface{}
	_ = json.Unmarshal(body, &data)
	var list []map[string]string
	for i := range data {
		m, ok := data[i].(map[string]interface{})
		if !ok {
			t.Errorf("want type map[string]interface{};  got %T", t)
		}
		temp := make(map[string]string)
		for k, v := range m {
			temp[k] = v.(string)
		}
		list = append(list, temp)
	}
	return list
}

func RemoveSpace(s string) string {
	rr := make([]rune, 0, len(s))
	for _, r := range s {
		if !unicode.IsSpace(r) {
			rr = append(rr, r)
		}
	}
	return string(rr)
}

func RemoveQuote(s string) string {
	rr := make([]rune, 0, len(s))
	for _, r := range s {
		if !unicode.Is(unicode.Quotation_Mark, r) {
			rr = append(rr, r)
		}
	}
	ans := string(rr)
	ans = strings.TrimRight(ans, "\n")
	return ans
}

func TestClearAll(t *testing.T) {

}

func init() {
	req, _ := http.NewRequest("GET", "/header/ClearAll", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controller.ClearAll)
	handler.ServeHTTP(rr, req)
	var expected string = "\"delete all data succeed\""
	if RemoveSpace(rr.Body.String()) != RemoveSpace(expected) {
		log.Fatal("init error")
	}
}

func TestAddHeader(t *testing.T) {
	jsonString := `[
     {
        "name":"a"
     },
     {
        "name":"b"
     },
     {
        "name":"c"
     }
   ]`
	req, err := http.NewRequest("POST", "/header", bytes.NewReader([]byte(jsonString)))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controller.AddHeader)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	var expected string = "\"insert a succeed\"" + "\"insert b succeed\"" + "\"insert c succeed\""
	if RemoveSpace(expected) != RemoveSpace(rr.Body.String()) {
		t.Errorf("error in TestAddHeader")
	}
}

func TestGetAllHeader(t *testing.T) {
	req, err := http.NewRequest("GET", "/header", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controller.GetAllHeader)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	body, err := ioutil.ReadAll(rr.Body)
	list := makeToList(body, t)
	if len(list) != 3 {
		t.Errorf("error in TestGetAllHeader")
	}
	if list[0]["id"] != "a" {
		t.Errorf("error in TestGetAllHeader")
	} else if list[1]["id"] != "b" {
		t.Errorf("error in TestGetAllHeader")
	} else if list[2]["id"] != "c" {
		t.Errorf("error in TestGetAllHeader")
	}
}

func TestGetHead(t *testing.T) {
	req, err := http.NewRequest("POST", "/header/GetHead", nil)
	if err != nil {
		t.Fatal(err)
	}

	vars := map[string]string{
		"header": "a",
	}

	req = mux.SetURLVars(req, vars)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controller.GetCertainHeader)

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	if RemoveSpace(rr.Body.String()) != "\"\"" {
		t.Errorf("TestGetHead")
	}
}

func TestDeleteHeader(t *testing.T) {
	req, err := http.NewRequest("DELETE", "/header", nil)
	if err != nil {
		t.Fatal(err)
	}

	vars := map[string]string{
		"header": "c",
	}

	req = mux.SetURLVars(req, vars)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controller.DeleteHeader)

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	req, err = http.NewRequest("GET", "/header", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(controller.GetAllHeader)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	body, err := ioutil.ReadAll(rr.Body)
	list := makeToList(body, t)
	if len(list) != 2 {
		t.Errorf("error in TestDeleteHeader")
	}
	if list[0]["id"] != "a" {
		t.Errorf("error in TestDeleteHeader")
	} else if list[1]["id"] != "b" {
		t.Errorf("error in TestDeleteHeader")
	}

}

func TestChangeHeaderName(t *testing.T) {
	jsonString := `[
     {
        "newHeader":"d"
     }
   ]`
	req, err := http.NewRequest("PUT", "/header", bytes.NewReader([]byte(jsonString)))
	if err != nil {
		t.Fatal(err)
	}

	vars := map[string]string{
		"header": "a",
	}

	req = mux.SetURLVars(req, vars)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controller.UpdateHeader)

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	req, err = http.NewRequest("GET", "/header", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(controller.GetAllHeader)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	body, err := ioutil.ReadAll(rr.Body)
	list := makeToList(body, t)
	if len(list) != 2 {
		t.Errorf("error in TestChangeHeaderName")
	}

	if list[0]["id"] != "b" {
		t.Errorf("error in TestChangeHeaderName")
	} else if list[1]["id"] != "d" {
		t.Errorf("error in TestChangeHeaderName")
	}

}
