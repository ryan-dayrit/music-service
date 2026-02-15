# music-service
practice Golang application for gRPC, REST, Kafka, PostgreSQL, Protobuf

# Processing Flow 
Components 
1. PostgreSQL database in minikube 
2. Kafka in docker
3. gRPC API (internal)
4. REST API (external)
5. Kafka consumer
   
Writes 
1. REST API POST/PUT receiver for json payloads
2. REST API publishes proto to Kafka
3. Kafka consumer consumes the proto from kafka and creates/updates PostgreSQL

Reads 
1. gRPC API which reads from the PostgreSQL database using Sqlx framework and returns protos in json format
2. REST API which reads from the PostgreSQL database using ORM Go framework and returns protos in json format
   
# CLI Testers
1. REST API client which sends POST/PUT requests
2. gRPC API client which calls gRPC API and prints response
3. Kakfa producer which sends marshalled protos to the Kafka topic
4. PostgreSQL client which reads gets all albums from the PostgreSQL database direcly
5. PostgreSQL client which reads gets album by id from the PostgreSQL database direcly
6. PostgreSQL client which inserts PostgreSQL database direcly
7. REST API client which sends put/post requests to REST API album endpoint
8. REST API client which sends put/post requests to REST API albums endpoint
   
# rest api using Fiber 
  * https://dev.to/koddr/build-a-restful-api-on-go-fiber-postgresql-jwt-and-swagger-docs-in-isolated-docker-containers-475j
  * https://github.com/koddr/tutorial-go-fiber-rest-api
  * https://medium.com/@christian.asterisk/building-a-scalable-monolithic-backend-with-go-fiber-folder-structure-explained-5d1023eafa5e

# deploying kafka using docker compose 
  * https://docs.docker.com/guides/kafka/
  * https://medium.com/@darshak.kachchhi/setting-up-a-kafka-cluster-using-docker-compose-a-step-by-step-guide-a1ee5972b122

# installing kafka 
  * https://medium.com/@Shamimw/kafka-a-complete-tutorial-part-1-installing-kafka-server-without-zookeeper-kraft-mode-using-6fc60272457f

# deploying kafka to minikube 
  * https://medium.com/globant/deploying-kafka-on-minikube-without-ip-hack-springboot-producer-consumer-6698489012dd
  * [Strimzi](https://saedhasan.medium.com/setting-up-kafka-on-minikube-k8s-using-strimzi-5cac7870d943) 
  * https://www.google.com/search?q=how+to+access+Strimzi+kafka+from+outside+minikube
  * https://stackoverflow.com/questions/77480906/how-to-access-strimzi-kafka-cluster-running-on-minikube-publically
  * [Helm](https://www.google.com/search?q=deploying+kafka+to+minikube+using+helm)

# UIs for kafka
  * [Provectus Kafka UI](https://github.com/provectus/kafka-ui/)
  * https://www.reddit.com/r/apachekafka/comments/x9sov2/a_list_of_gui_tools_for_working_with_apache_kafka/
  * https://aiven.io/blog/top-kafka-ui

# structuring a repo for kafka processing 
  * https://www.google.com/search?q=how+to+structure+golang+repo+for+kafka+processing

# Steps to Generate a kustomization.yaml from Live Resources
  * https://www.google.com/search?q=how+to+get+kustomization.yaml+of+all+resources+in+k8s
  
# postgres in minikube 
  * https://www.google.com/search?q=how+to+host+postgres+in+minikube
  
# gRPC service in Kubernetes 
  * https://www.google.com/search?q=how+to+host+a+grpc+service+in+kubernetes

# Basic Self-Hosted Docker Registry 
  * https://www.google.com/search?q=personal+docker+registry
  * https://hub.docker.com/_/registry

# Makefile 
  * https://tutorialedge.net/golang/makefiles-for-go-developers/

# PostgreSQL
  * https://www.commandprompt.com/education/postgresql-list-all-tables/

# GOPRIVATE 
  * https://medium.com/@jasei/my-golang-experience-managing-multiple-private-repositories-with-github-and-docker-6d9f61452b81

# replace in go.mod 
  * https://thewebivore.com/using-replace-in-go-mod-to-point-to-your-local-module/
  
# Golang Project structure
  * https://github.com/golang-standards/project-layout/tree/master
  * https://pkg.go.dev/github.com/loveyourstack/northwind

# Golang packages 
  * https://medium.com/the-godev-corner/how-to-create-publish-a-go-public-package-9034e6bfe4a9

# SQL File Naming Conventions
  * https://www.google.com/search?q=naming+convention+for+insert+sql+scripts+in+golang+repo

# Repository / Data Access Pattern
  * https://go.dev/doc/tutorial/database-access
  * https://medium.com/@dewirahmawatie/connecting-to-postgresql-in-golang-59d7b208bad2
  * https://medium.com/@eikhapoetra/building-a-scalable-backend-with-the-repository-pattern-in-golang-4c30e735034c
  * https://medium.com/@rseanjustice/data-access-in-go-d39d8945b078
  * https://threedots.tech/post/database-transactions-in-go/

# Protocol Buffers 
  * https://grpc.io/docs/languages/go/basics/
  * https://dev.to/davidsbond/golang-structuring-repositories-with-protocol-buffers-3012
  * https://protobuf.dev/programming-guides/proto3/
  * https://github.com/grpc/grpc-go/blob/master/examples/route_guide/routeguide/route_guide.proto

# gRPC 
  * https://grpc.io/docs/languages/go/basics/
  * https://bbengfort.github.io/2017/03/secure-grpc/
  * https://grpc.io/docs/guides/auth/
  * https://github.com/sahansera/go-grpc

# Cobra CLI
  * https://www.jetbrains.com/guide/go/tutorials/cli-apps-go-cobra/creating_cli/
  * https://www.digitalocean.com/community/tutorials/how-to-use-the-cobra-package-in-go
  * https://cobra.dev/docs/tutorials/getting-started/
  * https://github.com/mwiater/golangcliscaffold/tree/step3




