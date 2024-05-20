FROM golang:alpine
LABEL authors="itisadb"
LABEL version="v0.0.1"

COPY . /itisadb
WORKDIR /itisadb

#RUN chmod +x bin/itisadb-all-linux-amd64 && \
#    mv bin/itisadb-all-linux-amd64 /bin/itisadb

RUN go build -o itisadb cmd/*.go
RUN chmod +x itisadb
RUN mv itisadb /bin/itisadb

EXPOSE 8888

CMD ["itisadb"]
