
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

## Demo: https://grpc-storage.egorpoletaikin.repl.co
  
## Main page
  
![img.png](img.png)

## Help command
set key value   
get key  
history
