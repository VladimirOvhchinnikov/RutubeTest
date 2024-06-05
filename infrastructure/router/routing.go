package router

import (
	"net/http"
	"rutube/controller"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type GoChiRouter struct {
	Logger  *zap.Logger
	Router  chi.Router
	Handler controller.HandlersInterface
}

func NewGoChiRouting(logger *zap.Logger, handler controller.HandlersInterface) *GoChiRouter {

	router := chi.NewRouter()
	router.Use(loggingMiddleware(logger))
	router.Post("/telegram-webhook", handler.CommandHandler)

	return &GoChiRouter{
		Logger: logger,
		Router: router,
	}
}

func (gc *GoChiRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	gc.Router.ServeHTTP(w, r)
}

func loggingMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Info("Incoming request", zap.String("method", r.Method), zap.String("url", r.URL.String()))
			next.ServeHTTP(w, r)
		})
	}
}
