# Go Events API Reference

Quick reference guide for all API endpoints in the Go Events application.

## Base URL

```
http://localhost:8080
```

## Authentication

Most endpoints require a JWT token in the Authorization header:

```
Authorization: <your-jwt-token>
```

## Endpoints Overview

| Method | Endpoint                  | Auth Required | Description                |
| ------ | ------------------------- | ------------- | -------------------------- |
| POST   | `/signup`                 | ❌            | Register a new user        |
| POST   | `/login`                  | ❌            | Login and get JWT token    |
| GET    | `/events`                 | ❌            | Get all events             |
| GET    | `/events/:id`             | ❌            | Get single event           |
| GET    | `/user/:id`               | ❌            | Get user by ID             |
| POST   | `/events`                 | ✅            | Create new event           |
| PUT    | `/events/:id`             | ✅            | Update event               |
| DELETE | `/events/:id`             | ✅            | Delete event               |
| POST   | `/events/:id/register`    | ✅            | Register for event         |
| DELETE | `/events/:id/cancel`      | ✅            | Cancel event registration  |
| GET    | `/notifications`          | ✅            | Get user notifications     |
| PUT    | `/notifications/:id/read` | ✅            | Mark notification as read  |
| POST   | `/notifications/trigger`  | ✅            | Trigger notification check |

---

## User Authentication

### Register User

**POST** `/signup`

Create a new user account.

```bash
curl -X POST http://localhost:8080/signup \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "securepassword"
  }'
```

**Response:**

```json
{
  "message": "User created"
}
```

### Login User

**POST** `/login`

Authenticate user and receive JWT token.

```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "securepassword"
  }'
```

**Response:**

```json
{
  "message": "login success",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### Get User by ID

**GET** `/user/:id`

Retrieve user information by ID.

```bash
curl http://localhost:8080/user/1
```

**Response:**

```json
{
  "ID": 1,
  "Email": "user@example.com",
  "Password": "$2a$14$hashedpassword..."
}
```

---

## Event Management

### Get All Events

**GET** `/events`

Retrieve all events (public endpoint).

```bash
curl http://localhost:8080/events
```

**Response:**

```json
[
  {
    "ID": 1,
    "Name": "Tech Conference 2024",
    "Description": "Annual technology conference",
    "Location": "Convention Center",
    "DateTime": "2024-12-20T09:00:00Z",
    "UserID": 1
  }
]
```

### Get Single Event

**GET** `/events/:id`

Retrieve a specific event by ID.

```bash
curl http://localhost:8080/events/1
```

**Response:**

```json
{
  "ID": 1,
  "Name": "Tech Conference 2024",
  "Description": "Annual technology conference",
  "Location": "Convention Center",
  "DateTime": "2024-12-20T09:00:00Z",
  "UserID": 1
}
```

### Create Event

**POST** `/events` 🔒

Create a new event (requires authentication).

```bash
curl -X POST http://localhost:8080/events \
  -H "Content-Type: application/json" \
  -H "Authorization: your-jwt-token" \
  -d '{
    "name": "Tech Conference 2024",
    "description": "Annual technology conference",
    "location": "Convention Center",
    "dateTime": "2024-12-20T09:00:00Z"
  }'
```

**Response:**

```json
{
  "message": "event created",
  "event": {
    "ID": 1,
    "Name": "Tech Conference 2024",
    "Description": "Annual technology conference",
    "Location": "Convention Center",
    "DateTime": "2024-12-20T09:00:00Z",
    "UserID": 1
  }
}
```

### Update Event

**PUT** `/events/:id` 🔒

Update an existing event (only by creator).

```bash
curl -X PUT http://localhost:8080/events/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: your-jwt-token" \
  -d '{
    "name": "Updated Conference Name",
    "description": "Updated description",
    "location": "New Location",
    "dateTime": "2024-12-21T09:00:00Z"
  }'
```

**Response:**

```json
{
  "message": "event updated"
}
```

### Delete Event

**DELETE** `/events/:id` 🔒

Delete an event (only by creator).

```bash
curl -X DELETE http://localhost:8080/events/1 \
  -H "Authorization: your-jwt-token"
```

**Response:**

```json
{
  "message": "event deleted"
}
```

---

## Event Registration

### Register for Event

**POST** `/events/:id/register` 🔒

Register the authenticated user for an event.

```bash
curl -X POST http://localhost:8080/events/1/register \
  -H "Authorization: your-jwt-token"
```

**Response:**

```json
{
  "message": "Event registration success"
}
```

### Cancel Event Registration

**DELETE** `/events/:id/cancel` 🔒

Cancel the user's registration for an event.

```bash
curl -X DELETE http://localhost:8080/events/1/cancel \
  -H "Authorization: your-jwt-token"
```

**Response:**

```json
{
  "message": "Event cancelled"
}
```

---

## Notifications

### Get User Notifications

**GET** `/notifications` 🔒

Retrieve all notifications for the authenticated user.

```bash
curl http://localhost:8080/notifications \
  -H "Authorization: your-jwt-token"
```

**Response:**

```json
[
  {
    "id": 1,
    "user_id": 123,
    "event_id": 456,
    "message": "Reminder: Your event 'Tech Conference' is in 2 hour(s) at 2:00 PM on Dec 20",
    "type": "upcoming_event",
    "is_read": false,
    "created_at": "2024-12-20T10:00:00Z"
  }
]
```

### Mark Notification as Read

**PUT** `/notifications/:id/read` 🔒

Mark a specific notification as read.

```bash
curl -X PUT http://localhost:8080/notifications/1/read \
  -H "Authorization: your-jwt-token"
```

**Response:**

```json
{
  "message": "Notification marked as read"
}
```

### Trigger Notification Check

**POST** `/notifications/trigger` 🔒

Manually trigger the notification system (for testing/development).

```bash
curl -X POST http://localhost:8080/notifications/trigger \
  -H "Authorization: your-jwt-token"
```

**Response:**

```json
{
  "message": "Notification check triggered successfully"
}
```

---

## Error Responses

The API returns consistent error responses:

### Authentication Errors

```json
{
  "message": "Not authorized"
}
```

### Validation Errors

```json
{
  "message": "Could not parse user data"
}
```

### Not Found Errors

```json
{
  "message": "Could not find event"
}
```

### Server Errors

```json
{
  "message": "Could not create event"
}
```

---

## HTTP Status Codes

| Code | Description                             |
| ---- | --------------------------------------- |
| 200  | OK - Request successful                 |
| 201  | Created - Resource created successfully |
| 400  | Bad Request - Invalid request data      |
| 401  | Unauthorized - Authentication required  |
| 403  | Forbidden - Access denied               |
| 404  | Not Found - Resource not found          |
| 500  | Internal Server Error - Server error    |

---

## Request/Response Headers

### Required Headers for Protected Endpoints

```
Authorization: <jwt-token>
Content-Type: application/json
```

### Response Headers

```
Content-Type: application/json
```

---

## Data Models

### User

```json
{
  "ID": 1,
  "Email": "user@example.com",
  "Password": "hashed_password"
}
```

### Event

```json
{
  "ID": 1,
  "Name": "Event Name",
  "Description": "Event description",
  "Location": "Event location",
  "DateTime": "2024-12-20T09:00:00Z",
  "UserID": 1
}
```

### Notification

```json
{
  "id": 1,
  "user_id": 123,
  "event_id": 456,
  "message": "Notification message",
  "type": "upcoming_event",
  "is_read": false,
  "created_at": "2024-12-20T10:00:00Z"
}
```

### Event Registration

```json
{
  "ID": 1,
  "EventID": 456,
  "UserID": 123
}
```
