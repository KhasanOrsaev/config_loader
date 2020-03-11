Configuration loader
====
Configuration loader from postgres storage.


####Methods

* _/v1/config/apply_ - create or update configuration. 

Request:
```json
{
	"service_name": "test",
	"service_env": "test",
	"data": {
		"abc": "cde",
		"fgj": "tnm"
	}
}
```

* _/v1/config/get_ - get configuration

Request:
```json
{
	"service_name": "test",
	"service_env": "test"
}
```

####Response
On success:
```json
{
  "status": {
    "code":0,
    "error":null
  },
  "content":{   "id":3,
               "data":{
                  "abc":"cde",
                  "fgj":"tnm"
               
            },
               "user_id": <user_id>,
               "dt_create": <date create>,
               "dt_update": <date update>,
               "service_env":"test",
               "service_name":"test"
            }
}
```

On error:
```json
{   "status":{
      "code":<code>,
      "error": {
        "message":<message>,
        "stack_trace": <stack_trace>  
      }   
    },
    "content":null
}
```

####Authorization
Basic authorization.
