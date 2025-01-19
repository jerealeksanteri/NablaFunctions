# Nabla Functions

## Introduction
This is a Golang HTTP Server project, which mocks AWS Lambda Functions. It uses docker images as functions. Docker containers are launched when the function is called.
Currently it supports only Golang and Python projects. Heavily influenced by https://www.youtube.com/@helloWorldGolang

## Usage
The server has two endpoint load and execute.

### ZipFile
The directory should have its "main" -file. The program doesnt recognize start file without it.

### Endpoints

#### Load
( Makefile has its operation )
http://localhost:8080/api/load

The zip file should be in the "code" multipart form.

The endpoint will output Docker Image ID and the ***Function ID***

#### Execute
( Makefile has its operation )
http://localhost:8080/api/execute?functionId=$(functionId)

***Function ID*** must be passed as a parameter.

The endpoint will output the functions output.

## Lane Diagram of usage
![NablaFunctions](https://github.com/user-attachments/assets/c054bd96-bd50-4e1c-8256-37b69be22967)


