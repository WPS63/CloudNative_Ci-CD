package main

import (
    "log"
    "github.com/go-openapi/loads"
    "github.com/go-openapi/runtime/middleware"
    "labs/lab6/rest-api-swagger/pkg/swagger/server/restapi"
    "labs/lab6/rest-api-swagger/pkg/swagger/server/restapi/operations"
)
//Health route returns OK
func Health(operations.CheckHealthParams) middleware.Responder {
    return operations.NewCheckHealthOK().WithPayload("OK\n")
}

//GetHelloUser returns Hello + your name
func GetHelloUser(user operations.GetHelloUserParams) middleware.Responder {
    return operations.NewGetHelloUserOK().WithPayload("Hello " + user.User + "!")
}
func main() {

    // Initialize Swagger
    swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "") 
    if err != nil {
        log.Fatalln(err)
    }
    api := operations.NewHelloAPIAPI(swaggerSpec)
    server := restapi.NewServer(api)
// Implement the CheckHealth handler
    api.CheckHealthHandler = operations.CheckHealthHandlerFunc(Health)

// Implement the GetHelloUser handler
    api.GetHelloUserHandler = operations.GetHelloUserHandlerFunc(GetHelloUser)

    defer func() {
        if err := server.Shutdown(); err != nil {
            // error handle
            log.Fatalln(err)
        }
    }() 

    server.Port = 8080

    // Start server which listening
    if err := server.Serve(); err != nil {
        log.Fatalln(err)
    }
}
