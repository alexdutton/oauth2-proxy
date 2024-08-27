package middleware

import (
	"fmt"
	"net/http"

	"github.com/justinas/alice"
	"github.com/oauth2-proxy/oauth2-proxy/v7/pkg/logger"
	"github.com/oauth2-proxy/oauth2-proxy/v7/providers/loader"
	"github.com/oauth2-proxy/oauth2-proxy/v7/providers/utils"
)

// middleware that loads the provider and stores it in the context
func NewProviderLoader(loader loader.Loader) alice.Constructor {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			ctx := req.Context()

			providerID := utils.ProviderIDFromContext(ctx)

			if providerID == "" {
				provider, err := loader.GetDefault(req.Context())
				if err == nil {
					ctx = utils.AppendProviderToContext(ctx, provider)
				}
			} else {
				provider, err := loader.Load(req.Context(), providerID)
				if err != nil {
					logger.Error(fmt.Sprintf("unable to load provider, id='%s': %s", providerID, err.Error()))
					rw.WriteHeader(http.StatusUnauthorized)
					return
				}

				ctx = utils.AppendProviderToContext(ctx, provider)
			}
			next.ServeHTTP(rw, req.WithContext(ctx))
		})
	}
}
