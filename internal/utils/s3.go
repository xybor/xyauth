package utils

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/xybor-x/xyerror"
)

func ReadS3(file string) ([]byte, error) {
	var path = file[5:]
	var bucket, item, found = strings.Cut(path, "/")
	if !found {
		return nil, xyerror.ValueError.Newf("not found item in path %s", path)
	}

	var sess, err = session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})

	if err != nil {
		return nil, err
	}

	var downloader = s3manager.NewDownloader(sess)
	var buf = aws.NewWriteAtBuffer([]byte{})
	_, err = downloader.Download(
		buf,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(item),
		})

	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
