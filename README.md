## Collabradoc

### Purpose

This application is for making a useable live collaboration editor. Mainly written for usage between small teams to edit private documents together. Was used between side project members initially and later shared to other friends.

### Current Features

- Live edit for documents via websockets, and using binary communication protocol for a smoother experience.
- Managing (save, edit) documents for each user.
- Full authentication flow for safely storing users and their documents.

### Planned Features

- Allowing users to invite and add others as collaborators seamlessly. [WIP]
- Tracks created documents and comments for each user seperately. [WIP]

### The Websocket Server

#### Why Websockets

I chose websockets over a custom tcp connection management because of simplicity and reliability, but still opting for performance over ease of development by choosing to use binary communication rather than utilizing JSON.

More details on this can be found in the Binary Communication Protocl section.

#### Design

The websocket setup sets up a single WS handler that servers each client that connects and joins the document edit functionality. The server then spins up a **goroutine** for each individual connected client, reading in their inidvidual payloads for handling.

The payloads are managed via a channel that recieves the payloads and are managed by another **goroutine** started at the root of the application. This goroutine handles each channel payload based on the pre-defined custom binary protocol that both the frontend nextjs application and this golang server follows.

### Custom Binary Communication Protocol
