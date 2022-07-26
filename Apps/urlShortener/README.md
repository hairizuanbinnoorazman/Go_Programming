# URL Shortener

This is a sample URL Shortener application for demonstration purposes. Definitely do not run this production due to various concerns with regards to how the application is run as well as it being incapable of processing multiple requests concurrently safely

## Testing

```bash
curl localhost:8080/add -X POST -d '{"url":"https://www.google.com"}'

curl localhost:8080/r/XXXXX -X DELETE

// To delete https://www.google.com
curl http://localhost:8080/r/ef7efc9839 -X DELETE
```