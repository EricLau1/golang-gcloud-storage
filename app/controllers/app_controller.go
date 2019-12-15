package controllers

import (
	"context"
	"fmt"
	"golang-gcloud-storage/app/render"
	"golang-gcloud-storage/bucket"
	"golang-gcloud-storage/models"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type AppController interface {
	Home(http.ResponseWriter, *http.Request)
	GetImages(http.ResponseWriter, *http.Request)
	GetVideos(http.ResponseWriter, *http.Request)
	GetUpload(http.ResponseWriter, *http.Request)
	PostUpload(http.ResponseWriter, *http.Request)
}

type appControllerImpl struct {
	bucketHandler bucket.Bucket
}

func NewAppController(bucketHandler bucket.Bucket) *appControllerImpl {
	return &appControllerImpl{bucketHandler}
}

func (a *appControllerImpl) Home(w http.ResponseWriter, r *http.Request) {
	render.Page(w, "home", nil)
}

func (a *appControllerImpl) GetImages(w http.ResponseWriter, r *http.Request) {
	images, err := a.bucketHandler.GetImages(context.Background())
	if err != nil {
		w.Write([]byte(err.Error()))

		return
	}

	render.Page(w, "images", images)
}

func (a *appControllerImpl) GetVideos(w http.ResponseWriter, r *http.Request) {
	videos, err := a.bucketHandler.GetVideos(context.Background())
	if err != nil {
		w.Write([]byte(err.Error()))

		return
	}

	render.Page(w, "videos", videos)
}

func (a *appControllerImpl) GetUpload(w http.ResponseWriter, r *http.Request) {
	uploads, err := a.bucketHandler.GetUploads(context.Background())
	if err != nil {
		log.Println(err)

		render.Page(w, "uploads", nil)
		return
	}

	render.Page(w, "uploads", uploads)
}

func (a *appControllerImpl) PostUpload(w http.ResponseWriter, r *http.Request) {
	f, fh, err := r.FormFile("file")
	if err != nil {
		w.Write([]byte(err.Error()))

		return
	}

	defer f.Close()

	uniqueName := fmt.Sprintf("%s-%s", uuid.New().String(), fh.Filename)

	uploadMidia := &models.Midia{
		Name: uniqueName,
		Type: fh.Header.Get("Content-Type"),
		Size: fh.Size,
		Link: "#",
		File: f,
	}

	ctx := context.Background()

	publicLink, err := a.bucketHandler.UploadMidia(ctx, uploadMidia)
	if err != nil {
		w.Write([]byte(err.Error()))

		return
	}

	http.Redirect(w, r, publicLink, 302)
}
