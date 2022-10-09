package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/smtp"
	"os"
	"path/filepath"
	"strings"

	"cloud.google.com/go/storage"
	"github.com/gorilla/mux"
)

func main() {
	log.Println("Begin server")

	emailHost := os.Getenv("EMAIL_HOST")
	emailPort := os.Getenv("EMAIL_PORT")
	emailServer := emailHost + ":" + emailPort
	log.Printf("Email host provided: %v\n", emailServer)
	emailUser := os.Getenv("EMAIL_USER")
	emailPass := os.Getenv("EMAIL_PASS")
	if emailPass == "" || emailUser == "" {
		panic("user and password expected to utilize this service")
	}
	bucketName := os.Getenv("BUCKET_NAME")
	log.Printf("Bucket name provided: %v\n", bucketName)

	cl, err := storage.NewClient(context.Background())
	if err != nil {
		panic("unexpected error starting gcs client")
	}

	r := mux.NewRouter()
	r.Handle("/send-email", &SendEmail{
		BucketName:  bucketName,
		StorageSvc:  cl,
		EmailServer: emailServer,
		EmailUser:   emailUser,
		EmailPass:   emailPass,
	}).Methods("POST")
	srv := &http.Server{
		Handler: r,
		Addr:    "0.0.0.0:8080",
	}
	log.Fatal(srv.ListenAndServe())
}

type Message struct {
	To          []string
	CC          []string
	BCC         []string
	Subject     string
	Body        string
	Attachments map[string][]byte
}

func NewMessage(s, b string) *Message {
	return &Message{Subject: s, Body: b, Attachments: make(map[string][]byte)}
}

func (m *Message) AttachFile(src string) error {
	b, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	_, fileName := filepath.Split(src)
	m.Attachments[fileName] = b
	return nil
}

func (m *Message) ToBytes() []byte {
	buf := bytes.NewBuffer(nil)
	withAttachments := len(m.Attachments) > 0
	buf.WriteString(fmt.Sprintf("Subject: %s\n", m.Subject))
	buf.WriteString(fmt.Sprintf("To: %s\n", strings.Join(m.To, ",")))
	if len(m.CC) > 0 {
		buf.WriteString(fmt.Sprintf("Cc: %s\n", strings.Join(m.CC, ",")))
	}

	if len(m.BCC) > 0 {
		buf.WriteString(fmt.Sprintf("Bcc: %s\n", strings.Join(m.BCC, ",")))
	}

	buf.WriteString("MIME-Version: 1.0\n")
	writer := multipart.NewWriter(buf)
	boundary := writer.Boundary()
	if withAttachments {
		buf.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\n", boundary))
		buf.WriteString(fmt.Sprintf("--%s\n", boundary))
	} else {
		buf.WriteString("Content-Type: text/plain; charset=utf-8\n")
	}

	buf.WriteString(m.Body)
	if withAttachments {
		for k, v := range m.Attachments {
			buf.WriteString(fmt.Sprintf("\n\n--%s\n", boundary))
			buf.WriteString(fmt.Sprintf("Content-Type: %s\n", http.DetectContentType(v)))
			buf.WriteString("Content-Transfer-Encoding: base64\n")
			buf.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=%s\n", k))

			b := make([]byte, base64.StdEncoding.EncodedLen(len(v)))
			base64.StdEncoding.Encode(b, v)
			buf.Write(b)
			buf.WriteString(fmt.Sprintf("\n--%s", boundary))
		}

		// buf.WriteString("--")
	}

	return buf.Bytes()
}

type SendEmail struct {
	BucketName  string
	StorageSvc  *storage.Client
	EmailServer string
	EmailUser   string
	EmailPass   string
}

func (s *SendEmail) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("start send-email")
	defer log.Println("end send-email")

	ctx := r.Context()

	type SendEmailRequest struct {
		To             []string `json:"to"`
		Subject        string   `json:"subject"`
		Body           string   `json:"body"`
		ReportFileName string   `json:"report_filename"`
	}

	var sendEmailReq SendEmailRequest
	raw, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(raw, &sendEmailReq)

	log.Printf("%+v\n", sendEmailReq)

	if sendEmailReq.ReportFileName != "" {
		err := s.downloadFile(ctx, sendEmailReq.ReportFileName)
		if err != nil {
			w.Write([]byte(fmt.Sprintf("unable to download file from gcs :: %v\n", err)))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer func() {
			os.Remove(sendEmailReq.ReportFileName)
		}()
	}

	m := NewMessage("Test", sendEmailReq.Body)
	m.To = sendEmailReq.To
	m.Subject = sendEmailReq.Subject
	if sendEmailReq.ReportFileName != "" {
		lol, _ := ioutil.ReadFile(sendEmailReq.ReportFileName)
		m.Attachments = map[string][]byte{
			sendEmailReq.ReportFileName: lol,
		}
	}

	err := smtp.SendMail(s.EmailServer, smtp.CRAMMD5Auth(s.EmailUser, s.EmailPass), "zzz@zzz.com", m.To, m.ToBytes())
	if err != nil {
		log.Println(err)
	}

	// buf := bytes.NewBufferString("This is the email body.")
	// if _, err = buf.WriteTo(wc); err != nil {
	// 	log.Fatal(err)
	// }
	type sendEmailResponse struct {
		Status string `json:"status"`
	}

	sendEmailRespRaw, _ := json.Marshal(sendEmailResponse{Status: "Sent"})
	w.Write(sendEmailRespRaw)
	w.WriteHeader(http.StatusOK)
}

func (s *SendEmail) downloadFile(ctx context.Context, fileName string) error {
	f, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("os.Open: %v", err)
	}
	defer f.Close()

	o := s.StorageSvc.Bucket(s.BucketName).Object(fileName)
	rc, err := o.NewReader(ctx)
	if err != nil {
		return fmt.Errorf("unable to create storage reader")
	}
	if _, err := io.Copy(f, rc); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	if err := rc.Close(); err != nil {
		return fmt.Errorf("reader.Close: %v", err)
	}
	return nil
}
