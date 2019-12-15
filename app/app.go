package app

import (
	"context"
	"flag"
	"fmt"
	"golang-gcloud-storage/app/controllers"
	"golang-gcloud-storage/app/render"
	"golang-gcloud-storage/app/routes"
	"golang-gcloud-storage/bucket"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var (
	port = flag.Int("p", 8080, "set app port")
)

func Run() {
	flag.Parse()

	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	render.Load("app/pages/*.html")

	ctx := context.Background()

	b := bucket.NewBucket(ctx)

	fmt.Println("App running on port", *port)

	appController := controllers.NewAppController(b)
	appRoutes := routes.NewAppRoutes(appController)

	router := mux.NewRouter()
	bind(router, appRoutes)

	p := fmt.Sprintf(":%d", *port)

	log.Fatal(http.ListenAndServe(p, router))

}

func bind(router *mux.Router, appRoutes routes.AppRoutes) {
	for _, route := range appRoutes.Routes() {
		router.HandleFunc(route.Path, route.Handler).Methods(route.Method)
	}
}
