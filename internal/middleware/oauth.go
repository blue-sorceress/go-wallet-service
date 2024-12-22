package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"

	"go-wallet-service/internal/object"
	s "go-wallet-service/internal/service"
	"go-wallet-service/utils"
)

func OAuth(userService *s.UserService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authorizationHeader := r.Header.Get("Authorization")

			if !strings.Contains(authorizationHeader, "Bearer") {
				log.Error().Str("error", "authorization_failed").Str("error_description", "OAuth authorization required").Send()
				utils.HttpErrorResponse(w, "invalid_authorization", "OAuth Bearer authorization required", http.StatusUnauthorized)
				return
			}

			authorization := strings.Split(authorizationHeader, " ")
			if len(authorization) != 2 {
				utils.HttpErrorResponse(w, "invalid_authorization", "OAuth Bearer authorization invalid format", http.StatusUnauthorized)
				return
			}

			token := authorization[1]
			userId, err := userService.GetUserByToken(token)
			if err != nil {
				log.Error().Str("error", "token_invalid").Str("error_description", fmt.Sprintf("Token [%s] invalid for user due to %s", token, err.Error())).Send()
				utils.HttpErrorResponse(w, "authorization_failed", "The token provided is invalid or user does not exists.", http.StatusUnauthorized)
				return
			}

			body, err := io.ReadAll(r.Body)
			if err != nil {
				log.Error().Str("error", "io_read_error").Str("error_description", "Unable to read body").Send()
				utils.HttpErrorResponse(w, "invalid_request", "Unable to read body", http.StatusForbidden)
				return
			}
			r.Body.Close()

			ctx := context.WithValue(r.Context(), "requestParams", string(body))
			r = r.WithContext(ctx)

			r.Body = io.NopCloser(bytes.NewReader(body))
			var OAuthParams object.OAuthParams
			OAuthParamsErr := json.NewDecoder(r.Body).Decode(&OAuthParams)
			if OAuthParamsErr == nil {
				if userId != OAuthParams.UserId {
					log.Error().Str("error", "invalid_authorization").Str("error_description", fmt.Sprintf("User [%d] is unauthorized", userId)).Send()
					utils.HttpErrorResponse(w, "invalid_authorization", "The current user is unauthorized", http.StatusForbidden)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
