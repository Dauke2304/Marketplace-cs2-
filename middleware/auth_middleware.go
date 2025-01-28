package middleware

import (
	"Marketplace-cs2-/services"
	"net/http"
)

func Authorize(r *http.Request) error {
	return services.ValidateAuthorization(r)
}
