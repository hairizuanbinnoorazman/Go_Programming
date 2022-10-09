# Create Report

Run locally:

```bash
BUCKET_NAME=XXXX go run main.go
```

To create a report

```bash
curl localhost:8080/create-report -X POST -d '{"title": "Test", "description": "test", "image": "lol.jpeg", "template_file_name": "haha.md"}'
```
