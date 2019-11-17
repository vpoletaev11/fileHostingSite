package users

import (
	"database/sql"
	"html/template"
	"net/http"

	"github.com/vpoletaev11/fileHostingSite/errhand"
)

const selectUsers = "SELECT username, rating FROM users ORDER BY rating DESC LIMIT 15;"

// absolute path to template file
const absPathTemplate = "/home/perdator/go/src/github.com/vpoletaev11/fileHostingSite/pages/users/template/users.html"

// TemplateUsers contains fields with warning message and username for users page handler template
type TemplateUsers struct {
	Warning  template.HTML
	Username string
	UserList []UserInfo
}

type UserInfo struct {
	Username string
	Rating   int
}

// Page returns HandleFunc with access to MySQL database for index page
func Page(db *sql.DB, username string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// creating template for index page
		page, err := template.ParseFiles(absPathTemplate)
		if err != nil {
			errhand.InternalError("users", "Page", username, err, w)
			return
		}
		switch r.Method {
		case "GET":
			rows, err := db.Query(selectUsers)
			if err != nil {
				errhand.InternalError("users", "Page", username, err, w)
				return
			}
			defer rows.Close()

			usersInfo := []UserInfo{}
			for rows.Next() {
				ui := UserInfo{}

				err := rows.Scan(
					&ui.Username,
					&ui.Rating,
				)
				if err != nil {
					errhand.InternalError("users", "Page", username, err, w)
					return
				}
				usersInfo = append(usersInfo, ui)
			}

			err = page.Execute(w, TemplateUsers{Username: username, UserList: usersInfo})
			if err != nil {
				errhand.InternalError("users", "Page", username, err, w)
				return
			}
			return
		}
	}
}