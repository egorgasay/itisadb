# ItisaDB ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/egorgasay/grpc-storage) ![GitHub issues](https://img.shields.io/github/issues/egorgasay/grpc-storage) ![License](https://img.shields.io/badge/license-MIT-green)

This is a system consisting of several microservices (Memory Balancer, Storage, WebApplication), which is a distributed key-value database. 

There can be an unlimited number of Storage instances, they are connected to the Memory Balancer via gRPC, which distributes the load between them. You can connect to the Web Application (Echo) via the Web interface to enter the necessary data manually. 

The system is fault-tolerant, guarantees complete data recovery even after a power outage.
<p align="center" >
<img src="https://github.com/egorgasay/itisadb/assets/102957432/2fe84aea-9068-4615-bb16-da94d8277aad"  width="1000" />
</p>

# Documentation

Documentation is available on [doc.itisadb](https://208bc35d-fa36-4993-bbf4-e05bd8ca153e-00-oamkdjnhcl6s.worf.replit.dev)

# Demo Version

The demo version is available on [demo.itisadb](https://f273a4ca-47cb-4cba-8916-38b8866d1c49-00-17hh7nh99cz9s.worf.replit.dev)

# Drivers  
  
- Go - [itisadb-go-sdk](http://github.com/egorgasay/itisadb-go-sdk)  
  
# Quick start

```bash
go run cmd/itisadb/main.go
```
