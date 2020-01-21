# MongoDB Field Level Encryption (FLE) Tutorial/Demo

Demo MongoDB Client-Side Field Level Encryption. Uses Golang + Ubuntu in a Docker container.  


## Run 

Note: The Dockerfile contains __all__ environment dependencies to run this demo.

1. Add values to required variables specified in `env.list.example` and rename the file to `env.list`
    - Need to have a MongoDB deployment running, if not, delopoy a [free one in Atlas](https://cloud.mongodb.com/user#/atlas/register/accountProfile) and grab the connection string
    - Need to have [AWS KMS](https://aws.amazon.com/kms/) configured 

2. Run the following:
```bash
docker run --rm  -it  -p 8888:8888 -p 27020:27020 --env-file env.list --hostname fle  nullstring/mongo-fle-demo
```

## `foobar` document

```json
{
    "_id": "string",
    "name":"string",
    "message": "string" 
}
```
Note: `message` is encrypted/decrypted if inserted/read via /foo else as-is.

## Endpoints

- `POST /foo`  -- Inserts a valid `foobar` document to the `tutorial.foobar` namespace and encrypts the `message` field.
- `GET /foo/{id}` -- Reads a `foobar` document with matching `id` and attempts to decrypt the `message` field.

- `POST /bar` -- Inserts a valid `foobar` document to the `tutorial.foobar` namespace. (sans encryption)
- `GET /bar/{id}` -- Reads a `foobar` document with matching `id` as-is. (sans decryption)


## Test

Import [Postman collection]().

For debugging/ad-hoc testing:
```bash
git clone https://github.com/desteves/fle.git
cd fle
docker run --rm -it -v $PWD:/go/src/github.com/desteves/fle --entrypoint /bin/bash -p 8777:8888  -p 27020:27020 --env-file env.list --hostname fle-testing nullstring/mongo-fle-demo
go build -tags cse main.go
./main
```


## References

- [General mongoDB Docs on FLE](https://docs.mongodb.com/manual/core/security-client-side-encryption/)
- [mongoDB University Guide](https://github.com/mongodb-university/csfle-guides)
- [mongoDB Labs go example](https://github.com/mongodb-labs/field-level-encryption-sandbox/tree/master/go)
- [mongoDB FLE Use Case Guide](https://docs.mongodb.com/ecosystem/use-cases/client-side-field-level-encryption-guide/)
- [mongoDB Using KMS](https://docs.mongodb.com/ecosystem/use-cases/client-side-field-level-encryption-local-key-to-kms/)
- [mongoDB Driver FLE](https://godoc.org/go.mongodb.org/mongo-driver/mongo#hdr-Client_Side_Encryption)
- [mongoDB libmongocrypt](https://github.com/mongodb/libmongocrypt#installing-libmongocrypt-from-distribution-packages)
- [Go Driver Tests on FLE](https://github.com/mongodb/mongo-go-driver/blob/c5b8476622aec25b142e39ae7cb3e6787ccabc74/data/client-side-encryption/README.rst)
