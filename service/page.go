package service

type Page struct {
	ID       string `json:"id"`
	Article  string `json:"article"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
}

func AddPage(page Page) bool {
	db := SetupDB()
	sqlStatement := `INSERT INTO "page" ("ID","article","next","previous") VALUES ($1,$2,$3,$4)`
	_, err := db.Exec(sqlStatement, page.ID, page.Article, page.Next, page.Previous)
	if err != nil {
		return false
	} else {
		return true
	}
}

func GetPage(id string) Page {
	db := SetupDB()
	sqlStatement := `SELECT * FROM "page" where "ID"=$1`
	res, err := db.Query(sqlStatement, id)
	var headerContainer []Page
	if err != nil {
		var response string = "page error"
		headerContainer[0].ID = response
		return headerContainer[0]
	} else {
		for res.Next() {
			var id string
			var article string
			var next string
			var previous string
			var time string

			err = res.Scan(&id, &article, &next, &time, &previous)
			headerContainer = append(headerContainer, Page{ID: RemoveSpace(id), Article: RemoveSpace(article), Next: RemoveSpace(next), Previous: RemoveSpace(previous)})
		}
		return headerContainer[0]
	}
}

func ModifyNext(id string, newNext string) bool {
	db := SetupDB()
	sqlStatement := `UPDATE "page" SET "next" = $2 WHERE "ID" = $1`
	_, err := db.Exec(sqlStatement, id, newNext)
	if err != nil {
		return false
	} else {
		return true
	}
}

func ModifyPrevious(id string, newPrevious string) bool {
	db := SetupDB()
	sqlStatement := `UPDATE "page" SET "previous" = $2 WHERE "ID" = $1`
	_, err := db.Exec(sqlStatement, id, newPrevious)
	if err != nil {
		return false
	} else {
		return true
	}
}

func ModifyID(id string, newID string, newArticle string) bool {
	db := SetupDB()
	sqlStatement := `UPDATE "page" SET "ID" = $2, "article" = $3 WHERE "ID" = $1`
	_, err := db.Exec(sqlStatement, id, newID, newArticle)
	if err != nil {
		return false
	} else {
		return true
	}
}

func DeletePage(id string) bool {
	db := SetupDB()
	sqlStatement := `DELETE FROM "page" WHERE "ID" = $1`
	_, err := db.Exec(sqlStatement, id)
	if err != nil {
		return false
	} else {
		return true
	}
}
