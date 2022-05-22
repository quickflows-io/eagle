## Error code design

> Reference Sina Open Platform [Error code](http://open.weibo.com/wiki/Error_code) the design of

#### Error return value format

```json
{
  "code": 10002,
  "message": "Error occurred while binding the request body to the struct."
}
```

#### Error code description

| 1 | 00 | 02 |
| :------ | :------ | :------ |
| service level error（1 is a system level error） | service module code | specific error code |

- Service level error: 1 is a system-level error; 2 is a common error, usually caused by illegal user operations
- The service module is two digits: the service module of a large system usually does not exceed two digits. If it exceeds, it means that the system should be split.
- The error code is two digits: to prevent a module from customizing too many error codes, it is not easy to maintain later
- `code = 0` indicates correct return, `code > 0` indicates error return
- Errors usually include system-level error codes and service-level error codes
- It is recommended to categorize errors by service module in the code
- Error codes are all numbers >= 0
- In this project, the HTTP Code is fixed as http.StatusOK, and the error code is represented by code.