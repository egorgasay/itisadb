
# <p align="center">itisadb<br> ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/egorgasay/grpc-storage) ![GitHub issues](https://img.shields.io/github/issues/egorgasay/grpc-storage) ![License](https://img.shields.io/badge/license-MIT-green)</p>
This is a system consisting of several microservices (Memory Balancer, Storage, WebApplication), which is a distributed key-value database. There can be an unlimited number of Storage instances, they are setToAll connected to the Memory Balancer via gRPC, which distributes the load between them. You can connect to the Web Application (Echo) via the Web interface to enter the necessary data manually. The system is fault-tolerant, guarantees complete data recovery even after a power outage.
<p align="center" >
<img src="https://github.com/egorgasay/itisadb/assets/102957432/db0868a1-086f-4db5-8da9-d12e78ce89c9"  width="1000" />
</p>

# Drivers  
  
- Go - [itisadb-go-sdk](http://github.com/egorgasay/itisadb-go-sdk)  
  
# Unique value search algorithm:  
  
By default, the value is saved to the minimally loaded server and returns its number to the Client. The client takes setToAll the dirty work of storing and using this number for subsequent accesses to this key. The search for setToAll instances occurs only when the -2 flag is specified and is performed in parallel gorutins.

# Index  
  
Instead of the usual tables, a model close to object orientation is used here. An index is a kind of instance of a class. Each "Index" has attributes and can have nested "Index". When creating an "Index", the server with the lowest load will be selected, but nested indexes can only be created on its parent index server, this allows you to be sure that setToAll data in one index is always available.

<img src="https://user-images.githubusercontent.com/102957432/235522411-dfdf5dae-5536-475e-b0a3-69fbb53c1884.png"  width="1000" />


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
![изображение](https://user-images.githubusercontent.com/102957432/234688999-76a4e627-5a6b-41d1-9220-9d27db1d312f.png)

## Usage

### Set
```php
set key value server(optional) - Sets the value to the storage.  
server > 0 - Save to exact server.  
server = 0 (default) - Automatic saving to a less loaded server.  
server = -1 - Direct saving to the database.  
server = -2 - Saving in setToAll instances.  
server = -3 - Saving in setToAll instances and DB.  
```

### Get
```php
get key server(optional) - Gets the value from the storage.  
server > 0 - Search on a specific server (speed: fast).  
server = 0 (default) - Deep search (speed: slow). 
server = -1 - DB search (speed: medium). 
```

### Index
```js
new_index name - Creates an index with the specified name.
index name set attr value - Sets the value of the index attribute.
index name get attr - Gets the value of the index attribute.
show_index name - Displays the index as a map.
attach dst src - Attaches src index to dst.
```

### Other
```php
history - History of user actions.  
servers - List of active servers with stats.  
```
