# Level

ItisaDB has a Level system for data access control.

You can manage security settings on the Web UI.

## Levels:
- D: Default
- R: Restricted
- S: Secret

### Default level

- Usually mentioned as `D`. 
- Level is used when no level is specified in the command.
- Requires only success authentication.

### Restricted level

- Usually mentioned as `R`. 
- Level is used when the level is specified in the command.
- Access to the values protected by this level is carried out after checking whether the user has the necessary access.

Example:
```go
SETO obj13 key52 "value42" R
CREATE USER user1 123456 R
```

### Secret level

- The secret level is usually mentioned as `S`. 
- The secret level is used when the level is specified in the command.
- Access to the values protected by this level is carried out after checking whether the user has the necessary access.
- It encrypts all stored values.

Example:
```go
SETO obj13 key52 "value42" S
CREATE USER user1 123456 S
```