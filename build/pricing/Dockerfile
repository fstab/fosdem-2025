FROM ubuntu:latest

RUN apt-get update && apt-get upgrade -y

COPY ./bin/pricing-service /pricing-service

CMD ["/pricing-service"]
