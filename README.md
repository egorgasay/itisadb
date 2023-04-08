
# <p align="center">gRPCis<br> ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/egorgasay/grpc-storage) ![GitHub issues](https://img.shields.io/github/issues/egorgasay/grpc-storage) ![License](https://img.shields.io/badge/license-MIT-green)</p>
This is a system consisting of several microservices (Memory Balancer, Storage, WebApplication), 
which is a distributed key-value database. 

There can be an unlimited number of Storage instances, they are all connected to the Memory Balancer via gRPC, 
which distributes the load between them. 

You can connect to the Web Application (Echo) via the Web interface to enter the necessary data manually.
<p align="center" >
<img src="https://user-images.githubusercontent.com/102957432/229231202-cd39d983-7d8a-480c-9225-480422179e24.png"  width="400" />
</p>

# <p align="center">Quick start</p>

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
go run cmd/grpc-storage/main.go -a=':PORT' -d='MONGODB_URI' -connect='BALANCER_IP:BALANCER_PORT'
```

# <p align="center"> Preview of the WebApplication  </p>

## Demo: https://grpc-storage.egorpoletaikin.repl.co (The launch of a web application can take up to 30 seconds)
  
## Main page
  
![изображение](https://user-images.githubusercontent.com/102957432/230742612-d5e3876e-5cd0-496c-b7ee-7ef448e272bd.png)

## Help command
set key value server(optional) - Sets the value to the storage.  
server > 0 - Save to exact server.  
server = 0 (default) - Automatic saving to a less loaded server.  
server = -1 - Direct saving to the database.  
server = -2 - Saving in all instances.  
server = -3 - Saving in all instances and DB.  
  
get key server(optional) - Gets the value from the storage.  
server > 0 - Search on a specific server. (speed: fast)  
server = 0 (default) - Deep search. (speed: slow)  
server = -1 - DB search. (speed: medium)  

history - History of user actions.  
servers - List of active servers with stats.  
