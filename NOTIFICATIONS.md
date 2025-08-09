# Notification System

This document describes the notification system that has been added to the Go Events application.

## Overview

The notification system automatically checks for upcoming events and creates notifications for registered users. It runs as a background job that checks every hour for events happening within the next 24 hours.

## Features

### Background Job

- **Automatic Processing**: Runs every hour to check for upcoming events
- **Smart Notifications**: Only creates one notification per day per event per user
- **Event Detection**: Finds events happening within the next 24 hours
- **User Targeting**: Notifies only users who are registered for the event

### Database Schema

A new `notifications` table has been created with the following structure:

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

### API Endpoints

The following authenticated endpoints have been added:

#### GET `/notifications`

- **Description**: Fetch all notifications for the authenticated user
- **Authorization**: Required (JWT token)
- **Response**: Array of notification objects

```json
[
  {
    "id": 1,
    "user_id": 123,
    "event_id": 456,
    "message": "Reminder: Your event 'Conference 2024' is in 2 hour(s) at 2:00 PM on Dec 15",
    "type": "upcoming_event",
    "is_read": false,
    "created_at": "2024-12-15T10:00:00Z"
  }
]
```

#### PUT `/notifications/:id/read`

- **Description**: Mark a specific notification as read
- **Authorization**: Required (JWT token)
- **Parameters**: `id` - notification ID
- **Security**: Ensures users can only mark their own notifications as read

#### POST `/notifications/trigger` (Development/Testing)

- **Description**: Manually trigger the notification processing
- **Authorization**: Required (JWT token)
- **Use Case**: Testing and manual runs of the notification system

## Message Types

The system generates contextual messages based on event timing:

- **Within 1 hour**: "Reminder: Your event 'EventName' is starting soon at 3:04 PM!"
- **Within 24 hours**: "Reminder: Your event 'EventName' is in X hour(s) at 3:04 PM on Jan 2"
- **Beyond 24 hours**: "Reminder: You have an upcoming event 'EventName' on January 2, 2006 at 3:04 PM"

## Implementation Details

### Background Job Service

- **Location**: `jobs/notification_job.go`
- **Service**: `NotificationService`
- **Frequency**: Every 1 hour
- **Startup**: Automatically starts when the application launches

### Models

- **Notification Model**: `models/notification.go`
- **Key Functions**:
  - `Save()`: Save new notification
  - `GetNotificationsByUserID()`: Fetch user notifications
  - `MarkNotificationAsRead()`: Mark notification as read
  - `GetUpcomingEventsForNotification()`: Find events needing notifications

### Integration

The notification service is automatically started in `main.go` when the application launches:

```go
notificationService := jobs.NewNotificationService()
notificationService.Start()
```

## Usage Example

1. **User Registration**: User registers for an event
2. **Background Processing**: System checks every hour for upcoming events
3. **Notification Creation**: If an event is within 24 hours, creates notification
4. **API Access**: User can fetch notifications via `/notifications` endpoint
5. **Mark as Read**: User can mark notifications as read via PUT endpoint

## Testing

To test the notification system:

1. Create an event with a date/time within the next 24 hours
2. Register a user for that event
3. Call `POST /notifications/trigger` to manually run the notification check
4. Call `GET /notifications` to see the created notification

## Future Enhancements

Potential improvements to consider:

- Email/SMS notifications
- Different notification types (reminders, cancellations, updates)
- Configurable notification timing
- Push notifications for mobile apps
- Notification preferences per user
