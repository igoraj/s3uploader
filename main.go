package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func main() {
	// display help function
	flag.Usage = func() {
		fmt.Printf("Usage: s3upload [options] <filename>\n\n")
		fmt.Printf("  Set AWS_ACCESS_KEY and AWS_SECRET_KEY environmental variables.\n")
		fmt.Printf("  before invoking this command.\n\n")
		flag.PrintDefaults()
	}

	// parse flags
	var (
		region   = flag.String("region", "eu-geo", "region of your request")
		endpoint = flag.String("endpoint", "s3.eu-geo.objectstorage.softlayer.net", "authentication endpoint")
		bucket   = flag.String("bucket", "", "bucket name")
	)
	flag.Parse()

	// parse args
	args := flag.Args()
	if len(args) != 1 {
		flag.Usage()
		os.Exit(1)
	}
	filename := args[0]

	// open filename
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("failed to open file", filename, err)
		os.Exit(1)
	}
	defer file.Close()

	// create session to object storage
	sess := session.Must(session.NewSession((&aws.Config{
		Region:   aws.String(*region),
		Endpoint: aws.String(*endpoint),
	})))
	svc := s3manager.NewUploader(sess)

	// upload file
	fmt.Println("uploading file...")
	result, err := svc.Upload(&s3manager.UploadInput{
		Bucket: bucket,
		Key:    aws.String(filepath.Base(filename)),
		Body:   file,
	})
	if err != nil {
		fmt.Println("error", err)
		os.Exit(1)
	}

	fmt.Printf("finished uploading %s to %s\n", filename, result.Location)
}
