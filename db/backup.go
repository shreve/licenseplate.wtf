package db

import (
	"log"
	"os"
	"time"
	_ "time/tzdata"

	// Additional imports needed for examples below

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/c2h5oh/datasize"
)

var daemonRunning = false

func s3Client() *s3.S3 {
	key := os.Getenv("SPACES_KEY")
	secret := os.Getenv("SPACES_SECRET")

	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(key, secret, ""),
		Endpoint:         aws.String("https://nyc3.digitaloceanspaces.com"),
		Region:           aws.String("us-east-1"),
		S3ForcePathStyle: aws.Bool(false), // // Configures to use subdomain/virtual calling format. Depending on your version, alternatively use o.UsePathStyle = false
	}

	newSession := session.New(s3Config)
	return s3.New(newSession)
}

func StartDaemon() error {
	if daemonRunning {
		return nil
	}
	daemonRunning = true

	// If we can't find data.sql, restore from backup
	stat, err := os.Stat(LOCATION)

	empty := stat.Size() == 0

	if err != nil || empty {
		log.Printf("No %s found, restoring from backup", LOCATION)
		err := Restore()
		if err != nil {
			log.Printf("Failed to restore database: %v", err)
			return err
		}
	} else {
		size := datasize.ByteSize(stat.Size())
		log.Printf("Database already found in place (%s)", size.String())
	}

	log.Println("Starting backup daemon.")

	// Otherwise, start the daemon
	go Daemon()
	return nil
}

func Daemon() {
	ET, _ := time.LoadLocation("America/Detroit")
	// ET := time.FixedZone("America/Detroit", -5*60*60)

	for {
		// Calculate the next midnight based on current time
		now := time.Now()
		midnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, ET)

		// Wait until then
		wait := midnight.Sub(time.Now())
		log.Printf("Backing up at %s in %s seconds", midnight, wait)
		time.Sleep(wait)

		log.Printf("Backing up now")
		err := Backup()
		if err != nil {
			log.Printf("Failed to backup database: %v", err)
		}
	}
}

func Backup() error {
	// If we can't find data.sql, there's nothing to back up
	stat, err := os.Stat(LOCATION)
	if os.IsNotExist(err) {
		log.Printf("No %s found, can't run backup", LOCATION)
		return err
	}

	// Use S3 to upload data.sql to a bucket called licenseplate-wtf
	client := s3Client()
	uploader := s3manager.NewUploaderWithClient(client)

	file, err := os.Open(LOCATION)
	if err != nil {
		log.Printf("Failed to open data.sql for backup: %v", err)
		return err
	}
	defer file.Close()

	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String("licenseplate-wtf"),
		Key:    aws.String("backup/" + LOCATION),
		Body:   file,
	})

	if err != nil {
		log.Printf("Failed to upload %s to S3: %v", LOCATION, err)
		return err
	}

	size := datasize.ByteSize(stat.Size())
	log.Printf("Backup uploaded to %s (%s)", aws.StringValue(&result.Location), size.String())
	return nil
}

func Restore() error {
	// Download the backup from S3 and restore it in place
	client := s3Client()
	downloader := s3manager.NewDownloaderWithClient(client)

	file, err := os.Create(LOCATION + ".bak")
	if err != nil {
		log.Printf("Failed to create %s.bak for restore: %v", LOCATION, err)
		return err
	}

	n, err := downloader.Download(file, &s3.GetObjectInput{
		Bucket: aws.String("licenseplate-wtf"),
		Key:    aws.String("backup/" + LOCATION),
	})
	if err != nil {
		log.Printf("Failed to download backup from S3: %v", err)
		return err
	}
	file.Close()

	// Prevent any further reads or writes while we move the file in place
	Lock.Lock()
	defer Lock.Unlock()

	err = os.Rename(LOCATION+".bak", LOCATION)
	if err != nil {
		log.Printf("Failed to replace data.sql from backup: %v", err)
		return err
	}

	size := datasize.ByteSize(n)
	log.Printf("Restored %s from backup (%v)", LOCATION, size.String())
	return nil
}
