package repositories

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/configuration"
	"github.com/swiftwaterlabs/engineering-intelligence-services/internal/pkg/models"
	"log"
	"strings"
)

type S3EngineeringObjectRepository struct {
	bucketName   string
	awsRegion    string
	fileUploader *s3manager.Uploader
}

func (r *S3EngineeringObjectRepository) init(config configuration.ConfigurationService) {
	r.bucketName = config.GetValue("engineering_intelligence_prd_blob_store")
	r.awsRegion = config.GetValue("aws_region")
	r.fileUploader = r.initS3Uploader()
}

func (r *S3EngineeringObjectRepository) Save(item *models.EngineeringObject) error {

	input := &s3manager.UploadInput{
		Bucket:      aws.String(r.bucketName),
		Key:         aws.String(r.getFileKey(item)),
		Body:        strings.NewReader(item.Data),
		ContentType: aws.String("application/json"),
	}
	_, err := r.fileUploader.UploadWithContext(context.Background(), input)
	if err != nil {
		return err
	}

	return nil
}

func (r *S3EngineeringObjectRepository) getFileKey(item *models.EngineeringObject) string {
	path := r.resolveItemPath(item)
	identifier := r.resolveItemId(item)

	return fmt.Sprintf("%v/%v", path, identifier)
}

func (r *S3EngineeringObjectRepository) resolveItemPath(item *models.EngineeringObject) string {
	if item == nil || item.Type == "" {
		return "unknown"
	}

	return strings.ToLower(item.Type)
}

func (r *S3EngineeringObjectRepository) resolveItemId(item *models.EngineeringObject) string {
	if item == nil || item.Id == "" {
		return uuid.New().String()
	}

	return strings.ToLower(item.Id)
}

func (r *S3EngineeringObjectRepository) initS3Uploader() *s3manager.Uploader {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(r.awsRegion)},
	)
	if err != nil {
		log.Fatalf("failed to create AWS session, %v", err)
	}

	uploader := s3manager.NewUploader(sess)
	return uploader
}

func (r *S3EngineeringObjectRepository) Destroy() {

}
