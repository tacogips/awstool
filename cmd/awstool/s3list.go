package awstool

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/go-xweb/log"
)

type S3Files struct {
	Bucket  string
	SumSize int64
	Count   int
	Objects []*s3.Object
}

func S3List(region, bucket, prefix string) (S3Files, error) {
	sess := session.New(&aws.Config{Region: aws.String(region)})

	s3client := s3.New(sess)
	files := S3Files{Bucket: bucket}

	appendToFiles := func(obj *s3.Object) {
		files.SumSize = files.SumSize + *obj.Size
		files.Count += 1
		files.Objects = append(files.Objects, obj)
	}

	println(prefix)
	input := &s3.ListObjectsInput{
		Bucket: aws.String(bucket),
		Prefix: aws.String(prefix),
	}

	output, err := s3client.ListObjects(input)
	if err != nil {
		log.Error(" ListObjects %s", err.Error())
		return files, err
	}
	for _, content := range output.Contents {
		appendToFiles(content)
	}

	for output.NextMarker != nil && len(*output.NextMarker) != 0 {
		input = &s3.ListObjectsInput{
			Marker: output.NextMarker,
			Bucket: aws.String(bucket),
			Prefix: aws.String(prefix),
		}

		output, err = s3client.ListObjects(input)
		if err != nil {
			log.Error(" ListObjects %s", err.Error())
			return files, err
		}
		for _, content := range output.Contents {
			appendToFiles(content)
		}
	}

	return files, nil
}
