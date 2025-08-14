# Postman Testing Guide - E-Commerce Analytics System

This guide explains how to test the Phase 1 e-commerce analytics and personalization system using Postman.

## Table of Contents
1. [Getting Started](#getting-started)
2. [Environment Setup](#environment-setup)
3. [Authentication](#authentication)
4. [Event Management](#event-management)
5. [Segment Management](#segment-management)
6. [Offer Management](#offer-management)
7. [WebSocket Testing](#websocket-testing)
8. [Complete Testing Workflow](#complete-testing-workflow)

## Getting Started

### Prerequisites
- Docker and Docker Compose installed
- Postman installed
- Git (to clone the repository)

### Starting the Backend

1. **Clone and start the backend:**
   ```bash
   git clone <repository-url>
   cd cusror_ai
   make run
   # OR
   docker-compose up --build
   ```

2. **Verify the backend is running:**
   - API should be available at: `http://localhost:8081`
   - WebSocket endpoint: `ws://localhost:8081/ws`
   - Health check: `GET http://localhost:8081/health`

3. **Create initial segments:**
   ```bash
   curl -X POST http://localhost:8081/api/v1/segments/predefined
   ```

## Environment Setup

### Postman Environment Variables
Create a new Postman environment with these variables:

```json
{
  "base_url": "http://localhost:8081",
  "api_base": "http://localhost:8081/api/v1",
  "ws_url": "ws://localhost:8081/ws",
  "auth_token": "",
  "user_id": "user123",
  "test_event_id": "",
  "test_segment_id": "",
  "test_offer_id": ""
}
```

## Authentication

### 1. Create Test User
**POST** `{{api_base}}/users`

**Body (JSON):**
```json
{
  "name": "Test User",
  "email": "test@example.com",
  "password": "password123"
}
```

### 2. Login
**POST** `{{base_url}}/login`

**Body (JSON):**
```json
{
  "email": "test@example.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Set Environment Variable:**
In Postman Tests tab, add:
```javascript
if (pm.response.code === 200) {
    const response = pm.response.json();
    pm.environment.set("auth_token", response.token);
}
```

### 3. Test Protected Endpoint
**GET** `{{base_url}}/dashboard`

**Headers:**
```
Authorization: Bearer {{auth_token}}
```

## Event Management

### 1. Create Product View Event
**POST** `{{api_base}}/events`

**Body (JSON):**
```json
{
  "user_id": "{{user_id}}",
  "event_type": "product_view",
  "product_id": "product-123",
  "category": "electronics",
  "metadata": {
    "product_name": "iPhone 15",
    "price": 999.99,
    "page_url": "/products/iphone-15"
  }
}
```

### 2. Create Purchase Event
**POST** `{{api_base}}/events`

**Body (JSON):**
```json
{
  "user_id": "{{user_id}}",
  "event_type": "purchase",
  "product_id": "product-123",
  "category": "electronics",
  "metadata": {
    "amount": 999.99,
    "order_id": "order-456",
    "currency": "USD"
  }
}
```

### 3. Create Cart Add Event
**POST** `{{api_base}}/events`

**Body (JSON):**
```json
{
  "user_id": "{{user_id}}",
  "event_type": "cart_add",
  "product_id": "product-789",
  "category": "electronics",
  "metadata": {
    "quantity": 2,
    "price": 299.99
  }
}
```

### 4. Simulate Shopify Event
**POST** `{{api_base}}/events/simulate/shopify`

**Body (JSON):**
```json
{
  "event_type": "order_created",
  "customer_id": "{{user_id}}",
  "order_data": {
    "order_id": "shopify-order-123",
    "total_price": 1299.98,
    "currency": "USD"
  }
}
```

### 5. Get User Events
**GET** `{{api_base}}/events/user/{{user_id}}?limit=10&offset=0`

### 6. Get User Event Statistics
**GET** `{{api_base}}/events/user/{{user_id}}/stats`

### 7. Get Recent Events
**GET** `{{api_base}}/events/recent?limit=20`

## Segment Management

### 1. Create Predefined Segments
**POST** `{{api_base}}/segments/predefined`

### 2. Get All Segments
**GET** `{{api_base}}/segments?active_only=true`

### 3. Create Custom Segment
**POST** `{{api_base}}/segments`

**Body (JSON):**
```json
{
  "name": "Tech Enthusiasts",
  "description": "Users who frequently purchase electronics",
  "rules": {
    "rules": [
      {
        "field": "product_views",
        "operator": "gte",
        "value": 5
      },
      {
        "field": "total_spent",
        "operator": "gte",
        "value": 500
      }
    ]
  }
}
```

### 4. Analyze User Segments
**POST** `{{api_base}}/segments/analyze/{{user_id}}`

### 5. Get User Segments
**GET** `{{api_base}}/segments/user/{{user_id}}`

### 6. Get User Interests
**GET** `{{api_base}}/segments/user/{{user_id}}/interests`

### 7. Get Segment Rules Template
**GET** `{{api_base}}/segments/rules/template`

## Offer Management

### 1. Get Active Offers
**GET** `{{api_base}}/offers`

### 2. Get User-Specific Offers
**GET** `{{api_base}}/offers/user/{{user_id}}`

### 3. Use an Offer
**POST** `{{api_base}}/offers/{{test_offer_id}}/use`

**Body (JSON):**
```json
{
  "user_id": "{{user_id}}",
  "order_id": "order-789",
  "amount": 25.50
}
```

### 4. Get Offer Statistics
**GET** `{{api_base}}/offers/stats`

## WebSocket Testing

### Testing Real-time Notifications

1. **Connect to WebSocket:**
   - URL: `ws://localhost:8081/ws?user_id={{user_id}}`
   - Use a WebSocket client or browser console

2. **JavaScript WebSocket Example:**
   ```javascript
   const ws = new WebSocket('ws://localhost:8081/ws?user_id=user123');
   
   ws.onopen = function(event) {
       console.log('Connected to WebSocket');
   };
   
   ws.onmessage = function(event) {
       const message = JSON.parse(event.data);
       console.log('Received:', message);
   };
   
   ws.onclose = function(event) {
       console.log('WebSocket connection closed');
   };
   ```

3. **Test Real-time Notifications:**
   - Create events via API and watch for `event_processed` messages
   - Analyze user segments and watch for `segment_changed` messages
   - Use offers and watch for `offer_used` messages

### WebSocket Message Types
- `event_processed`: When a new event is processed
- `segment_changed`: When user segments change
- `offer_generated`: When new offers are created
- `offer_used`: When offers are used
- `heartbeat`: Periodic keepalive messages

## Complete Testing Workflow

### Scenario: New User Journey

1. **Create a new user** (Authentication step 1)
2. **Login** (Authentication step 2)
3. **Create initial events** to simulate user behavior:
   ```json
   [
     { "event_type": "page_view", "metadata": {"page_url": "/"} },
     { "event_type": "product_view", "product_id": "prod1", "category": "electronics" },
     { "event_type": "product_view", "product_id": "prod2", "category": "electronics" },
     { "event_type": "cart_add", "product_id": "prod1", "category": "electronics" },
     { "event_type": "purchase", "product_id": "prod1", "category": "electronics", "metadata": {"amount": 299.99}}
   ]
   ```

4. **Analyze user segments** to see which segments the user belongs to
5. **Check for generated offers** specific to the user
6. **Use an offer** to complete the personalization loop
7. **Monitor WebSocket** for real-time notifications throughout

### Performance Testing

1. **Bulk Event Creation:**
   Create a Postman collection runner to simulate multiple events rapidly

2. **Load Testing:**
   Use the Shopify simulation endpoint to create realistic e-commerce scenarios

3. **WebSocket Stress Test:**
   Open multiple WebSocket connections and monitor performance

## API Response Examples

### Event Creation Response
```json
{
  "event_id": "123e4567-e89b-12d3-a456-426614174000",
  "user_id": "user123",
  "event_type": "product_view",
  "product_id": "product-123",
  "category": "electronics",
  "timestamp": "2024-01-15T10:30:00Z",
  "metadata": {
    "product_name": "iPhone 15",
    "price": 999.99
  },
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

### Segment Analysis Response
```json
{
  "user_id": "user123",
  "segments": [
    {
      "segment_id": "456e7890-e89b-12d3-a456-426614174001",
      "segment_name": "High Spender",
      "description": "Users who spend more than $500 on average per session",
      "is_active": true
    }
  ],
  "count": 1
}
```

### Offer Response
```json
{
  "offer_id": "789e0123-e89b-12d3-a456-426614174002",
  "segment_id": "456e7890-e89b-12d3-a456-426614174001",
  "user_id": "user123",
  "title": "VIP Exclusive: 15% Off Premium Items",
  "description": "Exclusive discount for our valued premium customers",
  "discount_percent": 15.0,
  "coupon_code": "OFFER789e0123",
  "valid_from": "2024-01-15T10:30:00Z",
  "valid_until": "2024-02-14T10:30:00Z",
  "is_active": true,
  "usage_count": 0
}
```

## Troubleshooting

### Common Issues

1. **Connection Refused (Port 8081):**
   - Ensure Docker containers are running: `docker-compose ps`
   - Check logs: `docker-compose logs app`

2. **Database Connection Errors:**
   - Wait for PostgreSQL to be ready (health check)
   - Check database credentials in docker-compose.yml

3. **WebSocket Connection Issues:**
   - Verify WebSocket URL format
   - Check browser console for CORS errors
   - Ensure WebSocket service is running

4. **No Segments Created:**
   - Run predefined segments creation endpoint first
   - Check segment rules are valid JSON

5. **No Offers Generated:**
   - Ensure user has enough events to match segment rules
   - Verify segment analysis was triggered
   - Check WebSocket notifications for debugging

### Health Check Endpoints

- **API Health:** `GET /health`
- **Database Status:** `GET /api/v1/users` (should return empty array)
- **WebSocket Status:** `GET /api/v1/ws/stats`

## Next Steps

After completing Phase 1 testing:
1. Collect performance metrics
2. Analyze user segmentation accuracy
3. Evaluate offer relevance and conversion rates
4. Prepare for Phase 2 AI-powered enhancements

---

For additional support or questions, refer to the project README or contact the development team. 