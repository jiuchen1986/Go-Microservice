###########################################
### This is for building a docker image ###
### from an alpine-glibc base image     ###
### containing the binary of the test   ###
### microservice with golang            ###
###########################################
From frolvlad/alpine-glibc

WORKDIR /usr/src/app/
RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh
    
# RUN git config --global http.proxy http://www-proxy.ericsson.se:8080

RUN git clone -b multi-chain https://github.com/jiuchen1986/Go-Microservice.git && \
    mv Go-Microservice/main ./main && \
    chmod +x ./main && \
    rm -rf Go-Microservice
CMD ["/usr/src/app/main"]    
ENTRYPOINT ["/bin/bash", "-c"]

