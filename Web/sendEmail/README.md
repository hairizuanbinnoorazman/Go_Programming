# Send Email

Reference: https://gist.github.com/douglasmakey/90753ecf37ac10c25873825097f46300

```bash
curl -X POST localhost:8080/send-email -d '{"to": "test@test.com", "subject": "This is another test", "body": "Report Generated", "report_filename": "Test.pdf"}'
```

If you're using mailhog to test this, you can download the base64 encoded file and then run the command:

IMPORTANT: REMOVE THE METADATA

```bash
base64 -i {input file} -d > {output file}
```

## Quickstart

To run the application

```bash
EMAIL_HOST=127.0.0.1 EMAIL_PORT=25 EMAIL_USER=test EMAIL_PASS=test BUCKET_NAME=xxx go run main.go
```