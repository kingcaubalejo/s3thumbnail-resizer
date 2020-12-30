package s3thumbnailresizer

import (
	"bytes"
	"log"
	"os"
	"math/rand"
	"time"
	"fmt"
	"image/jpeg"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/nfnt/resize"
	
	"golang-thumbnail-creation/model"
)

func Create(options map[string]interface{}) error {
	
	s, err := session.NewSession(&aws.Config{
		Region: aws.String(fmt.Sprintf("%v", options["region"])),
	})
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Open(fmt.Sprintf("%v", options["file_dir"]))
	if err != nil {
		return err
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	var size int64 = fileInfo.Size()
	buffer := make([]byte, size)
	file.Read(buffer)

	_, errUpload := s3.New(s).PutObject(&s3.PutObjectInput{
		Bucket						: aws.String(fmt.Sprintf("%v", options["bucket_name"])),
		Key							: aws.String(fmt.Sprintf("%v", options["file_dir"])),
		ACL							: aws.String("private"),
		Body						: bytes.NewReader(buffer),
		ContentLength				: aws.Int64(size),
		ContentDisposition			: aws.String("attachment"),
		ServerSideEncryption		: aws.String("AES256"),
	})

	if errUpload != nil {
		return errUpload
	}

	file_details := model.FileDetails{
		FileId		: generateFileId(),
		Name		: fileInfo.Name(),
		Size		: fileInfo.Size(),
	}

	table_name := fmt.Sprintf("%v", options["table_name"])
	err = addMetaDataToDynamoDB(s, &file_details, table_name)
	log.Fatal(err)

	s3_bucket_url := fmt.Sprintf("https://%v.s3-%v.amazonaws.com/%v", options["bucket_name"], options["region"], options["file_dir"]) 
	response := fmt.Sprintf("UPLOAD TO %s BUCKET IS SUCCESSFULL\nAWS S3 LINK: %s\nIMAGE METADATA IS SUCCESSFULLY INSERTED TO TABLE %v\n", options["bucket_name"], s3_bucket_url, options["table_name"])
	
	fmt.Println(response)

	return err
}

func addMetaDataToDynamoDB(s *session.Session, fileDetails *model.FileDetails, table string) error {
	svc := dynamodb.New(s)

	av, err := dynamodbattribute.MarshalMap(fileDetails)
	if err != nil {
		log.Println("Got error in marshalling file details")
		log.Fatal(err)
		return err
	}

	fmt.Println(table)

	tableName := table
	input := &dynamodb.PutItemInput{
		Item: av,
		TableName: aws.String(tableName),
	}

	_, errDynamo := svc.PutItem(input)
	if errDynamo != nil {
		log.Println("Got error calling [PutItem] Check if table is existing.")
		log.Fatal(errDynamo)
	}
	return errDynamo
}

func generateFileId() string {
	rand.Seed(time.Now().UnixNano())
	min := 10
	max := 30
	
	return fmt.Sprintf("%05d", rand.Intn(max - min + 1) + min)
}

func ResizeImage(fileDir string, width uint, height uint) error {

	file, errOpen := os.Open(fileDir)
	if errOpen != nil {
		log.Println("[err-open-file001]")
		return errOpen
	}

	img, errResize := jpeg.Decode(file)
	if errResize != nil {
		log.Println("[err-resize-file001]")
		return errResize
	}

	var file_location = generateNewFileName(file.Name())
	
	m 					:= resize.Resize(width, height, img, resize.Lanczos3)
	_fDir 				:= file_location
	out, errCreate 		:= os.Create(_fDir)

	if errCreate != nil {
		log.Println("[err-open-file003]")
		return errCreate
	}
	defer out.Close()

	jpeg.Encode(out, m, nil)

	return nil
}

func generateNewFileName(fileDir string) string {
	r_backslash 	:= strings.Split(fileDir, "/")
	r_dot 			:= strings.Split(r_backslash[1], ".")
	file_name 		:= r_dot[0] + "-resized_001" + "." + r_dot[1]

	return r_backslash[0] + "/" + file_name
}