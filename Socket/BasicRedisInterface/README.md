
Redis Clients don't just send `hello 3`. They send additional metadata - i'd imagine this is done to inform redis how much resources or data to be expected to come over. WIth that, the redis server can then allocate the right amount of resources to store said items.

```
Start server
Obtained input: *2
-ERR unknown command '*2', with args beginning with []:
Obtained input: $5
-ERR unknown command '$5', with args beginning with []:
Obtained input: hello
-ERR unknown command 'hello', with args beginning with []:
Obtained input: $1
-ERR unknown command '$1', with args beginning with []:
Obtained input: 3
-ERR unknown command '3', with args beginning with []:
```