# Transaction Logger

Protection against data loss in case of various hardware problems is achieved by using Transaction Logger. Each operation in the background is written to disk and performed again when the server is turned on after a failure (in other cases this does not happen).

```table
operation | key | value
```  

!!! DO NOT USE temporary directories for tlog_dir !!!