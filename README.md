# Go Events API

A comprehensive event management system built with Go, featuring event creation, user registration, authentication, and automated notification system.

## 🚀 Features

### Core Features

- **Event Management**: Create, read, update, and delete events
- **User Authentication**: JWT-based authentication with secure password hashing
- **Event Registration**: Users can register for and cancel event registrations
- **Automated Notifications**: Background job system for upcoming event reminders
- **Database Integration**: MySQL database with automatic table creation

### Key Capabilities

- ✅ User signup and login with JWT authentication
- ✅ CRUD operations for events
- ✅ Event registration and cancellation
- ✅ Automatic notification system for upcoming events
- ✅ Secure password hashing with bcrypt
- ✅ Authentication middleware for protected routes
- ✅ RESTful API design

## 🛠️ Technology Stack

- **Language**: Go 1.22.2
- **Web Framework**: Gin
- **Database**: MySQL
- **Authentication**: JWT (JSON Web Tokens)
- **Password Hashing**: bcrypt
- **Database Driver**: go-sql-driver/mysql

## 📋 Prerequisites

- Go 1.22.2 or higher
- MySQL 5.7 or higher
- Git

## 🔧 Installation & Setup

### 1. Clone the Repository

```bash
git clone <repository-url>
cd go-events
```

### 2. Install Dependencies

```bash
go mod tidy
```

### 3. Database Setup

Ensure MySQL is running and create a database:

```sql
CREATE DATABASE go_events;
```

### 4. Configuration

The application uses the following default database configuration:

- **Host**: 127.0.0.1
- **Port**: 3306
- **Username**: root
- **Password**: (empty)
- **Database**: go_events

To modify these settings, edit the DSN in `db/db.go`:

```go
dsn := "root:@tcp(127.0.0.1:3306)/go_events?parseTime=true"
```

### 5. Run the Application

```bash
go run main.go
```

The server will start on `http://localhost:8080`

## 📊 Database Schema

The application automatically creates the following tables:

### Users Table

```sql
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL
);
```

### Events Table

```sql
CREATE TABLE events (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    location VARCHAR(255) NOT NULL,
    dateTime DATETIME NOT NULL,
    user_id INT,
    FOREIGN KEY (user_id) REFERENCES users(id)
);
```

### Events Registry Table

```sql
CREATE TABLE events_registry (
    id INT AUTO_INCREMENT PRIMARY KEY,
    event_id INT,
    user_id INT,
    FOREIGN KEY (event_id) REFERENCES events(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);
```

### Notifications Table

```sql
CREATE TABLE notifications (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    event_id INT NOT NULL,
    message TEXT NOT NULL,
    type VARCHAR(50) NOT NULL,
    is_read BOOLEAN DEFAULT FALSE,
    created_at DATETIME NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (event_id) REFERENCES events(id)
);
```

## 🔗 API Documentation

### Authentication

All protected endpoints require a JWT token in the Authorization header:

```
Authorization: <your-jwt-token>
```

### Public Endpoints

#### User Registration

```http
POST /signup
Content-Type: application/json

{
    "email": "user@example.com",
    "password": "securepassword"
}
```

#### User Login

```http
POST /login
Content-Type: application/json

{
    "email": "user@example.com",
    "password": "securepassword"
}

Response:
{
    "message": "login success",
    "token": "eyJhbGciOiJIUzI1NiIs..."
}
```

#### Get All Events

```http
GET /events

Response:
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

#### Get Single Event

```http
GET /events/:id

Response:
{
    "ID": 1,
    "Name": "Tech Conference 2024",
    "Description": "Annual technology conference",
    "Location": "Convention Center",
    "DateTime": "2024-12-20T09:00:00Z",
    "UserID": 1
}
```

#### Get User by ID

```http
GET /user/:id

Response:
{
    "ID": 1,
    "Email": "user@example.com",
    "Password": "$2a$14$..." // hashed password
}
```

### Protected Endpoints (Require Authentication)

#### Create Event

```http
POST /events
Authorization: <jwt-token>
Content-Type: application/json

{
    "name": "Tech Conference 2024",
    "description": "Annual technology conference",
    "location": "Convention Center",
    "dateTime": "2024-12-20T09:00:00Z"
}

Response:
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

#### Update Event

```http
PUT /events/:id
Authorization: <jwt-token>
Content-Type: application/json

{
    "name": "Updated Conference Name",
    "description": "Updated description",
    "location": "New Location",
    "dateTime": "2024-12-21T09:00:00Z"
}

Response:
{
    "message": "event updated"
}
```

#### Delete Event

```http
DELETE /events/:id
Authorization: <jwt-token>

Response:
{
    "message": "event deleted"
}
```

#### Register for Event

```http
POST /events/:id/register
Authorization: <jwt-token>

Response:
{
    "message": "Event registration success"
}
```

#### Cancel Event Registration

```http
DELETE /events/:id/cancel
Authorization: <jwt-token>

Response:
{
    "message": "Event cancelled"
}
```

### Notification Endpoints

#### Get User Notifications

```http
GET /notifications
Authorization: <jwt-token>

Response:
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

#### Mark Notification as Read

```http
PUT /notifications/:id/read
Authorization: <jwt-token>

Response:
{
    "message": "Notification marked as read"
}
```

#### Trigger Notification Check (Development/Testing)

```http
POST /notifications/trigger
Authorization: <jwt-token>

Response:
{
    "message": "Notification check triggered successfully"
}
```

## 🔔 Notification System

The application includes an automated notification system with the following features:

### How It Works

1. **Background Job**: Runs every hour automatically
2. **Event Detection**: Finds events happening within the next 24 hours
3. **User Targeting**: Notifies only users registered for the event
4. **Smart Notifications**: Prevents duplicate notifications (one per day per event per user)
5. **Contextual Messages**: Generates different messages based on event timing

### Message Types

- **Within 1 hour**: "Reminder: Your event 'EventName' is starting soon at 3:04 PM!"
- **Within 24 hours**: "Reminder: Your event 'EventName' is in X hour(s) at 3:04 PM on Jan 2"
- **Beyond 24 hours**: "Reminder: You have an upcoming event 'EventName' on January 2, 2006 at 3:04 PM"

For detailed notification system documentation, see [NOTIFICATIONS.md](NOTIFICATIONS.md)

## 🏗️ Project Structure

```
go-events/
├── db/
│   └── db.go                 # Database connection and table creation
├── jobs/
│   └── notification_job.go   # Background notification service
├── middlewares/
│   └── auth.go              # JWT authentication middleware
├── models/
│   ├── event.go             # Event model and database operations
│   ├── event-register.go    # Event registration model
│   ├── notification.go      # Notification model
│   └── user.go              # User model and authentication
├── routes/
│   ├── events.go            # Event-related routes
│   ├── notifications.go     # Notification routes
│   ├── registration.go      # Event registration routes
│   ├── routes.go            # Main route registration
│   └── users.go             # User authentication routes
├── utils/
│   ├── hash.go              # Password hashing utilities
│   └── jwt.go               # JWT token utilities
├── main.go                  # Application entry point
├── go.mod                   # Go module dependencies
├── go.sum                   # Dependency checksums
├── README.md                # This documentation
├── NOTIFICATIONS.md         # Notification system documentation
└── TESTING.md               # Comprehensive testing guide
```

## 🧪 Testing

The project includes a comprehensive test suite with over **50 unit tests** covering all components:

### Quick Test Commands

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test ./... -cover

# Run specific package tests
go test ./utils/... -v
go test ./models/... -v
go test ./middlewares/... -v
go test ./jobs/... -v
```

### Test Coverage

- ✅ **Utils**: Password hashing, JWT tokens (100% coverage)
- ✅ **Models**: All CRUD operations with database mocking (100% coverage)
- ✅ **Middleware**: Authentication flow testing (100% coverage)
- ✅ **Jobs**: Background notification service (100% coverage)

For detailed testing information, see [TESTING.md](TESTING.md)

### Manual Testing Steps

1. **User Registration and Authentication**

   ```bash
   # Register a new user
   curl -X POST http://localhost:8080/signup \
     -H "Content-Type: application/json" \
     -d '{"email":"test@example.com","password":"password123"}'

   # Login to get JWT token
   curl -X POST http://localhost:8080/login \
     -H "Content-Type: application/json" \
     -d '{"email":"test@example.com","password":"password123"}'
   ```

2. **Event Management**

   ```bash
   # Create an event (replace TOKEN with actual JWT)
   curl -X POST http://localhost:8080/events \
     -H "Content-Type: application/json" \
     -H "Authorization: TOKEN" \
     -d '{
       "name":"Test Event",
       "description":"A test event",
       "location":"Test Location",
       "dateTime":"2024-12-25T10:00:00Z"
     }'

   # Get all events
   curl http://localhost:8080/events
   ```

3. **Event Registration**

   ```bash
   # Register for an event (replace :id with actual event ID)
   curl -X POST http://localhost:8080/events/:id/register \
     -H "Authorization: TOKEN"
   ```

4. **Notification Testing**

   ```bash
   # Create an event happening within 24 hours
   # Register for the event
   # Trigger notification check
   curl -X POST http://localhost:8080/notifications/trigger \
     -H "Authorization: TOKEN"

   # Check notifications
   curl http://localhost:8080/notifications \
     -H "Authorization: TOKEN"
   ```

## 🔒 Security Features

- **Password Hashing**: Uses bcrypt with cost factor 14
- **JWT Authentication**: Secure token-based authentication with 2-hour expiration
- **Authorization Middleware**: Protects sensitive endpoints
- **Input Validation**: Request data validation using Gin's binding
- **User Isolation**: Users can only access their own data
- **Event Ownership**: Only event creators can modify their events

## 🚀 Deployment

### Build for Production

```bash
go build -o go-events .
```

### Environment Variables

For production, consider using environment variables for:

- Database connection string
- JWT secret key
- Server port

Example:

```bash
export DB_HOST=your-db-host
export DB_PASSWORD=your-db-password
export JWT_SECRET=your-secret-key
export PORT=8080
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🆘 Troubleshooting

### Common Issues

1. **Database Connection Error**

   - Ensure MySQL is running
   - Check database credentials in `db/db.go`
   - Verify database exists

2. **Authentication Issues**

   - Ensure JWT token is included in Authorization header
   - Check token expiration (2 hours)
   - Verify secret key consistency

3. **Build Errors**
   - Run `go mod tidy` to sync dependencies
   - Check Go version compatibility

### Support

For additional support or questions, please open an issue in the repository.
