# ItisaDB ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/egorgasay/grpc-storage) ![GitHub issues](https://img.shields.io/github/issues/egorgasay/grpc-storage) ![License](https://img.shields.io/badge/license-MIT-green)

This is a system consisting of several microservices (Memory Balancer, Storage, WebApplication), which is a distributed key-value database. There can be an unlimited number of Storage instances, they are connected to the Memory Balancer via gRPC, which distributes the load between them. You can connect to the Web Application (Echo) via the Web interface to enter the necessary data manually. The system is fault-tolerant, guarantees complete data recovery even after a power outage.
<p align="center" >
<img src="https://github.com/egorgasay/itisadb/assets/102957432/2fe84aea-9068-4615-bb16-da94d8277aad"  width="1000" />
</p>

# Drivers  
  
- Go - [itisadb-go-sdk](http://github.com/egorgasay/itisadb-go-sdk)  
  
# Unique value search algorithm:  
  
By default, the value is saved to the minimally loaded server and returns its number to the Client. The client takes setToAll the dirty work of storing and using this number for subsequent accesses to this key. The search for setToAll instances occurs only when the -2 flag is specified and is performed in parallel gorutins.

# Objects  
  
Instead of the usual tables, a model close to object orientation is used here. An object is a kind of instance of a class. Each "Object" has attributes and can have nested "Objects". When creating an "Object", the server with the lowest load will be selected, but nested objects can only be created on its parent object server, this allows you to be sure that all data in one object is always available.

<img src="https://github.com/egorgasay/itisadb/assets/102957432/ad8cd23b-3b4a-49c9-895f-3917ead8a7c7"  width="1000" />


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
go run cmd/itisadb/main.go
```
  
!!! DO NOT USE temporary directories for tlog_dir !!!

## Syntax

### SET - Sets the value to the storage.
```php
SET key "value" [ MODE - NX | RO | XX ] [ LEVEL - R | S ] [ SERVER - [0-9]+ ]  

MODE - Defines the mode of the operation. 
- `NX` - If the key already exists, it won't be overwritten. 
- `RO` - If the key already exists, an error will be returned.
- `XX` - If the key doesn't exist, it won't be created.

LEVEL - Defines the level of permission. 
- `R` (Restricted) - NO encryption, ACL validation
- `S` (Secret) - encryption, ACL validation
By default - NO encryption, NO ACL validation

SERVER - Defines server number to use. 
- Automaticly saving to a less loaded server by default.
```

### GET - Gets the value from the storage.
```php
GET key [ FROM - [0-9]+ ] [ LEVEL - D | R | L ]  

LEVEL - Defines the level of permission.
- `D` (Default) - NO encryption, NO ACL validation
- `R` (Restricted) - NO encryption, ACL validation
- `S` (Secret) - encryption, ACL validation

FROM - Defines server number to use.
- `> 0` - Search on a specific server (speed: fast).  
- `= 0` (default) - Deep search (speed: slow). 
```

### DEL - Deletes the key-value pair from the storage.
```php
DEL key [ FROM - [0-9]+ ] [ LEVEL - D | R | L ]  

LEVEL - Defines the level of permission.
- `D` (Default) - NO encryption, NO ACL validation
- `R` (Restricted) - NO encryption, ACL validation
- `S` (Secret) - encryption, ACL validation

FROM - Defines server number to use.
- `> 0` - Search on a specific server (speed: fast).  
- `= 0` (default) - Deep search (speed: slow). 
```

### NEW OBJECT - Creates an object with the specified name.
```php
NEW OBJECT name [ ON - [0-9]+ ] [ LEVEL - D | R | L ]

LEVEL - Defines the level of permission.
- `D` (Default) - NO encryption, NO ACL validation
- `R` (Restricted) - NO encryption, ACL validation
- `S` (Secret) - encryption, ACL validation

ON - Defines server number to use. 
- Automaticly saving to a less loaded server by default.
```

### DELETE OBJECT - Deletes the object with the specified name.
```php
DELETE OBJECT name [ ON - [0-9]+ ] [ LEVEL - D | R | L ]

LEVEL - Defines the level of permission.
- `D` (Default) - NO encryption, NO ACL validation
- `R` (Restricted) - NO encryption, ACL validation
- `S` (Secret) - encryption, ACL validation

ON - Defines server number to use. 
- Automaticly saving to a less loaded server by default.
```

### SETO - Sets the value to the object.
```php
SETO name key "value" [ LEVEL - D | R | L ]

LEVEL - Defines the level of permission.
- `D` (Default) - NO encryption, NO ACL validation
- `R` (Restricted) - NO encryption, ACL validation
- `S` (Secret) - encryption, ACL validation
```

### GETO - Gets the value from the object.
```php
GETO name key [ LEVEL - D | R | L ]

LEVEL - Defines the level of permission.
- `D` (Default) - NO encryption, NO ACL validation
- `R` (Restricted) - NO encryption, ACL validation
- `S` (Secret) - encryption, ACL validation
```

### DELO - Deletes the object key.
```php
DELO name key [ LEVEL - D | R | L ]

LEVEL - Defines the level of permission.
- `D` (Default) - NO encryption, NO ACL validation
- `R` (Restricted) - NO encryption, ACL validation
- `S` (Secret) - encryption, ACL validation
```

### MARSHAL OBJECT - Displays the object as JSON.
```php
MARSHAL OBJECT name [ LEVEL - D | R | L ]

LEVEL - Defines the level of permission.
- `D` (Default) - NO encryption, NO ACL validation
- `R` (Restricted) - NO encryption, ACL validation
- `S` (Secret) - encryption, ACL validation
```

### ATTACH - Attaches the src object to the dst.
```php
ATTACH dst [ LEVEL - D | R | L ] src [ LEVEL - D | R | L ]

LEVEL - Defines the level of permission.
- `D` (Default) - NO encryption, NO ACL validation
- `R` (Restricted) - NO encryption, ACL validation
- `S` (Secret) - encryption, ACL validation
```

### Other
```php
HISTORY - History of user actions.  
SERVERS - List of active servers with stats.  
```

# Preview of the WebApplication
(The launch of a web application can take up to 30 seconds)
(still on v0.7, use index keyword instead of object)
https://grpc-storage.egorpoletaikin.repl.co

## Main page
![изображение](https://user-images.githubusercontent.com/102957432/231824845-3c4f064d-2de9-433e-a616-05ca79edbef7.png)
