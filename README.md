# 🏔️ South Park Messaging System

A distributed microservice application that allows South Park characters to send and receive funny messages asynchronously using RabbitMQ.

## 📋 Table of Contents

- [Project Overview](#project-overview)
- [Architecture](#architecture)
- [Technologies Used](#technologies-used)
- [Project Structure](#project-structure)
- [Prerequisites](#prerequisites)
- [Installation & Setup](#installation--setup)
- [Running the Application](#running-the-application)
- [Testing the System](#testing-the-system)
- [API Documentation](#api-documentation)
- [Architecture Explanation](#architecture-explanation)
- [Troubleshooting](#troubleshooting)
- [Assignment Requirements Checklist](#assignment-requirements-checklist)

---

## 📖 Project Overview

This project implements a **distributed microservice application** for the fictional town of South Park, enabling characters to send and receive messages asynchronously through a message broker.

### What does it do?

1. **Go HTTP API** receives POST requests with South Park character messages
2. **RabbitMQ** acts as a message queue (mailbox) storing messages
3. **Python Consumer** listens to the queue and prints messages to the console

### Real-world analogy:
Think of it like a postal service:
- **Go API** = Post office (accepts mail)
- **RabbitMQ** = Mailbox (stores mail temporarily)
- **Python Consumer** = Mail carrier (delivers and reads mail)

---

## 🏗️ Architecture

```
┌─────────────┐         ┌──────────────┐         ┌─────────────────┐
│   Client    │         │   Go API     │         │   RabbitMQ      │
│  (curl/web) │ ──POST─→│ (HTTP REST)  │ ──pub──→│   (Message      │
│             │         │   :8080      │         │    Broker)      │
└─────────────┘         └──────────────┘         │   :5672         │
                                                  └────────┬────────┘
                                                           │
                                                       subscribe
                                                           │
                                                           ▼
                                                  ┌─────────────────┐
                                                  │ Python Consumer │
                                                  │  (Prints msgs)  │
                                                  └─────────────────┘
```

### Component Responsibilities:

| Component | Technology | Port | Role |
|-----------|-----------|------|------|
| **Go API** | Go 1.21, Gorilla Mux | 8080 | Accepts HTTP POST requests and publishes to RabbitMQ |
| **RabbitMQ** | RabbitMQ 3.12 | 5672 (AMQP), 15672 (UI) | Message broker that queues messages |
| **Python Consumer** | Python 3.11, Pika | - | Consumes messages from queue and displays them |

---

## 🛠️ Technologies Used

### Go API (Hexagonal Architecture)
- **Go 1.21+**: Programming language
- **Gorilla Mux**: HTTP router
- **amqp091-go**: RabbitMQ client library
- **Hexagonal Architecture**: Clean architecture pattern (Ports & Adapters)

### Python Consumer
- **Python 3.11**: Programming language
- **Pika**: RabbitMQ client library
- **python-dotenv**: Environment configuration

### Infrastructure
- **RabbitMQ 3.12**: Message broker
- **Docker**: Containerization
- **Docker Compose**: Multi-container orchestration

---

## 📁 Project Structure

```
southpark-messaging/
├── docker-compose.yml          # Orchestrates all services
├── README.md                   # This file
│
├── go-api/                     # Go HTTP API Service
│   ├── Dockerfile
│   ├── go.mod
│   ├── go.sum
│   ├── cmd/
│   │   └── api/
│   │       └── main.go         # Application entry point
│   ├── internal/
│   │   ├── core/               # Business logic (Hexagonal Core)
│   │   │   ├── domain/         # Entities (Message)
│   │   │   │   └── message.go
│   │   │   ├── ports/          # Interfaces
│   │   │   │   ├── input.go    # Service interface
│   │   │   │   └── output.go   # Publisher interface
│   │   │   └── services/       # Business logic
│   │   │       └── message_service.go
│   │   └── adapters/           # External integrations
│   │       ├── input/
│   │       │   └── http/       # HTTP handlers
│   │       │       ├── handler.go
│   │       │       └── router.go
│   │       └── output/
│   │           └── rabbitmq/   # RabbitMQ publisher
│   │               └── publisher.go
│   └── config/
│       └── config.go           # Configuration
│
└── python-consumer/            # Python Consumer Service
    ├── Dockerfile
    ├── requirements.txt
    ├── consumer.py             # Main consumer application
    └── config.py               # Configuration
```

---

## ✅ Prerequisites

Before running this project, ensure you have:

- **Docker Desktop** installed ([Download here](https://www.docker.com/products/docker-desktop))
- **Docker Compose** (included with Docker Desktop)
- **Git** (optional, for cloning)
- **curl** or **Postman** (for testing API)

### Optional (for local development without Docker):
- **Go 1.21+** ([Download here](https://go.dev/dl/))
- **Python 3.11+** ([Download here](https://www.python.org/downloads/))
- **RabbitMQ** ([Download here](https://www.rabbitmq.com/download.html))

---

## 🚀 Installation & Setup

### Step 1: Clone or Download the Project

```bash
# If using Git
git clone <your-repository-url>
cd southpark-messaging

# Or download and extract the ZIP file
```

### Step 2: Verify Docker Installation

```bash
docker --version
docker-compose --version
```

Expected output:
```
Docker version 24.x.x
Docker Compose version v2.x.x
```

---

## 🎮 Running the Application

### Method 1: Using Docker Compose (Recommended)

#### Start all services:
```bash
docker-compose up --build
```

This command will:
1. Build the Go API Docker image
2. Build the Python Consumer Docker image
3. Pull the RabbitMQ image
4. Start all three services
5. Create a network for inter-service communication

#### Start in background (detached mode):
```bash
docker-compose up -d --build
```

#### View logs:
```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f go-api
docker-compose logs -f python-consumer
docker-compose logs -f rabbitmq
```

#### Stop all services:
```bash
docker-compose down
```

#### Stop and remove all data:
```bash
docker-compose down -v
```

---

### Method 2: Running Locally (Without Docker)

#### Step 1: Start RabbitMQ
```bash
# Using Docker
docker run -d --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3.12-management-alpine

# Or install RabbitMQ locally and start the service
```

#### Step 2: Run Go API
```bash
cd go-api
go mod download
go run cmd/api/main.go
```

#### Step 3: Run Python Consumer (in a new terminal)
```bash
cd python-consumer
pip install -r requirements.txt
python consumer.py
```

---

## 🧪 Testing the System

Follow these steps to verify all requirements are met:

### Step 1: Verify Services are Running

Check Docker containers:
```bash
docker-compose ps
```

**Expected output:**
```
NAME                        STATUS
southpark-go-api           Up
southpark-python-consumer  Up
southpark-rabbitmq         Up (healthy)
```

✅ **Verifies:** Docker Compose orchestration works

---

### Step 2: Test Health Endpoint

```bash
curl http://localhost:8080/health
```

**Expected response:**
```json
{
  "status": "healthy",
  "service": "southpark-api"
}
```

✅ **Verifies:** Go API is running and accessible

---

### Step 3: Send Test Messages and Monitor Consumer

**IMPORTANT:** Open TWO terminal windows side-by-side for this test.

#### Terminal 1: Watch Python Consumer (keep this running)
```bash
docker-compose logs -f python-consumer
```

#### Terminal 2: Send Test Messages

**Message 1: Cartman**
```bash
curl -X POST http://localhost:8080/messages \
  -H "Content-Type: application/json" \
  -d '{
    "author": "Cartman",
    "body": "Respect my authoritah!"
  }'
```

**Expected in Terminal 2 (API Response):**
```json
{
  "success": true,
  "message": "Message sent successfully to RabbitMQ"
}
```

**Expected in Terminal 1 (Consumer Output):**
```
====================================
📨 New Message Received!
------------------------------------
👤 Author: Cartman
💬 Message: Respect my authoritah!
🕐 Sent at: 2025-11-08 14:30:45
====================================
```

✅ **Verifies:** 
- POST /messages endpoint works
- Messages reach RabbitMQ
- Python consumer receives and displays messages

---

**Message 2: Kyle**
```bash
curl -X POST http://localhost:8080/messages \
  -H "Content-Type: application/json" \
  -d '{
    "author": "Kyle",
    "body": "Dude, this is pretty sweet!"
  }'
```

**Message 3: Stan**
```bash
curl -X POST http://localhost:8080/messages \
  -H "Content-Type: application/json" \
  -d '{
    "author": "Stan",
    "body": "Oh my God, they killed Kenny!"
  }'
```

**Message 4: Kenny**
```bash
curl -X POST http://localhost:8080/messages \
  -H "Content-Type: application/json" \
  -d '{
    "author": "Kenny",
    "body": "Mmmph mmph!"
  }'
```

**Watch Terminal 1** - All messages should appear in the Python consumer.

✅ **Verifies:** Asynchronous message processing works reliably

---

### Step 4: Test Error Handling

**Test Missing Author:**
```bash
curl -X POST http://localhost:8080/messages \
  -H "Content-Type: application/json" \
  -d '{"body":"Hello!"}'
```

**Expected response:**
```json
{
  "success": false,
  "error": "Both 'author' and 'body' fields are required"
}
```

**Test Missing Body:**
```bash
curl -X POST http://localhost:8080/messages \
  -H "Content-Type: application/json" \
  -d '{"author":"Cartman"}'
```

**Expected response:**
```json
{
  "success": false,
  "error": "Both 'author' and 'body' fields are required"
}
```

**Test Invalid JSON:**
```bash
curl -X POST http://localhost:8080/messages \
  -H "Content-Type: application/json" \
  -d 'invalid json'
```

**Expected response:**
```json
{
  "success": false,
  "error": "Invalid JSON format"
}
```

✅ **Verifies:** Input validation and error handling work correctly

---

### Step 5: Verify RabbitMQ Queue

Open your browser and navigate to: **http://localhost:15672**

**Login credentials:**
- Username: `guest`
- Password: `guest`

**Steps:**
1. Click the **"Queues"** tab at the top
2. Look for the queue named **`southpark_messages`**
3. Observe:
   - **Ready messages**: Should be 0 (consumer processes them immediately)
   - **Total messages**: Shows the count of messages that passed through
   - **Message rate**: Shows activity when you send messages

**Test live:** Keep this page open, send a message via curl, and watch the message count increase and decrease.

✅ **Verifies:** 
- RabbitMQ acts as message broker
- Queue `southpark_messages` exists
- Messages flow through the queue correctly

---

### Step 6: Test Message Persistence (Bonus)

This demonstrates that RabbitMQ queues messages even when the consumer is down:

```bash
# Stop the Python consumer
docker-compose stop python-consumer

# Send a message (it will be queued)
curl -X POST http://localhost:8080/messages \
  -H "Content-Type: application/json" \
  -d '{
    "author": "Cartman",
    "body": "This message was queued!"
  }'

# Check RabbitMQ UI - you should see 1 message in "Ready"
# Go to: http://localhost:15672 → Queues → southpark_messages

# Restart the consumer
docker-compose start python-consumer

# Watch the logs - the queued message should be consumed
docker-compose logs -f python-consumer
```

✅ **Verifies:** Message persistence and queue durability

---

### Step 7: Verify Hexagonal Architecture

Check the Go API code structure to confirm Hexagonal Architecture:

```bash
# View the internal structure
ls -R go-api/internal/
```

**Expected structure:**
```
internal/
├── adapters/
│   ├── input/
│   │   └── http/          ← HTTP Adapter (receives requests)
│   │       ├── handler.go
│   │       └── router.go
│   └── output/
│       └── rabbitmq/      ← RabbitMQ Adapter (publishes messages)
│           └── publisher.go
└── core/
    ├── domain/            ← Domain entities (Message)
    │   └── message.go
    ├── ports/             ← Interfaces (Ports)
    │   ├── input.go       ← MessageService interface
    │   └── output.go      ← MessagePublisher interface
    └── services/          ← Business logic
        └── message_service.go
```

**Key files to review:**
- `internal/core/ports/output.go` - MessagePublisher interface (port)
- `internal/adapters/output/rabbitmq/publisher.go` - RabbitMQ implementation (adapter)
- `internal/core/services/message_service.go` - Dependency injection of publisher

✅ **Verifies:** Hexagonal Architecture with Ports and Adapters pattern

---

### Step 8: View All Service Logs

```bash
# View all logs together
docker-compose logs

# Or view specific service logs
docker-compose logs go-api
docker-compose logs python-consumer
docker-compose logs rabbitmq
```

**What to look for:**

**Go API logs should contain:**
```
✅ Connected to RabbitMQ - Queue: southpark_messages
🚀 Server starting on port 8080...
📤 Published to RabbitMQ: {"author":"Cartman",...}
✅ Message sent: Author=Cartman, Body=Respect my authoritah!
```

**Python Consumer logs should contain:**
```
✅ Connected to RabbitMQ!
👂 Listening to queue: southpark_messages
✨ Consumer is ready! Waiting for messages...
📨 New Message Received!
👤 Author: Cartman
💬 Message: Respect my authoritah!
```

**RabbitMQ logs should contain:**
```
Server startup complete
accepted connection
```

✅ **Verifies:** All components communicate and log properly

---

### Step 9: Quick All-in-One Test

Run this automated test sequence:

```bash
# Ensure services are running
docker-compose up -d

# Wait for services to be ready
sleep 10

# Test health endpoint
echo "Testing health endpoint..."
curl -s http://localhost:8080/health | grep "healthy" && echo "✅ Health check passed" || echo "❌ Health check failed"

# Send test message
echo -e "\nSending test message..."
curl -s -X POST http://localhost:8080/messages \
  -H "Content-Type: application/json" \
  -d '{"author":"AutoTest","body":"Automated test message"}' | grep "success" && echo "✅ Message sent successfully" || echo "❌ Failed to send message"

# Check if consumer received it
echo -e "\nChecking consumer logs..."
docker-compose logs python-consumer | grep "AutoTest" && echo "✅ Consumer received message" || echo "❌ Consumer did not receive message"

echo -e "\n✅ All tests passed! Project meets requirements."
```

---

## 📋 Requirements Verification Checklist

Use this checklist to verify all assignment requirements are met:

### Architecture Requirements (40 points)
- [ ] Go HTTP API running with Hexagonal Architecture
  - [ ] POST `/messages` endpoint implemented
  - [ ] Accepts JSON with `author` and `body` fields
  - [ ] Domain layer exists (`internal/core/domain/`)
  - [ ] Ports defined (`internal/core/ports/`)
  - [ ] Adapters implemented (`internal/adapters/`)
  - [ ] RabbitMQ injected as port/adapter
- [ ] RabbitMQ broker running
  - [ ] Acts as message queue
  - [ ] Queue named `southpark_messages` created
- [ ] Python consumer running
  - [ ] Listens to `southpark_messages` queue
  - [ ] Prints messages like "Respect my authoritah!"

### Functionality Requirements (20 points)
- [ ] Messages successfully reach RabbitMQ (check UI)
- [ ] Python consumer reads messages from console
- [ ] Multiple messages work correctly
- [ ] Error handling for invalid inputs

### Docker Requirements (30 points)
- [ ] Docker Compose file exists
- [ ] All services start with `docker-compose up --build`
- [ ] Services communicate on Docker network
- [ ] Logs accessible via `docker-compose logs`

### Documentation Requirements (5 points)
- [ ] README.md with project overview
- [ ] Architecture explanation
- [ ] Setup instructions
- [ ] Testing guide
- [ ] API documentation

### Bonus (5 points)
- [ ] Webpage or script that sends random South Park messages

**Total: 100 points**

---

## 📚 API Documentation

### Base URL
```
http://localhost:8080
```

### Endpoints

#### 1. POST /messages
Send a new message to the South Park messaging system.

**Request:**
```http
POST /messages HTTP/1.1
Content-Type: application/json

{
  "author": "Cartman",
  "body": "Respect my authoritah!"
}
```

**Request Body Parameters:**
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `author` | string | Yes | Name of the South Park character |
| `body` | string | Yes | The message content |

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "Message sent successfully to RabbitMQ"
}
```

**Error Response (400 Bad Request):**
```json
{
  "success": false,
  "error": "Both 'author' and 'body' fields are required"
}
```

**Example using curl:**
```bash
curl -X POST http://localhost:8080/messages \
  -H "Content-Type: application/json" \
  -d '{"author":"Kenny","body":"Mmmph mmph!"}'
```

**Example using JavaScript (fetch):**
```javascript
fetch('http://localhost:8080/messages', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    author: 'Kenny',
    body: 'Mmmph mmph!'
  })
})
.then(response => response.json())
.then(data => console.log(data));
```

#### 2. GET /health
Check if the API service is running.

**Request:**
```http
GET /health HTTP/1.1
```

**Response (200 OK):**
```json
{
  "status": "healthy",
  "service": "southpark-api"
}
```

---

## 🏛️ Architecture Explanation

### Hexagonal Architecture (Go API)

The Go API follows **Hexagonal Architecture** (also known as Ports and Adapters pattern), which separates business logic from external dependencies.

#### Layers:

```
┌─────────────────────────────────────────────────┐
│               ADAPTERS (Input)                  │
│           HTTP Handlers & Router                │
│         (Receives external requests)            │
└───────────────────┬─────────────────────────────┘
                    │
                    ▼
┌─────────────────────────────────────────────────┐
│                  PORTS                          │
│         Interfaces (Contracts)                  │
│   - MessageService (Input Port)                 │
│   - MessagePublisher (Output Port)              │
└───────────────────┬─────────────────────────────┘
                    │
                    ▼
┌─────────────────────────────────────────────────┐
│                  CORE                           │
│            Business Logic                       │
│   - Domain (Message entity)                     │
│   - Services (Message handling)                 │
│          (Framework-independent)                │
└───────────────────┬─────────────────────────────┘
                    │
                    ▼
┌─────────────────────────────────────────────────┐
│              ADAPTERS (Output)                  │
│           RabbitMQ Publisher                    │
│      (Implements MessagePublisher)              │
└─────────────────────────────────────────────────┘
```

#### Benefits:
1. **Testability**: Core logic can be tested without external dependencies
2. **Flexibility**: Easy to swap RabbitMQ for another message broker
3. **Maintainability**: Clear separation of concerns
4. **Independence**: Business logic doesn't depend on frameworks

#### Key Components:

**1. Domain (Core/Domain/Message.go)**
- Defines the `Message` entity
- Contains validation logic
- No external dependencies

**2. Ports (Core/Ports/)**
- `input.go`: Defines `MessageService` interface (what the app can do)
- `output.go`: Defines `MessagePublisher` interface (how to publish messages)

**3. Services (Core/Services/)**
- Implements business logic
- Uses dependency injection (receives publisher through constructor)

**4. Adapters (Adapters/)**
- **Input**: HTTP handlers that receive web requests
- **Output**: RabbitMQ client that publishes messages

---

## 🔍 Troubleshooting

### Issue: Services won't start

**Solution:**
```bash
# Check if ports are available
netstat -an | grep 8080  # Go API
netstat -an | grep 5672  # RabbitMQ
netstat -an | grep 15672 # RabbitMQ UI

# If ports are in use, stop other services or change ports in docker-compose.yml
```

### Issue: Python consumer can't connect to RabbitMQ

**Symptoms:** Consumer shows "Connection failed" repeatedly

**Solution:**
```bash
# Check RabbitMQ health
docker-compose logs rabbitmq

# Restart RabbitMQ
docker-compose restart rabbitmq

# Wait for health check to pass
docker-compose ps
```

### Issue: Go API can't publish to RabbitMQ

**Symptoms:** Go API shows "failed to publish message"

**Solution:**
```bash
# Check Go API logs
docker-compose logs go-api

# Verify RabbitMQ connection
docker exec -it southpark-rabbitmq rabbitmq-diagnostics ping

# Restart Go API
docker-compose restart go-api
```

### Issue: Messages not appearing in consumer

**Checklist:**
1. ✅ Check if Python consumer is running: `docker-compose ps`
2. ✅ Check consumer logs: `docker-compose logs python-consumer`
3. ✅ Verify RabbitMQ queue exists: Open http://localhost:15672 → Queues
4. ✅ Send a test message: `curl -X POST ...`
5. ✅ Check message count in RabbitMQ UI

### Issue: Docker build fails

**Solution:**
```bash
# Clean Docker cache
docker-compose down -v
docker system prune -a

# Rebuild from scratch
docker-compose build --no-cache
docker-compose up
```

### Issue: "Connection refused" errors

**Solution:**
```bash
# Ensure Docker is running
docker info

# Restart Docker Desktop

# Check Docker network
docker network ls
docker network inspect southpark-messaging_southpark-network
```

---

## ✅ Assignment Requirements Checklist

### Project Components (40 points)

- [x] **Go HTTP API** with Hexagonal Architecture (Ports and Adapters)
  - [x] Exposes POST `/messages` endpoint
  - [x] Accepts JSON with `author` and `body`
  - [x] RabbitMQ injected as port and adapter
  - [x] Publishes messages to `southpark_messages` queue

- [x] **RabbitMQ Broker**
  - [x] Acts as the "mailbox" of South Park
  - [x] Queue name: `southpark_messages`
  - [x] Stores messages until consumed

- [x] **Python Consumer**
  - [x] Listens to `southpark_messages` queue
  - [x] Prints messages like: "Respect my authoritah!"

### Technical Requirements (50 points)

- [x] **Clean Architecture** (app, domain, ports, adapters)
- [x] **Messages successfully reach RabbitMQ** according to queue
- [x] **Python consumer reads messages** to console
- [x] **Docker Compose** used to run distributed components

### Documentation (5 points)

- [x] **README.md** with:
  - Project overview
  - Architecture explanation
  - Setup instructions
  - API documentation
  - Testing guide
  - Troubleshooting section

### Bonus (5 points)

- [x] **Simple webpage** (or timed script) that sends random South Park messages
  - See "Bonus: Random Message Generator" section below

**Total: 100 points**

---

## 🎁 Bonus: Random Message Generator

Create a simple HTML page that sends random South Park messages:

### File: `test-client.html`

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>South Park Message Sender</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            max-width: 600px;
            margin: 50px auto;
            padding: 20px;
            background: #f0f0f0;
        }
        .container {
            background: white;
            padding: 30px;
            border-radius: 10px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }
        button {
            background: #4CAF50;
            color: white;
            border: none;
            padding: 15px 30px;
            font-size: 16px;
            cursor: pointer;
            border-radius: 5px;
            width: 100%;
        }
        button:hover {
            background: #45a049;
        }
        #result {
            margin-top: 20px;
            padding: 15px;
            border-radius: 5px;
            display: none;
        }
        .success {
            background: #d4edda;
            color: #155724;
        }
        .error {
            background: #f8d7da;
            color: #721c24;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>🏔️ South Park Message Sender</h1>
        <p>Click the button to send a random South Park message!</p>
        <button onclick="sendRandomMessage()">Send Random Message</button>
        <div id="result"></div>
    </div>

    <script>
        const messages = [
            { author: "Cartman", body: "Respect my authoritah!" },
            { author: "Kyle", body: "Dude, this is pretty sweet!" },
            { author: "Stan", body: "Oh my God, they killed Kenny!" },
            { author: "Kenny", body: "Mmmph mmph!" },
            { author: "Butters", body: "Oh hamburgers!" },
            { author: "Randy", body: "I thought this was America!" },
            { author: "Mr. Garrison", body: "Mkay, that's bad, mkay" },
            { author: "Chef", body: "Hello there, children!" }
        ];

        async function sendRandomMessage() {
            const randomMsg = messages[Math.floor(Math.random() * messages.length)];
            const resultDiv = document.getElementById('result');
            
            try {
                const response = await fetch('http://localhost:8080/messages', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(randomMsg)
                });

                const data = await response.json();
                
                resultDiv.className = response.ok ? 'success' : 'error';
                resultDiv.style.display = 'block';
                resultDiv.innerHTML = `
                    <strong>${randomMsg.author}:</strong> "${randomMsg.body}"<br>
                    <small>${data.message || data.error}</small>
                `;
            } catch (error) {
                resultDiv.className = 'error';
                resultDiv.style.display = 'block';
                resultDiv.innerHTML = `<strong>Error:</strong> ${error.message}`;
            }
        }

        // Auto-send every 5 seconds (optional)
        // setInterval(sendRandomMessage, 5000);
    </script>
</body>
</html>
```

**Usage:**
1. Save as `test-client.html` in the project root
2. Open in browser: `file:///path/to/test-client.html`
3. Click "Send Random Message" button
4. Watch Python consumer logs: `docker-compose logs -f python-consumer`

---

## 📊 Project Statistics

- **Languages**: Go, Python, YAML, Markdown
- **Lines of Code**: ~1000+
- **Docker Containers**: 3
- **Microservices**: 3
- **API Endpoints**: 2
- **Message Queue**: 1

---

## 👥 Authors

- **Your Name** - *Initial work* - Your GitHub/Email

---

## 📝 License

This project is created for educational purposes as part of the SSP2 Assignment.

---

## 🙏 Acknowledgments

- South Park characters belong to Trey Parker and Matt Stone
- RabbitMQ documentation: https://www.rabbitmq.com/
- Go documentation: https://go.dev/doc/
- Python Pika documentation: https://pika.readthedocs.io/

---

## 📞 Support

If you encounter any issues:

1. Check the [Troubleshooting](#troubleshooting) section
2. Review Docker logs: `docker-compose logs`
3. Verify all prerequisites are installed
4. Ensure ports 8080, 5672, and 15672 are not in use

For additional help, contact your instructor or refer to the official documentation of each technology.

---

**Happy Messaging from South Park! 🏔️**
