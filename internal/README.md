# internals

this folder contains the type definitions, interfaces and factories that implement all
of the business logic needed by the API.

the code written here should be private code that should not be used by outside packages, 
with the only exception being `/handlers` that exposes the function handlers to 
external requests.

all API endpoints are defined inside of `handlers/api` and all other code pretty much 
exists just to implement the logic that these endpoints need.

here is a brief description of all folders and what they are for:

## handlers

## middleware

## models

## services

## tools
