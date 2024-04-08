# Objects commands

### SETO

_Sets the value to the object._

```go
//                    MODE                  LEVEL         SERVER
SETO name key "value" [ RO | UQ | NX | XX ] [ D | R | L ] [ [0-9]+ ]
```

`MODE` - Defines the mode of the operation.
- `RO` - Mark the key as read-only and create it if it doesn't exist.
- `UQ` - If the key already exists, an error will be returned.
- `NX` - If the key already exists, it won't be overwritten.
- `XX` - If the key doesn't exist, it won't be created.

`LEVEL` - Defines the level of permission.
- `D` (Default) - NO encryption, NO ACL validation
- `R` (Restricted) - NO encryption, ACL validation
- `S` (Secret) - encryption, ACL validation

`SERVER` - Defines server number to use.
- `> 0` - Use a specific server.
- `= 0` (default) - Automaticly saving to a less loaded server.

Example:
```go
SETO obj13 key52 "value42" UQ S 1
```

### GETO

_Gets the value from the object._

```go
//            SERVER
GETO name key [ [0-9]+ ]
```

`SERVER` - Defines server number to use.
- `> 0` - Search on a specific server (speed: fast).
- `= 0` (default) - Deep search (speed: slow).

Example:
```go
GETO obj13 key52
```

### DELO

_Deletes the object key._

```go
//            SERVER
DELO name key [ [0-9]+ ]
```

`SERVER` - Defines server number to use.
- `> 0` - Search on a specific server (speed: fast).
- `= 0` (default) - Deep search (speed: slow).

Example:
```go
DELO obj13 key52
```