## Sync Gateway User Migration Script
🚀 Features built in concurrency to speed up your migration!!!

Simple script to migrate users from one [Sync Gateway](https://docs.couchbase.com/sync-gateway/current/rest-api-admin.html) instance to another. 

#### Download users from your source Couchbase instance 
into a `files/sample.json` file with the following structure:

```json
[
    {
      "admin_channels": [
        "some",
        "admin",
        "channels"
      ],
      "admin_roles": [
        "some",
        "admin",
        "roles"
      ],
      "name": "some_username"
   }
]
```

### Fetch Dependecies

```bash
go get
```

### Build the binary

```bash
go build main.go
```

### Run the binary 
with `url` and `db` flags for the destination Sync Gateway Instance

```bash
./main -url=my_url -db=my_db
```

😊 Happy Migrating!
