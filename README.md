# Aardvark Infocenter Task

## Table of Contents
* [Introduction](#Introduction)
* [Project Structure](#Project-Structure)
* [About the Project](#about-the-project)
* [How to Run](#how-to-run)
* [Author](#author)

# Introduction

Infocenter service challenge You're tasked to create a backend service using the Go programming language, which will allow its clients to achieve almost real-time communication between each other by sending messages. Concurrency is key when developing this service as it should be able to handle multiple clients sending or receiving messages at the same time.

# Project Structure
Infocenter/<br>
│── domain/          - Business logic & entities  
│── application/     - Use case layer  
│── infrastructure/  - HTTP handlers  
│── main.go          - Application entry point  
│── README.md        - Documentation  

# About the project

Infocenter is a lightweight, real-time message streaming application built using Go with Domain-Driven Design (DDD) principles. It allows clients to publish messages to topics and subscribe to real-time updates via Server-Sent Events (SSE).

# Features

Publish Messages – Clients can send messages to specific topics.
Subscribe to Messages – Clients receive real-time updates via SSE.
Message Caching – Stores up to the last 10 messages per topic.
Thread-Safe – Uses synchronization to handle concurrent requests.
Scalable & Modular – Designed with DDD principles for easy extension.

# How to run

Launching the application:
1. Clone the application using command line command "git clone (GitHubURL)"
2. Change the folder to Infocenter using command "cd Infocenter".
3. Launch the application by running command "go run main.go".

Testing:
1. To post a message use the command "curl -X POST -d "labas" http://localhost:8080/infocenter/adv".
2. To receive messages use the command "curl -N http://localhost:8080/infocenter/adv".

# Author

Erikas Pomeliaika.
