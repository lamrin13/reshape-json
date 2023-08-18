# reshape-json
Takes JSON string and mapping config, returns the reshaped JSON according to config.

## Syntax of mapping config
If the source JSON is,
```
{
    "name": "John",
    "email": "john.doe@abc.com",
    "phone": "6471663421",
    "street": "50 Dummy Street",
    "unit": 12
}
```
And if we want the output JSON to be,
```
{
    "address": "12-50 Dummy Street",
    "user": {
        "emailAddress": "john.doe@abc.com",
        "firstName": "John",
        "phoneNumber": "6471663421"
    }
}
```
The mapping config for this operation would be,
```
{
    "user.firstName":    "name",
    "address":           "unit+street/-",
    "user.emailAddress": "email",
    "user.phoneNumber":  "phone",
}
```

