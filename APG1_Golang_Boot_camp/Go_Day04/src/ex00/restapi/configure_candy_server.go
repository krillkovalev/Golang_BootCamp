// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"candy/restapi/operations"
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
)

//go:generate swagger generate server --target ../../ex00 --name CandyServer --spec ../swagger.yml --principal interface{}

func configureFlags(api *operations.CandyServerAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func buyCandy(params *operations.BuyCandyParams) middleware.Responder {
	candy_type := *params.Order.CandyType
	candy_count := *params.Order.CandyCount
	money := *params.Order.Money
	candies := map[string]int{"CE": 10, "AA": 15, "NT": 17, "DE": 21, "YR": 23}
	val, ok := candies[candy_type]
	if candy_count < 0 || !ok {
		str := "wrong candyType or candyCount"
		err := &operations.BuyCandyBadRequestBody{
			Error: str,
		}
		return operations.NewBuyCandyBadRequest().WithPayload(err)
	}
	sum := int64(val) * money
	if sum > money {
		str := fmt.Sprintf("You need %d more money!", sum-*params.Order.Money)
		err := &operations.BuyCandyPaymentRequiredBody{
			Error: str,
		}
		return operations.NewBuyCandyPaymentRequired().WithPayload(err)
	}

	if sum <= money {
		money = money - sum
		res := &operations.BuyCandyCreatedBody{
			Change: int64(money),
			Thanks: "Thank you!",
		}
		return operations.NewBuyCandyCreated().WithPayload(res)
	}

	return nil
}

func configureAPI(api *operations.CandyServerAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.UseSwaggerUI()
	// To continue using redoc as your UI, uncomment the following line
	// api.UseRedoc()

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	if api.BuyCandyHandler == nil {
		api.BuyCandyHandler = operations.BuyCandyHandlerFunc(func(params operations.BuyCandyParams) middleware.Responder {
			return buyCandy(&params)
		})
	}

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix".
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation.
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics.
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
