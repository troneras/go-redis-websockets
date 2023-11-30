# gorews

## Overview
goreWS is a scalable server application integrating Redis and WebSockets using Go. It's designed to handle real-time data communication efficiently and robustly.

## Features
- Modular design with distinct components for configuration, data handling, WebSocket communication, and Redis interactions.
- Efficient real-time communication between Redis channels and WebSocket clients.
- Scalable architecture suitable for high-concurrency environments.

## Getting Started

### Prerequisites
- Go 1.21 or higher
- Redis server

### Installation
Clone the repository:
```bash
git clone https://github.com/troneras/goreWS.git
cd goreWS
```

### Running the application
```
# Start the server
go run main.go
```

### Modules
#### Config
Manages application configuration settings.

#### Data
Handles data modeling and interactions.

#### WebSockets
Manages WebSocket connections and communication.

#### Redis-Client
Handles interactions with the Redis server, including subscribing to and publishing in channels.

#### License
This project is licensed under the MIT License - see the LICENSE.md file for details.