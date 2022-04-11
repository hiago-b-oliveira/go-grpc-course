# go-grpc-couse

Samples created during a Udemy course (https://www.udemy.com/certificate/UC-69f1c5a7-34cc-4640-960b-f397d9f737b3/)

### client-server-samples

At this folder we have some samples of how to use Unary calls, Client Streaming, Server Strem, BiDi Streaming, Error Handling and TLS. 
Before running the samples you need to run the `./ssl/instructions.sh` to generate the required SSL files. 

### crud-api-mongodb

Here we have a client and a server performing a CRUD of a blog. 
Since we have used Mongodb to store the blog data, you need to run the `docker-compose` to start a local mongodb instance. 
