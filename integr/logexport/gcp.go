package logexport

import (
	"fmt"
//	"log"
// Imports the Google Cloud Storage client package.
	"cloud.google.com/go/storage"
	"golang.org/x/net/context"
	//"google.golang.org/api/iterator"
	"time"
	//"github.com/shirou/gopsutil/host"

	"os"
	"io/ioutil"
	"io"
)

type GcpObjectStorage struct {
	bucket * storage.BucketHandle
	hostId string
	ctx context.Context
	client *storage.Client
}

func NewGcpObjectStorage(bucketName string) (*GcpObjectStorage ,error) {
	var err error
	ost := GcpObjectStorage{}
	ost.ctx = context.Background()
	ost.client, err = storage.NewClient(ost.ctx)
	if err != nil {
		fmt.Println("Failed to create client: %v", err)
		return nil , err
	}
	ost.bucket = ost.client.Bucket(bucketName)

	//hostInfo ,err := host.Info()
	//if err == nil {
	//	ost.hostId = hostInfo.HostID
	//}else {
	//	ost.hostId = "unknown"
	//}
	return &ost , nil
}

func (ost *GcpObjectStorage)UploadLogSnapshot(files []string,username string,sizeLimit int64)[]string{
	ts := time.Now().Format("2006-01-02T15:04")
	snapshotName := username+"/"+ts
	var statusReport []string
	for i := range files {
		metadata := map[string]string {
			"fimpui-username": username,
		}
		err := ost.UploadTextFile(snapshotName,files[i],metadata,sizeLimit)
		uploadStatus := "OK"
		if err != nil {
			uploadStatus = err.Error()
			fmt.Println("Error while uploading file :",err)
		}

		statusReport = append(statusReport,fmt.Sprintf("File: %s , upload status = %s ",files[i],uploadStatus))
	}
	return statusReport
}


func (ost * GcpObjectStorage)UploadTextFile(objectPrefix string,filePath string , metadata map[string]string,sizeLimit int64) error {

	var fileBody []byte
	var err error
	if sizeLimit > 0 {
		file, err := os.Open(filePath)
		if err != nil {
			fmt.Println("Can't open file .Error:",err)
			return err
		}
		defer file.Close()

		buf := make([]byte,sizeLimit)
		stat, err := os.Stat(filePath)
		if err != nil {
			fmt.Println("Can't stat file .Error:",err)
			return err
		}
		var start int64
		start = 0
		if stat.Size()>sizeLimit {
			start = stat.Size() - sizeLimit
		}

		fileSize, err := file.ReadAt(buf, start)
		fmt.Println("NUmber of bytes were read = ",fileSize)
		if err != nil && err != io.EOF {
			fmt.Println(err)
			return err
		}
		fileBody = buf[:fileSize]
	} else {
		fileBody , err = ioutil.ReadFile(filePath)

	}
	if err != nil {
		return err
	}
	//fmt.Println(filepath.Dir(filePath))
	//name := objectPrefix+filepath.Dir(filePath)+"/"+ts+"_"+filepath.Base(filePath)
	name := objectPrefix+filePath
	wc := ost.bucket.Object(name).NewWriter(ost.ctx)
	wc.ContentType = "text/plain"
	//wc.Metadata = map[string]string{
	//	"fimpui-username": username,
	//}
	if metadata != nil {
		wc.Metadata = metadata
	}
	if _, err := wc.Write(fileBody); err != nil {
		fmt.Printf("createFile: unable to write data to bucket %q, file %q: %v", ost.bucket, "test", err)
		return err
	}
	err = wc.Close()
	if err != nil {
		fmt.Println("Failed to create object: %v", err)
	}
	return nil

}

//func UploadLogToGcp() {
//	ctx := context.Background()
//
//	// Sets your Google Cloud Platform project ID.
//	//projectID := "63247381933"
//	// Creates a client.
//	client, err := storage.NewClient(ctx)
//	if err != nil {
//		log.Fatalf("Failed to create client: %v", err)
//	}
//
//	// Sets the name for the new bucket.
//	bucketName := "fh-cube-log"
//
//	// Creates a Bucket instance.
//	bucket := client.Bucket(bucketName)
//	wc := bucket.Object("test").NewWriter(ctx)
//	wc.ContentType = "text/plain"
//
//	wc.Metadata = map[string]string{
//		"x-goog-meta-foo": "foo",
//		"x-goog-meta-bar": "bar",
//	}
//
//	if _, err := wc.Write([]byte("abcde\n")); err != nil {
//		fmt.Printf("createFile: unable to write data to bucket %q, file %q: %v", bucketName, "test", err)
//		return
//	}
//	err = wc.Close()
//	if err != nil {
//		log.Fatalf("Failed to create object: %v", err)
//	}
//	//// Creates the new bucket.
//	//if err := bucket.Create(ctx, projectID, nil); err != nil {
//	//	log.Fatalf("Failed to create bucket: %v", err)
//	//}
//
//
//	fmt.Printf("Bucket %v created.\n", bucketName)
//}


