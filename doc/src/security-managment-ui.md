# Security management (Web UI)
 
Through the Web UI, you can manage users and their permissions.   
- Create new user
- Delete user
- Change user password
- Change user level

### Create new user

The user is created using the CREATE USER command

```bash
NEW USER <name> <password> <optional level>
```

Example:
```bash
NEW USER user1 123456
NEW USER user2 123456 R
NEW USER user3 123456 S
```

### Delete user

The user is deleted using the DELETE USER command

```bash
DELETE USER <name>
```

Example:
```bash
DELETE USER user1
```

### Change user password

The user password is changed using the CHANGE PASSWORD command

```bash
CHANGE PASSWORD <name> <password>
```

Example:
```bash
CHANGE PASSWORD user1 654321
```

### Change user level

The user level is changed using the CHANGE LEVEL command

```bash
CHANGE LEVEL <name> <level>
```

Example:
```bash
CHANGE LEVEL user1 R
```
