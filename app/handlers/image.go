package handlers

import (
	"os"

	"abibby.com/mangadb/app/models"
	"gosalusa.com/request"
)

type ImageRequest struct {
	Image *models.Image `inject:"image_id"`
}

var Image = request.Handler(func(r *ImageRequest) (request.Responder, error) {
	f, err := os.Open(r.Image.Path)
	if err != nil {
		return nil, err
	}
	return request.NewResponse(f).AddHeader("Content-Type", "image/"+r.Image.Format), nil
})
