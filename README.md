# music-service
practice Golang application for gRPC, REST, Kafka, PostgreSQL, Protobuf

# Processing Flow 
Components 
1. PostgreSQL database in minikube with DDL/DML initialization scripts
2. Kafka in docker with a topic and a consumer group
3. gRPC API (internal) to service the data in the Postgres datanase 
4. REST API (external) using fiber library to receive JSON payloads and enqueue a Protobuf message to a kakfa topic
5. Kafka consumer using IBM/sarama library to process Protobuf messages from a Kafka topic
6. Kafka consumer using confluent library to process Protobuf messages from a Kafka topic
   
Writes 
1. REST API POST/PUT receiver for json payloads
2. REST API publishes proto to Kafka
3. Kafka consumers consume the proto from the kafka topic and creates/updates PostgreSQL

Reads 
1. gRPC API which reads from the PostgreSQL database using Sqlx library and returns protos in json format
2. REST API which reads from the PostgreSQL database using ORM library and returns protos in json format
   
# CLI Testers
1. REST API client which sends POST/PUT requests
2. gRPC API client which calls gRPC API and prints response
3. Kakfa producer which sends marshalled protos to the Kafka topic
4. PostgreSQL client which reads gets all albums from the PostgreSQL database direcly
5. PostgreSQL client which reads gets album by id from the PostgreSQL database direcly
6. PostgreSQL client which inserts PostgreSQL database direcly
7. REST API client which sends put/post requests to REST API album endpoint
8. REST API client which sends put/post requests to REST API albums endpoint

