# microstore

Microstore is a universal agentic tool store that focuses on being extremely usable and
reliable. We sell tools for AI agents hosted on f4rmhouse.com. Agents
use tools and pay by usage (pay/request) 

## getting started
If you want to get a deep understanding or what f4rmhouse is you can visit our 
documentations page

If you want to run microstore locally just clone the repo 

```
git clone git@github.com:fizlip/microstore.git
```

Then cd into microstore and run
```
go run ./cmd/api/main.go
```

This will start the microstore on port 8000.

## structure
This is the directory structure we use for the project
```
.
├── api
│   ├── mocks
│   ├── models
│   └── services
├── cmd
│   └── api
└── internal
    ├── handlers
    │   ├── review_handlers
    │   └── vendor_handlers
    ├── middleware
    └── tools
```

### api
The ´api´ directory contains the code for all our entities. The main entities currently 
are: 
    * Analytics
    * Product
    * User
    * Review
    * PaymentExecutor

All of these are located in separate files. There also is an `api.go` file that
contains some utility structs to handle responses and requests.

### cmd

### internal

## testing

## current sprint
| feature | status |
| -------- | ------- |
| testing |  ❌ |
| documentation | ✅|

legend: 

❌ -- not finished

✅ -- finished

View previous sprints and more detailed descriptions in the [projects](https://github.com/fizlip/microstore/projects?query=is%3Aopen) tab.

## license
