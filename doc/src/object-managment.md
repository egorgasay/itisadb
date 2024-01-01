# Data manipulation

### NEW OBJECT - Creates an object with the specified name.
```go
//              LEVEL     SERVER
NEW OBJECT name [ R | S ] [ [0-9]+ ] 
```

`LEVEL` - Defines the level of permission.
- `R` (Restricted) - NO encryption, ACL validation
- `S` (Secret) - encryption, ACL validation
- Empty - NO encryption, NO ACL validation

`SERVER` - Defines server number to use.
- Automaticly saving to a less loaded server by default.

Example:
```go
NEW OBJECT obj R 6
```

### DELETE OBJECT - Deletes the object with the specified name.
```go
//                 SERVER
DELETE OBJECT name [ [0-9]+ ] 
```

`SERVER` - Defines server number to use.
- Automaticly saving to a less loaded server by default.

Example:
```go
DELETE OBJECT obj
```