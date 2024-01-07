# Security management (Config)

Through the configuration file, you can manage security only after restarting the server.
- Enabling or disabling mandatory security
- Enabling or disabling mandatory encryption
- Enabling or disabling demo user

### Security section:
```toml
# Security settings.
[Security]
On = true

# Mandatory authorization.
# If false, authentication is not required for keys and objects that has Default level.
MandatoryAuthorization = true

# Mandatory encryption.
# If false, encryption is not mandatory for keys and objects that has Secret level.
MandatoryEncryption = true
```