
# <p align="center">gRPCis<br> ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/egorgasay/grpc-storage) ![GitHub issues](https://img.shields.io/github/issues/egorgasay/grpc-storage) ![License](https://img.shields.io/badge/license-MIT-green)</p>
This is a system consisting of several microservices (Memory Balancer, Storage, WebApplication), which is a distributed key-value database. There can be an unlimited number of Storage instances, they are all connected to the Memory Balancer via gRPC, which distributes the load between them. You can connect to the Web Application (Echo) via the Web interface to enter the necessary data manually. The system is fault-tolerant, guarantees complete data recovery even after a power outage.
<p align="center" >
<img src="https://user-images.githubusercontent.com/102957432/231838305-8903c6b4-8590-43a0-a7d5-ce7840be0070.png"  width="1000" />
</p>

# Drivers  
  
- Go - [grpcis-go-sdk](github.com/egorgasay/grpcis-go-sdk)

# Unique value search algorithm:  
  
By default, the value is saved to the minimally loaded server and returns its number to the Client. The client takes all the dirty work of storing and using this number for subsequent accesses to this key. The search for all instances occurs only when the -2 flag is specified and is performed in parallel gorutins.

# Transaction Logger
Protection against data loss in case of various hardware problems is achieved by using Transaction Logger. Each operation in the background is written to disk and performed again when the server is turned on after a failure (in other cases this does not happen).

The data is stored in a database or in a text file (depending on the type that was passed as a flag):  
  
```table
operation | key | value
```  

Important! When using a file logger, it is worth considering that when using keys and strings with a newline character, it is unacceptable.
# Quick start

## Server with Memory Balancer
```bash
go run cmd/balancer/main.go -a=':PORT' -d='MONGODB_URI'
```

## Server with WebApplication
```bash
go run cmd/grpc-storage-cli/main.go -a=':PORT' -b='BALANCER_IP:BALANCER_PORT'
```

## Server (or servers) with Storage instance
```bash
go run cmd/grpc-storage/main.go -a=':PORT' -d='MONGODB_URI' -connect='BALANCER_IP:BALANCER_PORT' -tlog_dir='DIR_FOR_TRANSACTION_LOGGER'
```
  
!!! DO NOT USE temporary directories for tlog_dir !!!

# Preview of the WebApplication  
(The launch of a web application can take up to 30 seconds)

## Go 
https://grpc-storage.egorpoletaikin.repl.co 
## PHP
https://grpc-web.egorpoletaikin.repl.co   

## Main page
## Go  
![изображение](https://user-images.githubusercontent.com/102957432/231824845-3c4f064d-2de9-433e-a616-05ca79edbef7.png)
  
## PHP

## Usage
set key value server(optional) - Sets the value to the storage.  
server > 0 - Save to exact server.  
server = 0 (default) - Automatic saving to a less loaded server.  
server = -1 - Direct saving to the database.  
server = -2 - Saving in all instances.  
server = -3 - Saving in all instances and DB.  
  
get key server(optional) - Gets the value from the storage.  
server > 0 - Search on a specific server (speed: fast).  
server = 0 (default) - Deep search (speed: slow). 
server = -1 - DB search (speed: medium). 


history - History of user actions.  
servers - List of active servers with stats.  
