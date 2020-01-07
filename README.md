# MongoDB Field Level Encryption (FLE) Tutorial/Demo


(IN PROGRESS!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!)




Demo MongoDB Client-Side Field Level Encryption. Uses Golang + Ubuntu in a Docker container


## Run 

Note: The Dockerfile contains all environment dependencies to run this demo.


1. Add values to required variables specified in `env.list.example` and rename the file to `env.list`

2. Run the following:
```bash
docker run --rm  -it  -p 8888:8888 --env-file ./env.list --hostname fle  nullstring/mongo-fle-demo
```
