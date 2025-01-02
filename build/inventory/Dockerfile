FROM ubuntu:latest

RUN apt-get update && apt-get upgrade -y

COPY ./bin/inventory-service /inventory-service

CMD ["/inventory-service"]
