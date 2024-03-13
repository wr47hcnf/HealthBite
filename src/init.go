package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"log"
	"net/http"
	"time"
)

var Svc *s3.S3

func init() {
	fmt.Printf("HealthBite backend server\n(C) 2024 Patrick Covaci a.k.a Ty3r0X\n%s\n", time.Now())

	err := dbConnect()

	if err != nil {
		log.Fatal("Failed to connect to database")
	}

	err = dbInit()

	if err != nil {
		log.Fatal("Failed to initialize db")
	}

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("eu-north-1"),
		Credentials: credentials.NewStaticCredentials(aws_access, aws_secret, ""),
	})
	if err != nil {
		log.Fatal(err)
	}

	Svc = s3.New(sess)

	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
}
