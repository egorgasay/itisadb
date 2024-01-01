# Basic commands

### SET - Sets the value to the storage.
```go
//              MODE                  LEVEL     SERVER
SET key "value" [ RO | UQ | NX | XX ] [ R | S ] [ [0-9]+ ]
```

`MODE` - Defines the mode of the operation.
- `UQ` - If the key already exists, an error will be returned.
- `RO` - Mark the key as read-only and create it if it doesn't exist.
- `NX` - If the key already exists, it won't be overwritten.
- `XX` - If the key doesn't exist, it won't be created.

`LEVEL` - Defines the level of permission.
- `R` (Restricted) - NO encryption, ACL validation
- `S` (Secret) - encryption, ACL validation
- Empty - NO encryption, NO ACL validation

`SERVER` - Defines server number to use.
- `> 0` - Use a specific server.
- `= 0` (default) - Automaticly saving to a less loaded server.

Example:
```go
SET key "value" UQ S 1
```

### GET - Gets the value from the storage.
```go
//      SERVER
GET key [ [0-9]+ ]
```

`SERVER` - Defines server number to use.
- `> 0` - Search on a specific server (speed: fast).
- `= 0` (default) - Deep search (speed: slow).

Example:
```go
GET key
```

### DEL - Deletes the key-value pair from the storage.
```go
//      SERVER
DEL key [ [0-9]+ ] 
```

`SERVER` - Defines server number to use.
- `> 0` - Search on a specific server (speed: fast).
- `= 0` (default) - Deep search (speed: slow). 

Example:
```go
DEL key
```