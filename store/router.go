package store

var controller = &Controller{Repository: Repository{}}

// Route struct
type Route struct {
	Name string
	Method string
	Pattern string
	HandlerFunc http.HandlerFunc
}

// List of routes
type Routes []Route

var routes = Routes {
	Route {
		"Authentication",
		"POST",
		"/get-token",
		controller.GetToken,
	},
	Route {
		"Index",
		"GET",
		"/",
		controller.Index,
	},
	Route {
		"AddProduct",
		"POST",
		"/AddProduct",
		AuthenticationMiddleware(controller.AddProduct),
	},
}

// DOCS: https://github.com/gorilla/mux#examples

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc

		router.
		Methods(route.Method).
		Path(route.Pattern).
		Name(route.Name).
		Handler(handler)
	}
	return router
}