# microstore

Microstore is an alternative app store that lets users upload Docker files that 
will act like tools for AI agents hosted on f4rmhouse.com. Agents
will use the tools and pay by usage (pay/request) 

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
│   └── mocks
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

## license
