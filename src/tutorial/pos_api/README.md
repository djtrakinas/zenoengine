# Multi-Tenant POS API Tutorial

## Overview

This tutorial demonstrates how to build a **multi-tenant Point of Sale (POS) API** using ZenoEngine with:
- Multi-tenant architecture (subdomain & header-based)
- JWT authentication
- Tenant isolation
- RESTful API design

## Architecture

### Multi-Tenancy Model

**Tenant Detection**:
1. `X-Tenant-ID` header (priority)
2. Subdomain fallback (e.g., `abc.yourdomain.com`)

**Database Structure**:
- `pos_system` - System database (tenant metadata)
- `pos_tenant_abc` - Tenant A database
- `pos_tenant_xyz` - Tenant B database

### Authentication Flow

```
1. Login → JWT Token (with tenant_id claim)
2. Request → auth.middleware validates token
3. Middleware → Sets $auth object
4. Controller → Uses $auth.user_id
```

## Setup

### 1. Database Configuration

**.env**:
```bash
# System DB (tenant metadata)
DB_SYSTEM_DRIVER=mysql
DB_SYSTEM_HOST=127.0.0.1:3306
DB_SYSTEM_USER=root
DB_SYSTEM_PASS=yourpassword
DB_SYSTEM_NAME=pos_system

# Tenant A
DB_TENANT_ABC_DRIVER=mysql
DB_TENANT_ABC_HOST=127.0.0.1:3306
DB_TENANT_ABC_USER=root
DB_TENANT_ABC_PASS=yourpassword
DB_TENANT_ABC_NAME=pos_tenant_abc

# JWT Secret
JWT_SECRET=your_secret_key_here
```

### 2. Create Databases

```sql
-- System database
CREATE DATABASE pos_system;
USE pos_system;

CREATE TABLE tenants (
    id INT PRIMARY KEY AUTO_INCREMENT,
    code VARCHAR(50) UNIQUE,
    name VARCHAR(100),
    db_connection_name VARCHAR(50),
    is_active BOOLEAN DEFAULT 1
);

INSERT INTO tenants (code, name, db_connection_name) VALUES
('abc', 'Tenant ABC', 'tenant_abc'),
('xyz', 'Tenant XYZ', 'tenant_xyz');

-- Tenant database
CREATE DATABASE pos_tenant_abc;
USE pos_tenant_abc;

CREATE TABLE users (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100),
    email VARCHAR(100) UNIQUE,
    password VARCHAR(255),
    role VARCHAR(50) DEFAULT 'user'
);

CREATE TABLE products (
    id INT PRIMARY KEY AUTO_INCREMENT,
    sku VARCHAR(50) UNIQUE,
    name VARCHAR(100),
    price DECIMAL(10,2),
    stock INT DEFAULT 0
);

CREATE TABLE sales (
    id INT PRIMARY KEY AUTO_INCREMENT,
    invoice_number VARCHAR(50) UNIQUE,
    user_id INT,
    total DECIMAL(10,2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE sale_items (
    id INT PRIMARY KEY AUTO_INCREMENT,
    sale_id INT,
    product_id INT,
    qty INT,
    price DECIMAL(10,2)
);
```

## API Endpoints

### Authentication

**POST /api/v1/auth/login**
```bash
curl -X POST http://localhost:3000/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -H "X-Tenant-ID: abc" \
  -d '{
    "email": "admin@pos.com",
    "password": "password"
  }'
```

Response:
```json
{
  "success": true,
  "token": "eyJhbGci...",
  "user": {
    "id": 1,
    "email": "admin@pos.com",
    "role": "admin"
  }
}
```

### Products

**GET /api/v1/products/**
```bash
curl http://localhost:3000/api/v1/products/ \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "X-Tenant-ID: abc"
```

**POST /api/v1/products/**
```bash
curl -X POST http://localhost:3000/api/v1/products/ \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "X-Tenant-ID: abc" \
  -H "Content-Type: application/json" \
  -d '{
    "sku": "PROD001",
    "name": "Product 1",
    "price": 10000,
    "stock": 100
  }'
```

### Sales

**POST /api/v1/sales/**
```bash
curl -X POST http://localhost:3000/api/v1/sales/ \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "X-Tenant-ID: abc" \
  -H "Content-Type: application/json" \
  -d '{
    "items": [
      {"product_id": 1, "qty": 5}
    ]
  }'
```

## Implementation Guide

### 1. Routes with Authentication

**routes.zl**:
```zenolang
http.post: /api/v1/sales/ {
    summary: "Create POS Transaction"
    auth.middleware: {
        secret: env("JWT_SECRET")
        do: {
            call: create_sale
        }
    }
}
```

### 2. Controller with $auth Object

**controller.zl**:
```zenolang
fn: create_sale {
    // $auth object available from middleware
    $user_id: $auth.user_id
    $tenant_id: $auth.tenant_id
    
    // Use tenant-specific database
    db.insert: {
        db: $CURRENT_TENANT_DB
        table: "sales"
        data: {
            user_id: $user_id
            total: $total
        }
        as: $sale
    }
}
```

### 3. Login Controller

**auth/controller.zl**:
```zenolang
fn: login {
    // Validate credentials
    db.select: {
        db: $CURRENT_TENANT_DB
        table: "users"
        where: { email: $email }
        as: $user
    }
    
    // Generate JWT
    jwt.sign: {
        secret: env("JWT_SECRET")
        claims: {
            user_id: $user.0.id
            email: $user.0.email
            tenant_id: $tenant_id
            role: $user.0.role
        }
        expires_in: 86400
        as: $token
    }
}
```

## Best Practices

### 1. Always Use env("JWT_SECRET")
```zenolang
auth.middleware: {
    secret: env("JWT_SECRET")  // ✅ Good
    // NOT: secret: "hardcoded"  // ❌ Bad
}
```

### 2. Validate Tenant
```zenolang
// Middleware auto-validates tenant from system DB
// Sets: $CURRENT_TENANT_DB, $CURRENT_TENANT_ID, $CURRENT_TENANT_NAME
```

### 3. Use $auth Object
```zenolang
// Available after auth.middleware:
$auth.user_id
$auth.email
$auth.tenant_id
$auth.role
```

### 4. Tenant Isolation
```zenolang
// Always use tenant-specific database
db.insert: {
    db: $CURRENT_TENANT_DB  // ✅ Isolated
    // NOT: db: "default"    // ❌ Shared
}
```

## Security

1. **JWT Secret**: Use strong, random secret in production
2. **HTTPS**: Always use HTTPS in production
3. **Token Expiry**: Set appropriate expiration (default: 24h)
4. **Tenant Validation**: Middleware validates tenant exists and is active
5. **Database Isolation**: Each tenant has separate database

## Testing

```bash
# 1. Login
TOKEN=$(curl -X POST http://localhost:3000/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -H "X-Tenant-ID: abc" \
  -d '{"email":"admin@pos.com","password":"password"}' \
  -s | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

# 2. Create Product
curl -X POST http://localhost:3000/api/v1/products/ \
  -H "Authorization: Bearer $TOKEN" \
  -H "X-Tenant-ID: abc" \
  -H "Content-Type: application/json" \
  -d '{"sku":"TEST","name":"Test Product","price":10000,"stock":50}'

# 3. Create Sale
curl -X POST http://localhost:3000/api/v1/sales/ \
  -H "Authorization: Bearer $TOKEN" \
  -H "X-Tenant-ID: abc" \
  -H "Content-Type: application/json" \
  -d '{"items":[{"product_id":1,"qty":5}]}'
```

## Troubleshooting

### Token Invalid
- Check JWT_SECRET matches in .env
- Verify token not expired
- Ensure X-Tenant-ID header present

### Database Connection
- Verify DB_TENANT_* config in .env
- Check tenant exists in system DB
- Ensure tenant is_active = 1

### Auth Not Working
- Use `auth.middleware` block (not `middleware: "auth"`)
- Include `secret: env("JWT_SECRET")`
- Wrap logic in `do` block

## Production Deployment

1. Set strong JWT_SECRET
2. Use environment variables
3. Enable HTTPS
4. Configure rate limiting
5. Monitor tenant databases
6. Backup regularly

## Summary

This tutorial demonstrates:
- ✅ Multi-tenant architecture
- ✅ JWT authentication
- ✅ Tenant isolation
- ✅ RESTful API
- ✅ Secure by default

**Status**: Production Ready 🚀
