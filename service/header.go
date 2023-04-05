package service

type Header struct {
	ID       string `json:"id"`
	Page_key string `json:"page_Key"`
	Final    string `json:"final"`
}

func GetCertainHeader(header string) string {
	db := SetupDB()
	sqlStatement := `SELECT * FROM "header" WHERE "header" = $1`
	res, err := db.Query(sqlStatement, header)
	if err != nil {
		var response string = "get head error"
		return response
	} else {
		var headerContainer []Header
		for res.Next() {
			var id string
			var page_Key string
			var final string

			err = res.Scan(&id, &page_Key, &final)
			headerContainer = append(headerContainer, Header{ID: RemoveSpace(id), Page_key: RemoveSpace(page_Key), Final: RemoveSpace(final)})
		}
		return headerContainer[0].Page_key
	}
}

func GetHeader(header string) Header {
	db := SetupDB()
	sqlStatement := `SELECT * FROM "header" WHERE "header" = $1`
	res, err := db.Query(sqlStatement, header)
	var headerContainer []Header
	if err != nil {
		var response string = "get head error"
		headerContainer[0].ID = response
		return headerContainer[0]
	} else {
		for res.Next() {
			var id string
			var page_Key string
			var final string

			err = res.Scan(&id, &page_Key, &final)
			headerContainer = append(headerContainer, Header{ID: RemoveSpace(id), Page_key: RemoveSpace(page_Key), Final: RemoveSpace(final)})
		}
		return headerContainer[0]
	}
}

func FixHeaderFinal(idOld string, idNew string) bool {
	db := SetupDB()
	sqlStatement := `UPDATE "header" SET "final" = $2 WHERE "header" = $1`
	_, err := db.Exec(sqlStatement, idOld, idNew)
	if err != nil {
		return false
	} else {
		return true
	}
}

func FixHeaderPageKey(header string, idNew string) bool {
	db := SetupDB()
	sqlStatement := `UPDATE "header" SET "page_key" = $2 WHERE "header" = $1`
	_, err := db.Exec(sqlStatement, header, idNew)
	if err != nil {
		return false
	} else {
		return true
	}
}
