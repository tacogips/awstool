package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

var sess = session.New(&aws.Config{Region: aws.String("ap-northeast-1")})

func main() {

}
