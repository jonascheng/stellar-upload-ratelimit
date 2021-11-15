package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/juju/ratelimit"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	uploadServer = kingpin.Flag("server", "Upload server endpoint.").Default("http://localhost:8080/upload?token=secret").String()
	srcFile      = kingpin.Flag("file", "Upload source file.").Default("./data/agent-telemetry-threat-info-flat-500-1636954894.json.gz").String()
	rateLimit    = kingpin.Flag("rate", "Upload rate limit in Mbps.").Default("5").Float64()
)

// wrapper to limit io.Reader speed
func NewRateReader(reader io.Reader, ratelimitMbs float64) io.Reader {
	NewRateReader := reader

	if ratelimitMbs > 0 {
		log.Printf("upload file in %vMbs\n", ratelimitMbs)
		tokenPS := ratelimitMbs / 8 * 1024 * 1024
		log.Printf("adding %v bucket tokens per second\n", tokenPS)
		bucket := ratelimit.NewBucketWithRate(float64(tokenPS), int64(tokenPS))
		NewRateReader = ratelimit.Reader(reader, bucket)
	}

	return NewRateReader
}

func UploadFile(server string, filePath string, ratelimitMbs float64) {
	// create file reader
	file, err := os.Open(filePath)
	checkError(err)
	defer file.Close()

	body := &bytes.Buffer{}

	// create multipart writter
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", fmt.Sprintf("%v-%v", time.Now().Unix(), filepath.Base(file.Name())))
	checkError(err)
	written, err := io.Copy(part, file)
	checkError(err)
	log.Printf("%v bytes has been copy to io.Writer\n", written)
	writer.Close()

	// create ratelimit reader
	ratelimitBody := NewRateReader(body, ratelimitMbs)

	log.Println("start http.NewRequest")
	req, err := http.NewRequest("POST", server, ratelimitBody)
	log.Println("end http.NewRequest")
	checkError(err)

	log.Println("start req.Header.Add")
	req.Header.Add("Content-Type", writer.FormDataContentType())
	log.Println("end req.Header.Add")

	client := &http.Client{}
	// transport := http.DefaultTransport.(*http.Transport)
	// log.Printf("%v\n", transport)
	// transport.WriteBufferSize = 1 * 100
	// transport.ReadBufferSize = 1 * 100
	// log.Printf("%v\n", transport)
	// client.Transport = transport

	{
		quit := make(chan bool)
		log.Println("start client.Do")
		go dots(quit)
		resp, err := client.Do(req)
		quit <- true
		log.Println("\nend client.Do")
		checkError(err)
		log.Println(*resp)
	}

}

func main() {
	kingpin.Version("1.0.0")
	kingpin.Parse()

	UploadFile(*uploadServer, *srcFile, *rateLimit)
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func dots(quit chan bool) {
	for {
		select {
		case <-quit:
			return
		default:
			time.Sleep(time.Second)
			fmt.Print(".")
		}
	}
}
