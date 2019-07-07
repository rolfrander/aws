package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/glacier"
)

func getSession() (sess *session.Session, err error) {
	sess, err = session.NewSession(&aws.Config{
		Region: aws.String("eu-west-2")},
	)
	return
}

func upload(sess *session.Session, filename, vault, description string) error {
	svc := glacier.New(sess)
	svc.UploadArchive(&glacier.UploadArchiveInput{
		AccountId:          aws.String("-"),
		ArchiveDescription: aws.String(description),
		Body:               nil,
		Checksum:           nil,
		VaultName:          aws.String(vault),
	})
	return nil
}
