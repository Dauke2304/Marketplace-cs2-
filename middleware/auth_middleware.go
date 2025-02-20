package middleware

import (
	"Marketplace-cs2-/database"
	"Marketplace-cs2-/repositories"
	"Marketplace-cs2-/services"
	"fmt"
	"net/http"
)

func Authorize(r *http.Request) error {
	return services.ValidateAuthorization(r)
}

func AdminMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionCookie, err := r.Cookie("sessiontoken")
		if err != nil || sessionCookie.Value == "" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		fmt.Println(sessionCookie)

		database.InitDB()
		rep := repositories.NewUserRepository(database.Client.Database("cs2_skins_marketplace"))

		user, err := rep.GetUserBySessionToken(sessionCookie.Value)
		if err != nil || user == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		if !user.IsAdmin {
			http.Redirect(w, r, "/main", http.StatusSeeOther)
			return
		}

		next(w, r)
	}
}
