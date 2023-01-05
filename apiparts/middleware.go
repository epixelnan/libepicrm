// This file is part of Epixel CRM Software

package epicrm_apiparts

import (
	"context"
	"net/http"
	"os"
	
	"github.com/google/uuid"

	"github.com/jackc/pgx/v4/pgxpool"
	
	"github.com/go-chi/chi/v5"
	chi_middleware "github.com/go-chi/chi/v5/middleware"
	middleware "github.com/deepmap/oapi-codegen/pkg/chi-middleware"
	"github.com/getkin/kin-openapi/openapi3filter"
)

func GetMiddlewareOptions() *middleware.Options {
	return &middleware.Options{
		Options: openapi3filter.Options {
			AuthenticationFunc: newAuthValidator(),
		},
	}
}

func validateAuth(cntxt context.Context, input *openapi3filter.AuthenticationInput) error {
	// TODO Make sure nginx has performed the auth validation
	
	// TODO check if the app has access to this module (and the orgnaization? but isn't that for the endpoints to check?); get the app info from the auth token, not the incoming request
	
	return nil
}

func newAuthValidator() openapi3filter.AuthenticationFunc {
	return func(cntxt context.Context, input *openapi3filter.AuthenticationInput) error {
		return validateAuth(cntxt, input)
	}
}

func AddCommonMiddlewares(r chi.Router, serviceName string) {
	r.Use(NewAttachServiceName(serviceName))
	r.Use(chi_middleware.Logger)
	r.Use(ReadOnlyMode)
}

func ExtractBearerToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		btokstr := ""
	
		btoken := extractBearerTokenFromHeader(r.Header.Get("Authorization"))
		if btoken != nil {
			btokstr = string(btoken)
		}

		newctx := context.WithValue(r.Context(), "bearerToken", btokstr)
		                            
		next.ServeHTTP(w, r.WithContext(newctx))
	})
}

// Use only when you are sure the UidFromAPIGatewayMust() middleware
// has been used. Or the server will panic for every unauthorized request.
func GetUidOrPanic(r *http.Request) string {
	buid := r.Context().Value("uidFromGateway").(string)
	_, err := uuid.Parse(buid)
	if err != nil {
		panic("uidFromTokstore not set")
	}
	
	return buid
}

// Use only when you are sure the UidFromAPIGatewayMust() middleware
// has been used. Or the server will panic for every unauthorized request.
func GetUidUuidOrPanic(r *http.Request) uuid.UUID {
	buid := r.Context().Value("uidFromGateway").(string)
	uuid, err := uuid.Parse(buid)
	if err != nil {
		panic("uidFromTokstore not set")
	}
	
	return uuid
}

func ReadOnlyMode(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if os.Getenv("EPICRM_READONLY") != "" {
			if r.Method != "GET" && r.Method != "OPTIONS" {
				Send500ReadOnlyMode(w)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

func UidFromAPIGatewayMust(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO HMAC check

		uidhdr := r.Header.Get("X-EPICRM-UID")
		_, err := uuid.Parse(uidhdr)
		if err != nil {
			Send401Bearer(w)
			return
		}

		newctx := context.WithValue(r.Context(), "uidFromGateway", string(uidhdr))

		next.ServeHTTP(w, r.WithContext(newctx))
	})
}

func NewUserMustBeA(role string, orgid uuid.UUID, authdbcon *pgxpool.Pool) func(next http.Handler) http.Handler {
	return func (next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			uidhdr := r.Header.Get("X-EPICRM-UID")
			uid, err := uuid.Parse(uidhdr)
			if err != nil {
				Send401Bearer(w)
				return
			}

			hasRole, err :=
				UserIsAnyOf(
					authdbcon, uid.String(), orgid.String(), []string{role})
			if err != nil {
				LogAndSendError(w,
					r.Context().Value("serviceName").(string),
					"NewUserMustBeA",
					err)

				return
			}

			if !hasRole {
				Send403(w)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func NewAttachServiceName(serviceName string) func(next http.Handler) http.Handler {
	return func (next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			newctx := context.WithValue(r.Context(), "serviceName", serviceName)
			next.ServeHTTP(w, r.WithContext(newctx))
		})
	}
}
