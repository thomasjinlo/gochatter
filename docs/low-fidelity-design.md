# Gochatter Design
Version: 0.1
Description: Low fidelity design of Gochatter, a chat service

## Introduction
This document describes the technical design for Gochatter service. Gochatter is
a chat service with focus on being minimal. Unlike other chat services, Gochatter
should be runnable on any machine.

This document will cover the high-level architecture design, design components, and the
functional and non-functional requirements. Alternate design considerations will also
be made including the tradeoffs which led to the final decision.

The target audience are software engineers, software architects, and anyone in a
technical role. Non-technical readers should be able to digest the high-level
architecture design and functional/non-functional requirements.

## Problem Statement
Modern day chat services typically require either a desktop application or a browser.
It's very common for personal machines to have both of these applications, but neither
are accessible on remote machines. Gochatter solves this problem by providing a command line
application that can run on any remote machine or personal lightweight machines with
minimal dependencies.

## Functional Requirements
1. Users can view available channels
2. Users can join available channels
3. Users can send messages to other users in the channel
4. Users can receive messages sent from other users in the channel

## Non-Functional Requirements
1. Messages should be encrypted using TLS 1.3
2. JWT Authorization required for client requests

## High-Level Architecture
### Client
The client is a terminal ui consisting of basic chat components such as text input,
message display view, and user and channel list views. The client will be written in go
and will use rivo/tview or bubbletea terminal ui libraries. Net/http, gorilla/websocket,
and crypto/tls will be used to send HTTPS requests and maintain TCP connections.

### Server
The server will be hosted on a single machine running an HTTP server listening on a single
port. The server will use TLS encryption and verify requests have proper authorization
using third-party authorization provider. The server will be written in Go and utilize
available data structures to maintain and manage TCP connections and message broadcasting.

### Communication Protocol
* HTTPS will be used to communication CRUD requests from client to server
* Client request for joining a channel will create a Websocket connection

## Component Design
### Client - API
Responsible for handling HTTPS requests and responses to create and list channels.

### Client - Socket
Responsible for sending and receiving messages to and from the server. There should
be a 1:1 mapping between number of joined channels and active sockets.

### Client - TUI
Responsible for rendering the terminal ui. Should expose an interface to display
messages received from server.

### Server - Sockets
Responsible for sending and receiving messages to and from clients.

### Server - HTTP Server
Responsible for listening and accepting connections. Should also delegate requests
to specific handlers.

### Server - Handlers
Responsible for handling client requests such as create, list, and joining channels.

### Server - Channel Model
Responsible for handling socket management and broadcasting messages sent from a client
in one channel to other clients in the same channel.
