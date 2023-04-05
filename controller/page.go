package controller

import (
	"Dcard/service"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func GetPage(w http.ResponseWriter, r *http.Request) {
	respBody, _ := ioutil.ReadAll(r.Body)
	var body []map[string]string
	json.Unmarshal(respBody, &body)
	var id string = body[0]["id"]
	page := service.GetPage(id)
	m := make(map[string]string)
	m["Article"] = page.Article
	m["Next"] = page.Next
	json.NewEncoder(w).Encode(m)
}

func UpdatePage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	header := vars["header"]
	first := service.GetHeader(header) //return Header
	respBody, _ := ioutil.ReadAll(r.Body)
	var body []map[string]string
	json.Unmarshal(respBody, &body)
	var originArticle string = body[0]["originArticle"]
	var newArticle string = body[0]["newArticle"]
	var originArticleID string = GetMD5Hash(originArticle + header)
	var newArticleID string = GetMD5Hash(newArticle + header)
	page := service.GetPage(originArticleID)
	if first.Page_key == originArticleID && first.Final == originArticleID {
		//fmt.Println("all equal")
		if !service.FixHeaderFinal(header, newArticleID) {
			json.NewEncoder(w).Encode("error in FixHeaderFinal")
		}
		if !service.FixHeaderPageKey(header, newArticleID) {
			json.NewEncoder(w).Encode("error in FixHeaderPageKey")
		}
	} else if first.Page_key == originArticleID {
		//fmt.Println("page key equal")
		if !service.ModifyPrevious(page.Next, newArticleID) {
			json.NewEncoder(w).Encode("error in ModifyPrevious")
		}
		if !service.FixHeaderPageKey(header, newArticleID) {
			json.NewEncoder(w).Encode("error in FixHeaderPageKey")
		}
	} else if first.Final == originArticleID {
		//fmt.Println("final equal")
		if !service.ModifyNext(page.Previous, newArticleID) {
			json.NewEncoder(w).Encode("error in ModifyNext")
		}
		if !service.FixHeaderFinal(header, newArticleID) {
			json.NewEncoder(w).Encode("error in FixHeaderFinal")
		}
	} else {
		//fmt.Println("none equal" + page.Previous + " " + page.Next)
		if !service.ModifyNext(page.Previous, newArticleID) {
			json.NewEncoder(w).Encode("error in ModifyNext")
		}
		if !service.ModifyPrevious(page.Next, newArticleID) {
			json.NewEncoder(w).Encode("error in ModifyPrevious")
		}
	}
	if !service.ModifyID(originArticleID, newArticleID, newArticle) {
		json.NewEncoder(w).Encode("error in ModifyID")
	}

}

func DeletePage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	header := vars["header"]
	first := service.GetHeader(header) //return Header
	respBody, _ := ioutil.ReadAll(r.Body)
	var body []map[string]string
	json.Unmarshal(respBody, &body)
	var article string = body[0]["article"]
	var articleID string = GetMD5Hash(article + header)
	page := service.GetPage(articleID)
	if first.Page_key == articleID && first.Final == articleID {
		//fmt.Println("all equal")
		if !service.FixHeaderFinal(header, "") {
			json.NewEncoder(w).Encode("error in FixHeaderFinal")
		}
		if !service.FixHeaderPageKey(header, "") {
			json.NewEncoder(w).Encode("error in FixHeaderPageKey")
		}
	} else if first.Page_key == articleID {
		//fmt.Println("page key equal")
		if !service.ModifyPrevious(page.Next, "") {
			json.NewEncoder(w).Encode("error in ModifyPrevious")
		}
		if !service.FixHeaderPageKey(header, page.Next) {
			json.NewEncoder(w).Encode("error in FixHeaderPageKey")
		}
	} else if first.Final == articleID {
		//fmt.Println("final equal")
		if !service.ModifyNext(page.Previous, "") {
			json.NewEncoder(w).Encode("error in ModifyNext")
		}
		if !service.FixHeaderFinal(header, page.Previous) {
			json.NewEncoder(w).Encode("error in FixHeaderFinal")
		}
	} else {
		//fmt.Println("none equal" + page.Previous + " " + page.Next)
		if !service.ModifyNext(page.Previous, page.Next) {
			json.NewEncoder(w).Encode("error in ModifyNext")
		}
		if !service.ModifyPrevious(page.Next, page.Previous) {
			json.NewEncoder(w).Encode("error in ModifyPrevious")
		}
	}
	if !service.DeletePage(articleID) {
		json.NewEncoder(w).Encode("error in DeletePage")
	}

}

func AddPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	header := vars["header"]
	respBody, _ := ioutil.ReadAll(r.Body)
	var body []map[string]string
	json.Unmarshal(respBody, &body)
	first := service.GetHeader(header) //return Header
	if first.Page_key != "" {
		var page service.Page
		var list []service.Page
		for i := 0; i < len(body); i++ {
			id := GetMD5Hash(body[i]["article"] + header)
			if i == 0 {
				page.Previous = first.Final
			} else {
				list[i-1].Next = id
				page.Previous = list[i-1].ID
			}
			page.Next = ""
			page.Article = body[i]["article"]
			page.ID = id
			list = append(list, page)
		}
		for j := range list {
			if !service.AddPage(list[j]) {
				json.NewEncoder(w).Encode("error in AddPage")
			}
		}
		if !service.FixHeaderFinal(first.ID, list[len(list)-1].ID) {
			json.NewEncoder(w).Encode("error in FixHeaderFinal")
		}
	} else {
		var page service.Page
		var list []service.Page
		for i := 0; i < len(body); i++ {
			id := GetMD5Hash(body[i]["article"] + header)
			if i == 0 {
				page.Previous = ""
			} else {
				list[i-1].Next = id
				page.Previous = list[i-1].ID
			}
			page.Next = ""
			page.Article = body[i]["article"]
			page.ID = id
			list = append(list, page)
		}
		for h := range list {
			if !service.AddPage(list[h]) {
				json.NewEncoder(w).Encode("error in AddPage")
			}
		}
		if !service.FixHeaderPageKey(first.ID, list[0].ID) {
			json.NewEncoder(w).Encode("error in FixHeaderPageKey")
		}
		if !service.FixHeaderFinal(first.ID, list[len(list)-1].ID) {
			json.NewEncoder(w).Encode("error in FixHeaderFinal")
		}
	}

}
