FROM golang:alpine
LABEL authors="itisadb"
LABEL version="v0.0.1"

COPY . /itisadb
WORKDIR /itisadb

RUN chmod +x bin/itisadb-all-linux-amd64 && \
    mv bin/itisadb-all-linux-amd64 /bin/itisadb
EXPOSE 8888

CMD ["itisadb"]
