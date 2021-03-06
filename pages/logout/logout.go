package logout

import (
	"net/http"

	"github.com/vpoletaev11/fileHostingSite/errhand"
	"github.com/vpoletaev11/fileHostingSite/session"
)

// Page returns HandleFunc that removes user cookie and redirect to login page
func Page(dep session.Dependency) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		_, err = dep.Redis.Do("DEL", cookie.Value)
		if err != nil {
			errhand.InternalError(err, w)
			return
		}

		newCookie := &http.Cookie{
			Name:   "session_id",
			MaxAge: -1,
		}
		http.SetCookie(w, newCookie)
		http.Redirect(w, r, "/login", http.StatusFound)
	}
}
