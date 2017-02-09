package awstool

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/go-xweb/log"
)

func S3DownloadPrefix(region, bucket, prefix string, destDir string) error {
	files, err := S3List(region, bucket, prefix)
	if err != nil {
		return err
	}

	sess := session.New(&aws.Config{Region: aws.String(region)})
	s3client := s3.New(sess)

	bulkDlCount := 10 //TODO to be changable

	fmt.Printf("download %d files. Total %d bytes \n", len(files.Objects), files.SumSize)

	var dlTargts []*s3.Object
	var rest []*s3.Object

	if len(files.Objects) < bulkDlCount {
		dlTargts, rest = files.Objects, nil
	} else {
		dlTargts, rest = files.Objects[:bulkDlCount], files.Objects[bulkDlCount:]
	}

	for len(dlTargts) > 0 {
		fmt.Printf("Downloading %d files\n", len(dlTargts))

		wg := new(sync.WaitGroup)
		for _, tgt := range dlTargts {
			obj := *tgt
			wg.Add(1)

			go func() {
				defer wg.Done()

				dlInput := &s3.GetObjectInput{
					Bucket: aws.String(bucket),
					Key:    obj.Key,
				}

				out, err := s3client.GetObject(dlInput)
				if err != nil {
					log.Error(err)
				}
				defer out.Body.Close()

				outputFileName := strings.Replace(*obj.Key, "/", "_", -1)

				destFile := filepath.Join(destDir, outputFileName)

				outf, err := os.Create(destFile)
				if err != nil {
					log.Error(err)
				}
				defer outf.Close()
				io.Copy(outf, out.Body)
			}()
		}
		wg.Wait()

		if len(rest) < bulkDlCount {
			dlTargts, rest = rest, nil
		} else {
			dlTargts, rest = rest[:bulkDlCount], rest[bulkDlCount:]
		}
	}

	return nil
}

// type GetObjectInput struct {
// 	_ struct{} `type:"structure"`

// 	// Bucket is a required field
// 	Bucket *string `location:"uri" locationName:"Bucket" type:"string" required:"true"`

// 	// Return the object only if its entity tag (ETag) is the same as the one specified,
// 	// otherwise return a 412 (precondition failed).
// 	IfMatch *string `location:"header" locationName:"If-Match" type:"string"`

// 	// Return the object only if it has been modified since the specified time,
// 	// otherwise return a 304 (not modified).
// 	IfModifiedSince *time.Time `location:"header" locationName:"If-Modified-Since" type:"timestamp" timestampFormat:"rfc822"`

// 	// Return the object only if its entity tag (ETag) is different from the one
// 	// specified, otherwise return a 304 (not modified).
// 	IfNoneMatch *string `location:"header" locationName:"If-None-Match" type:"string"`

// 	// Return the object only if it has not been modified since the specified time,
// 	// otherwise return a 412 (precondition failed).
// 	IfUnmodifiedSince *time.Time `location:"header" locationName:"If-Unmodified-Since" type:"timestamp" timestampFormat:"rfc822"`

// 	// Key is a required field
// 	Key *string `location:"uri" locationName:"Key" min:"1" type:"string" required:"true"`

// 	// Part number of the object being read. This is a positive integer between
// 	// 1 and 10,000. Effectively performs a 'ranged' GET request for the part specified.
// 	// Useful for downloading just a part of an object.
// 	PartNumber *int64 `location:"querystring" locationName:"partNumber" type:"integer"`

// 	// Downloads the specified range bytes of an object. For more information about
// 	// the HTTP Range header, go to http://www.w3.org/Protocols/rfc2616/rfc2616-sec14.html#sec14.35.
// 	Range *string `location:"header" locationName:"Range" type:"string"`

// 	// Confirms that the requester knows that she or he will be charged for the
// 	// request. Bucket owners need not specify this parameter in their requests.
// 	// Documentation on downloading objects from requester pays buckets can be found
// 	// at http://docs.aws.amazon.com/AmazonS3/latest/dev/ObjectsinRequesterPaysBuckets.html
// 	RequestPayer *string `location:"header" locationName:"x-amz-request-payer" type:"string" enum:"RequestPayer"`

// 	// Sets the Cache-Control header of the response.
// 	ResponseCacheControl *string `location:"querystring" locationName:"response-cache-control" type:"string"`

// 	// Sets the Content-Disposition header of the response
// 	ResponseContentDisposition *string `location:"querystring" locationName:"response-content-disposition" type:"string"`

// 	// Sets the Content-Encoding header of the response.
// 	ResponseContentEncoding *string `location:"querystring" locationName:"response-content-encoding" type:"string"`

// 	// Sets the Content-Language header of the response.
// 	ResponseContentLanguage *string `location:"querystring" locationName:"response-content-language" type:"string"`

// 	// Sets the Content-Type header of the response.
// 	ResponseContentType *string `location:"querystring" locationName:"response-content-type" type:"string"`

// 	// Sets the Expires header of the response.
// 	ResponseExpires *time.Time `location:"querystring" locationName:"response-expires" type:"timestamp" timestampFormat:"iso8601"`

// 	// Specifies the algorithm to use to when encrypting the object (e.g., AES256).
// 	SSECustomerAlgorithm *string `location:"header" locationName:"x-amz-server-side-encryption-customer-algorithm" type:"string"`

// 	// Specifies the customer-provided encryption key for Amazon S3 to use in encrypting
// 	// data. This value is used to store the object and then it is discarded; Amazon
// 	// does not store the encryption key. The key must be appropriate for use with
// 	// the algorithm specified in the x-amz-server-side​-encryption​-customer-algorithm
// 	// header.
// 	SSECustomerKey *string `location:"header" locationName:"x-amz-server-side-encryption-customer-key" type:"string"`

// 	// Specifies the 128-bit MD5 digest of the encryption key according to RFC 1321.
// 	// Amazon S3 uses this header for a message integrity check to ensure the encryption
// 	// key was transmitted without error.
// 	SSECustomerKeyMD5 *string `location:"header" locationName:"x-amz-server-side-encryption-customer-key-MD5" type:"string"`
//

// 	// VersionId used to reference a specific version of the object.
// 	VersionId *string `location:"querystring" locationName:"versionId" type:"string"`
// }
