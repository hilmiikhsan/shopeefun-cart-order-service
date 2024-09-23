package routes

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/hilmiikhsan/shopeefun-cart-order-service/config"
	"github.com/hilmiikhsan/shopeefun-cart-order-service/handlers/cart"
	"github.com/hilmiikhsan/shopeefun-cart-order-service/handlers/order"
	"github.com/hilmiikhsan/shopeefun-cart-order-service/utils/middleware"
	"github.com/spf13/viper"
)

type Routes struct {
	Router *http.ServeMux
	Cart   *cart.Handler
	Order  *order.Handler
}

func URLRewriter(baseURLPath string, next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.TrimPrefix(r.URL.Path, baseURLPath)

		next.ServeHTTP(w, r)
	}
}

func (r *Routes) SetupBaseURL() {
	baseURL := viper.GetString("BASE_URL_PATH")

	if baseURL != "" && baseURL != "/" {
		r.Router.HandleFunc(baseURL+"/", URLRewriter(baseURL, r.Router))
	}
}

func (r *Routes) cartRoutes() {
	r.Router.HandleFunc("GET /cart/{user_id}", middleware.ApplyMiddleware(r.Cart.GetCartByUserID, middleware.EnabledCors, middleware.LoggerMiddleware()))
	r.Router.HandleFunc("POST /cart/add", middleware.ApplyMiddleware(r.Cart.AddCart, middleware.EnabledCors, middleware.LoggerMiddleware()))
	r.Router.HandleFunc("PUT /cart/update/{user_id}", middleware.ApplyMiddleware(r.Cart.UpdateCart, middleware.EnabledCors, middleware.LoggerMiddleware()))
	r.Router.HandleFunc("DELETE /cart/delete/{user_id}", middleware.ApplyMiddleware(r.Cart.DeleteCart, middleware.EnabledCors, middleware.LoggerMiddleware()))
}

func (r *Routes) orderRoutes() {
	r.Router.HandleFunc("POST /order/create", middleware.ApplyMiddleware(r.Order.CreateOrder, middleware.EnabledCors, middleware.LoggerMiddleware()))
	r.Router.HandleFunc("POST /order/callback", middleware.ApplyMiddleware(r.Order.UpdateOrder, middleware.EnabledCors, middleware.LoggerMiddleware()))
}

func (r *Routes) SetupRouter() {
	r.Router = http.NewServeMux()
	r.SetupBaseURL()
	r.cartRoutes()
	r.orderRoutes()
}

func (r *Routes) Run(port string) {
	r.SetupRouter()

	log.Printf("[Running-Success] clients on localhost on port :%s", port)

	srv := &http.Server{
		Handler:      r.Router,
		Addr:         "localhost:" + port,
		WriteTimeout: config.WriteTimeout() * time.Second,
		ReadTimeout:  config.ReadTimeout() * time.Second,
	}

	log.Panic(srv.ListenAndServe())
}
