package utils

import (
	"chat-app/internal/shared/config"
	"chat-app/internal/shared/logger"
	"context"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type ImageUpload struct {
	cloud *cloudinary.Cloudinary
	log   *logger.Logger
}

func (i *ImageUpload) UploadToCloudinary(ctx context.Context, imageData *multipart.FileHeader) (string, error) {

	folderName := "profile-image"

	file, err := imageData.Open()

	if err != nil {
		return "", err
	}
	defer file.Close()

	resp, err := i.cloud.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder: folderName,
	})

	if err != nil {
		i.log.Error("Failed to upload image")
	}

	i.log.Info("image uploaded %v",resp.SecureURL)

	return resp.SecureURL, nil

}

func NewImageUpload(cfg config.CloudinaryConfig, log *logger.Logger) (*ImageUpload, error) {
	cloud, err := cloudinary.NewFromURL(cfg.URL)

	if err != nil {
		log.Error("Cloud connection failed %v",err)
			return nil, err
		
	}
	return &ImageUpload{cloud: cloud, log: log}, nil
}
