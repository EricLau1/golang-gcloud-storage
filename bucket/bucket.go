package bucket

import (
	"context"
	"fmt"
	"golang-gcloud-storage/models"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/iterator"

	"cloud.google.com/go/storage"
)

type Bucket interface {
	GetImages(context.Context) ([]*models.Midia, error)
	GetVideos(context.Context) ([]*models.Midia, error)
	UploadMidia(context.Context, *models.Midia) (string, error)
	GetUploads(context.Context) ([]*models.Midia, error)
}

type bucketImpl struct {
	handler *storage.BucketHandle
}

const (
	images_prefix  = "images/"
	videos_prefix  = "videos/"
	uploads_prefix = "uploads/"
)

func NewBucket(ctx context.Context) *bucketImpl {
	cli, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}

	return &bucketImpl{cli.Bucket(getBucketName())}
}

func (b *bucketImpl) setReadOnlyAccessPublic(ctx context.Context, objectName string) {
	acl := b.handler.Object(objectName).ACL()
	err := acl.Set(ctx, storage.AllUsers, storage.RoleReader)
	if err != nil {
		log.Println(err)
	}
}

func (b *bucketImpl) getMidias(ctx context.Context, prefix string) ([]*models.Midia, error) {

	q := &storage.Query{
		Prefix:    prefix,
		Delimiter: "/",
	}

	it := b.handler.Objects(ctx, q)

	midias := []*models.Midia{}

	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			return nil, err
		}

		b.setReadOnlyAccessPublic(ctx, attrs.Name)

		midias = append(midias, buildMidia(attrs))
	}

	return midias, nil
}

func buildMidia(attrs *storage.ObjectAttrs) *models.Midia {
	m := &models.Midia{
		Name: attrs.Name,
		Type: attrs.ContentType,
		Link: "#",
		Size: attrs.Size,
	}

	link, err := makePublicLink(m.Name)
	if err != nil {
		log.Println(err)

		return m
	}

	m.Link = link
	return m
}

func makePublicLink(objectName string) (string, error) {
	key, err := ioutil.ReadFile(getCredentials())
	if err != nil {
		return "", err
	}

	cfg, err := google.JWTConfigFromJSON(key)
	if err != nil {
		return "", err
	}

	urlOptions := &storage.SignedURLOptions{
		GoogleAccessID: cfg.Email,
		PrivateKey:     cfg.PrivateKey,
		Expires:        time.Now().Add(time.Minute * 10),
		Method:         http.MethodGet,
	}

	return storage.SignedURL(getBucketName(), objectName, urlOptions)
}

func (b *bucketImpl) GetImages(ctx context.Context) ([]*models.Midia, error) {
	images, err := b.getMidias(ctx, images_prefix)
	if err != nil {
		return nil, err
	}

	return filterMidias(images, "image"), nil
}

func (b *bucketImpl) GetVideos(ctx context.Context) ([]*models.Midia, error) {
	videos, err := b.getMidias(ctx, videos_prefix)
	if err != nil {
		return nil, err
	}

	return filterMidias(videos, "video"), nil
}

func (b *bucketImpl) GetUploads(ctx context.Context) ([]*models.Midia, error) {
	uploads, err := b.getMidias(ctx, uploads_prefix)
	if err != nil {
		return nil, err
	}

	filtered := filterMidias(uploads, "image")
	filtered = append(filtered, filterMidias(uploads, "video")...)

	return filtered, nil
}

func filterMidias(midias []*models.Midia, contentType string) []*models.Midia {
	filtered := []*models.Midia{}
	for _, m := range midias {
		if strings.Contains(m.Type, contentType) {
			filtered = append(filtered, m)
		}
	}

	return filtered
}

func (b *bucketImpl) UploadMidia(ctx context.Context, midia *models.Midia) (string, error) {
	writer := b.handler.Object(fmt.Sprintf("uploads/%s", midia.Name)).NewWriter(ctx)
	writer.ContentType = midia.Type
	writer.CacheControl = "public, max-age=86400"
	writer.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}

	if _, err := io.Copy(writer, midia.File); err != nil {
		return "#", err
	}
	if err := writer.Close(); err != nil {
		return "#", err
	}

	const publicLink = "https://storage.googleapis.com/%s/uploads/%s"
	return fmt.Sprintf(publicLink, getBucketName(), midia.Name), nil
}
