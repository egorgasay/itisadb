# Advanced object commands

### MARSHAL OBJECT - Displays the object as JSON.
```go
//                  SERVER
MARSHAL OBJECT name [ [0-9]+ ]
```

`SERVER` - Defines server number to use.
- `> 0` - Search on a specific server (speed: fast).
- `= 0` (default) - Deep search (speed: slow).

Example:
```go
MARSHAL OBJECT qwe
<- {
    "name": "qwe",
    "value": [
        {
            "name": "jkl",
            "value": "opi"
        }
    ]
}
```

### ATTACH - Attaches the src object to the dst.
```go
ATTACH dst src
```

Example:
```go
ATTACH obj53 obj56
<- OK
```