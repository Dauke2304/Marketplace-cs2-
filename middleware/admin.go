package middleware

import (
	"Marketplace-cs2-/services"
	"net/http"
)

func AdminOnly(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Validate authorization first
		if err := services.ValidateAuthorization(r); err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Get current user
		user := services.GetCurrentUser(r)
		if user == nil {
			http.Error(w, "User not found", http.StatusUnauthorized)
			return
		}

		// Check admin status
		if !user.IsAdmin {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		// Proceed to admin panel
		next(w, r)
	}
}
