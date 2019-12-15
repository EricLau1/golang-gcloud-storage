package routes

import (
	"golang-gcloud-storage/app/controllers"
	"net/http"
)

type Route struct {
	Path    string
	Method  string
	Handler http.HandlerFunc
}

type AppRoutes interface {
	Routes() []*Route
}

type appRoutesImpl struct {
	appController controllers.AppController
}

func NewAppRoutes(appController controllers.AppController) *appRoutesImpl {
	return &appRoutesImpl{appController}
}

func (a *appRoutesImpl) Routes() []*Route {
	return []*Route{
		&Route{
			Path:    "/",
			Method:  http.MethodGet,
			Handler: a.appController.Home,
		},
		&Route{
			Path:    "/images",
			Method:  http.MethodGet,
			Handler: a.appController.GetImages,
		},
		&Route{
			Path:    "/videos",
			Method:  http.MethodGet,
			Handler: a.appController.GetVideos,
		},
		&Route{
			Path:    "/uploads",
			Method:  http.MethodGet,
			Handler: a.appController.GetUpload,
		},
		&Route{
			Path:    "/uploads",
			Method:  http.MethodPost,
			Handler: a.appController.PostUpload,
		},
	}
}
