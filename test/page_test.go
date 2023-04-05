package test

import (
	"Dcard/controller"
	"Dcard/service"
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func init() {
	req, _ := http.NewRequest("GET", "/header/Clear", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controller.ClearPage)
	handler.ServeHTTP(rr, req)
	var expected string = "\"delete all data succeed\""
	if RemoveSpace(rr.Body.String()) != RemoveSpace(expected) {
		log.Fatal("init error")
	}
}

func TestAddPage(t *testing.T) {
	jsonString := `[
     {
        "article":"0"
     },
     {
        "article":"1"
     },
     {
        "article":"2"
     }
   ]`
	req, err := http.NewRequest("POST", "/page/setPage/", bytes.NewReader([]byte(jsonString)))
	if err != nil {
		t.Fatal(err)
	}
	vars := map[string]string{
		"header": "b",
	}

	req = mux.SetURLVars(req, vars)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controller.AddPage)

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	if RemoveSpace("") != RemoveSpace(rr.Body.String()) {
		t.Errorf("error in TestAddHeader")
	}
}

func iterate(next string, t *testing.T) []map[string]string {
	var allList []map[string]string
	for next != "" {
		jsonString := `[{"id": "` + next + `"}]`
		req, err := http.NewRequest("POST", "/page/getPage/", bytes.NewReader([]byte(jsonString)))
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(controller.GetPage)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}
		var x map[string]interface{}
		json.Unmarshal([]byte(rr.Body.String()), &x)
		node := make(map[string]string)
		node["Article"] = x["Article"].(string)
		node["Next"] = x["Next"].(string)
		next = x["Next"].(string)
		allList = append(allList, node)
	}
	return allList
}

func TestIteratePage(t *testing.T) {
	req, err := http.NewRequest("POST", "/header/GetHead", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	var target string = "b"
	vars := map[string]string{
		"header": target,
	}
	req = mux.SetURLVars(req, vars)
	handler := http.HandlerFunc(controller.GetCertainHeader)
	handler.ServeHTTP(rr, req)
	bodyBytes, err := io.ReadAll(rr.Body)
	if err != nil {
		t.Errorf("TestIteratePage")
		log.Fatal(err)
	}
	var next string = RemoveQuote(string(bodyBytes))
	if next == "" {
		t.Errorf("TestIteratePage")
		log.Fatal("TestIteratePage error")
	} else {
		var allList []map[string]string
		allList = iterate(next, t)
		if len(allList) != 3 {
			t.Errorf("TestIteratePage")
		}
		if allList[0]["Article"] != "0" && allList[0]["Next"] != "" {
			t.Errorf("TestIteratePage")
		} else if allList[1]["Article"] != "1" && allList[1]["Next"] != "" {
			t.Errorf("TestIteratePage")
		} else if allList[2]["Article"] != "2" && allList[2]["Next"] == "" {
			t.Errorf("TestIteratePage")
		}

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
		for i := range list {
			if list[i]["id"] == target {
				page := service.GetPage(allList[len(allList)-2]["Next"])
				if list[i]["final"] != page.ID {
					t.Errorf("TestIteratePage")
				}
			}
		}
	}
}

func updatePage(origin string, new string, t *testing.T) {
	jsonString := `[{"originArticle":"` + origin + `", "newArticle":"` + new + `"}]`
	req, err := http.NewRequest("PUT", "/header", bytes.NewReader([]byte(jsonString)))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	var target string = "b"
	vars := map[string]string{
		"header": target,
	}
	req = mux.SetURLVars(req, vars)
	handler := http.HandlerFunc(controller.UpdatePage)
	handler.ServeHTTP(rr, req)
}

func TestUpdatePageAndIterate(t *testing.T) {
	updatePage("0", "10", t)
	updatePage("1", "11", t)
	updatePage("2", "12", t)
	req, err := http.NewRequest("POST", "/header/GetHead", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	var target string = "b"
	vars := map[string]string{
		"header": target,
	}
	req = mux.SetURLVars(req, vars)
	handler := http.HandlerFunc(controller.GetCertainHeader)
	handler.ServeHTTP(rr, req)
	bodyBytes, err := io.ReadAll(rr.Body)
	if err != nil {
		t.Errorf("TestUpdatePageAndIterate")
		log.Fatal(err)
	}
	var next string = RemoveQuote(string(bodyBytes))
	if next == "" {
		t.Errorf("TestUpdatePageAndIterate")
		log.Fatal("TestUpdatePageAndIterate in next")
	} else {
		var allList []map[string]string
		allList = iterate(next, t)
		if len(allList) != 3 {
			t.Errorf("TestUpdatePageAndIterate in total 3")
		}
		if allList[0]["Article"] != "10" && allList[0]["Next"] != "" {
			t.Errorf("TestUpdatePageAndIterate in compare 10")
		} else if allList[1]["Article"] != "11" && allList[1]["Next"] != "" {
			t.Errorf("TestUpdatePageAndIterate in compare 11")
		} else if allList[2]["Article"] != "12" && allList[2]["Next"] == "" {
			t.Errorf("TestUpdatePageAndIterate in compare 12")
		}

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
		for i := range list {
			if list[i]["id"] == target {
				page := service.GetPage(allList[len(allList)-2]["Next"])
				if list[i]["final"] != page.ID {
					t.Errorf("TestUpdatePageAndIterate")
				}
			}
		}
	}
}

func TestPageDeleteMiddle(t *testing.T) {
	jsonString := `[{"article":"11"}]`
	req, err := http.NewRequest("DELETE", "/header", bytes.NewReader([]byte(jsonString)))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	var target string = "b"
	vars := map[string]string{
		"header": target,
	}
	req = mux.SetURLVars(req, vars)
	handler := http.HandlerFunc(controller.DeletePage)
	handler.ServeHTTP(rr, req)
	bodyBytes, err := io.ReadAll(rr.Body)
	if err != nil {
		t.Errorf("TestPageDeleteMiddle")
		log.Fatal(err)
	}

	req, err = http.NewRequest("DELETE", "/header/GetHead/", nil)
	if err != nil {
		t.Errorf("TestPageDeleteMiddle")
		log.Fatal(err)
	}
	rr = httptest.NewRecorder()
	vars = map[string]string{
		"header": target,
	}
	req = mux.SetURLVars(req, vars)
	handler = http.HandlerFunc(controller.GetCertainHeader)
	handler.ServeHTTP(rr, req)
	bodyBytes, err = io.ReadAll(rr.Body)
	if err != nil {
		t.Errorf("TestPageDeleteMiddle")
		log.Fatal(err)
	}
	var next string = RemoveQuote(string(bodyBytes))
	if next == "" {
		t.Errorf("TestPageDeleteMiddle")
		log.Fatal("TestPageDeleteMiddle error in next")
	} else {
		var allList []map[string]string
		allList = iterate(next, t)
		if len(allList) != 2 {
			t.Errorf("TestPageDeleteMiddle in total")
		}
		if allList[0]["Article"] != "10" && allList[0]["Next"] != "" {
			t.Errorf("TestPageDeleteMiddle in compare 10")
		} else if allList[1]["Article"] != "12" && allList[1]["Next"] == "" {
			t.Errorf("TestPageDeleteMiddle in compare 12")
		}

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
		for i := range list {
			if list[i]["id"] == target {
				page := service.GetPage(allList[len(allList)-2]["Next"])
				if list[i]["final"] != page.ID {
					t.Errorf("TestPageDeleteMiddle")
				}
			}
		}
	}
}

func TestPageDeleteFirst(t *testing.T) {
	jsonString := `[{"article":"10"}]`
	req, err := http.NewRequest("DELETE", "/header", bytes.NewReader([]byte(jsonString)))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	var target string = "b"
	vars := map[string]string{
		"header": target,
	}
	req = mux.SetURLVars(req, vars)
	handler := http.HandlerFunc(controller.DeletePage)
	handler.ServeHTTP(rr, req)
	bodyBytes, err := io.ReadAll(rr.Body)
	if err != nil {
		t.Errorf("TestPageDeleteFirst")
		log.Fatal(err)
	}
	req, err = http.NewRequest("DELETE", "/header/GetHead/", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	vars = map[string]string{
		"header": target,
	}
	req = mux.SetURLVars(req, vars)
	handler = http.HandlerFunc(controller.GetCertainHeader)
	handler.ServeHTTP(rr, req)
	bodyBytes, err = io.ReadAll(rr.Body)
	if err != nil {
		t.Errorf("TestPageDeleteFirst")
		log.Fatal(err)
	}
	var next string = RemoveQuote(string(bodyBytes))
	if next == "" {
		t.Errorf("TestPageDeleteFirst")
		log.Fatal("TestPageDeleteFirst error")
	} else {
		var allList []map[string]string
		allList = iterate(next, t)
		if len(allList) != 1 {
			t.Errorf("TestPageDeleteFirst error")
		}
		if allList[0]["Article"] != "12" && allList[0]["Next"] == "" {
			t.Errorf("TestPageDeleteFirst error")
		}

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
		for i := range list {
			if list[i]["id"] == target {
				page := service.GetPage(next)
				if list[i]["final"] != page.ID {
					t.Errorf("TestPageDeleteFirst error")
				}
			}
		}
	}
}

func TestPageDeleteFinal(t *testing.T) {
	jsonString := `[{"article":"12"}]`
	req, err := http.NewRequest("DELETE", "/header", bytes.NewReader([]byte(jsonString)))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	var target string = "b"
	vars := map[string]string{
		"header": target,
	}
	req = mux.SetURLVars(req, vars)
	handler := http.HandlerFunc(controller.DeletePage)
	handler.ServeHTTP(rr, req)

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
	for i := range list {
		if list[i]["id"] == "b" {
			if list[i]["page_key"] != "" {
				t.Errorf("TestPageDeleteFinal error")
			}
			if list[i]["final"] != "" {
				t.Errorf("TestPageDeleteFinal error")
			}
		}
	}
}
