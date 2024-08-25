## Collabradoc

### Purpose

This application is for making a useable live collaboration editor. Mainly written for usage between small teams to edit private documents together. Was used between side project members initially and later shared to other friends.

### Current Features

**v0.3.0**

- Live edit for documents via websockets, and using binary communication protocol for a smoother experience.
- Managing (save, edit) documents for each user.
- Full authentication flow for safely storing users and their documents.

**Latest in recent version**

- Allowing users to invite and add others as collaborators seamlessly.
- Tracks created documents and comments for each user seperately.

### Upcoming Features

**v0.3.1**

- Bugfixes and quality of life improvements like updating response values.

**v0.4.0**

- Add email system for collaboration.
- Add the initial stages of advanced live-doc editing features.

### The Websocket Server

#### Why Websockets

I chose websockets over a custom tcp connection management because of simplicity and reliability, but still opting for performance over ease of development by choosing to use binary communication rather than utilizing JSON.

More details on this can be found in the Binary Communication Protocl section.

#### Design

The websocket setup sets up a single WS handler that servers each client that connects and joins the document edit functionality. The server then spins up a **goroutine** for each individual connected client, reading in their inidvidual payloads for handling.

The payloads are managed via a channel that recieves the payloads and are managed by another **goroutine** started at the root of the application. This goroutine handles each channel payload based on the pre-defined custom binary protocol that both the frontend nextjs application and this golang server follows.

### Custom Binary Communication Protocol

_Details coming soon._

### Live Editing Documents

As detailed in the websocket server section live editing documents was created via websockets connections and concurrently managing connected clients.

A **live session** is created once per document and is used to initialize and authorize a document to be edited live.
Collaborators are added to this live session when invited via email. This live session essentially creates a unique
instance between these collaborators and they have their own goroutine in the websocket server handling their
requests.
