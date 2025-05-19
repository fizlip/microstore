# microstore

[![Version](https://img.shields.io/badge/version-1.0.0-blue.svg)](https://github.com/f4rmhouse/f4rmhouse)
[![Commits/Month](https://img.shields.io/github/commit-activity/m/f4rmhouse/f4rmhouse)](https://github.com/f4rmhouse/f4rmhouse/commits/main)
[![Issues Closed](https://img.shields.io/github/issues-closed/f4rmhouse/f4rmhouse)](https://github.com/f4rmhouse/f4rmhouse/issues?q=is%3Aissue+is%3Aclosed)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](https://github.com/f4rmhouse/f4rmhouse/pulls)

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

Each folder has it's own README that describes what it does

### api
JSON schema files, protocol definition files.

### cmd
Contains main.go file and also used for binaries.

### internal
Contains private code that should not be used by outside packages. Subdirectories
for models, services and handlers.

## testing

## current sprint
| feature | status |
| -------- | ------- |
| testing |  ❌ |
| documentation | ✅|
| auth | ✅|
| docker deployments | ✅ |

legend: 

❌ -- not finished

✅ -- finished

View previous sprints and more detailed descriptions in the [projects](https://github.com/fizlip/microstore/projects?query=is%3Aopen) tab.

## license
