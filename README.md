## Collabradoc

### Purpose

This application is for making a useable live collaboration editor for my teams that are working on side projects.

### Current Features

- Live edit for documents via websockets, and using binary communication protocol for a smoother experience.
- Managing (save, edit) documents for each user.
- Full authentication flow for safely storing users and their documents.

### Planned Features

- Allowing users to invite and add others as collaborators seamlessly.
- Add functionality for separating

### The Websocket Server

#### Why Websockets

I chose websockets over a custom tcp connection management because of simplicity and reliability, but we still opted for performance over ease of development by opting out of JSON for communication.

More details on this can be found in the Binary Communication Protocl section.

#### Design

The websocket setup is simple, setting up a single WS handler that servers each client that connects to the document edit functionality. The server then spins up a **goroutine** for each individual connected client, reading in their inidvidual payloads for handling.

The payloads are managed via a channel that recieves the payloads and are managed by another **goroutine** started at the root of the application. This goroutine will handle each channel message based on the pre-defined custom binary protocol.

#### Custom Binary Communication Protocol
