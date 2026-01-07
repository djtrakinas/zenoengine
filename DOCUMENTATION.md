# Complete ZenoLang & ZenoEngine Documentation
**Official Guide: Full Edition**

---

## 1. Introduction
ZenoLang is a slot-based configuration language designed for **fast**, **declarative**, and **readable** backend development. ZenoLang runs on **ZenoEngine**, a high-performance Go runtime.

Core Philosophy:
- **Human Friendly:** Code should read like simple English instructions.
- **Slot-Based:** All commands are "slots" (similar to functions) that accept arguments (key-value) and children.
- **Batteries Included:** Built-in features for Database, Auth, HTTP, File System, and Background Jobs.

---

## 2. Core Syntax & Control Flow
This section covers basic commands for programming logic.

### 2.1 Variables & Logs (`var`, `log`)
Storing data and debugging.

```javascript
// Creating a variable
var: $user_name {
  val: "John Doe"
}

// Legacy alias (still supported)
scope.set: $age {
  val: 25
}

// Logging to terminal
log: "User " + $user_name + " is " + $age + " years old"
```

### 2.2 Conditionals (`if`)
Logical branching. Supports operators: `==`, `!=`, `>`, `<`, `>=`, `<=`.

```javascript
if: $age >= 17 {
  then: {
    log: "Adult"
  }
  else: {
    log: "Minor"
  }
}
```

### 2.3 Loops (`while`, `loop`, `for`, `forelse`)
ZenoLang supports several ways to perform loops.

**a. While / Loop (Conditional)**
A loop that runs as long as the condition is true. `while` and `loop` are identical aliases.

```javascript
var: $i { val: 1 }

while: "$i <= 5" {
  do: {
    log: "Iteration " + $i
    math.calc: $i + 1 { as: $i }
    cast.to_int: $i { as: $i } // Ensure it remains an integer
  }
}
```

**b. For (List Iteration / Foreach)**
Iterate over an array or database result list.

```javascript
// $users is a list from the database
for: $users {
  as: $user         // Variable for current item (default: $item)
  do: {
    log: "Name: " + $user.name
  }
}
```

**c. For (C-Style)**
ZenoLang also supports loops in C-style format.

```javascript
for: "$i=1; $i<=3; $i++" {
  do: {
    log: "Iteration " + $i
  }
}
```

**d. Forelse (List with Fallback)**
Similar to `for`, but has an `empty:` block if the data is empty.

```javascript
forelse: $users {
  as: $user
  do: {
    log: $user.name
  }
  empty: {
    log: "No users found."
  }
}
```

**e. Break & Continue (Conditional)**
Use `break` to force stop, and `continue` to move to the next iteration. Now supports direct condition checking!

```javascript
for: $items {
  do: {
    // Stop if ID is 5
    break: "$item.id == 5"
    
    // Skip if status is 'draft'
    continue: "$item.status == 'draft'"
    
    log: $item.title
  }
}
```

### 2.4 Advanced Logic & Conditionals
ZenoLang provides core slots for elegant data checking without using long `if` statements.

**a. Isset & Empty**
Check if a variable exists or is empty.

```javascript
isset: $user {
  do: { log: "User is defined" }
}

empty: $cart {
  do: { log: "Shopping cart is empty" }
}
```

**b. Unless**
The opposite of `if`. Runs the block if the condition is **false**.

```javascript
unless: $is_admin {
  do: { log: "Access denied!" }
}
```

**c. Switch Case**
Brancning values more neatly than nested `if-else`.

```javascript
switch: $status {
  case: "pending"  { do: { log: "Processing" } }
  case: "approved" { do: { log: "Approved" } }
  default:         { do: { log: "Unknown status" } }
}
```

### 2.5 Access Control & Permissions (`auth`, `guest`, `can`)
Special slots for handling login status and user permissions (RBAC).

```javascript
auth: {
  do: { log: "For logged-in users only" }
}

guest: {
  do: { log: "Please login first" }
}

can: "edit-post" {
  resource: $post
  do: { log: "You are allowed to edit this post" }
}
```

**d. Auth Check & User**
Check login status or retrieve user data within ZenoLang scripts.

```javascript
// Check boolean
auth.check: { as: $is_logged_in }

// Get user data
auth.user: { as: $user }
log: "Hello " + $user.username
```

### 2.6 Debugging Tools (`dd`, `dump`)
Laravel-style helpers for quick variable inspection.

- `dump:` Displays variable content to console without stopping.
- `dd:` (Dump and Die) Displays variable content and stops the script immediately.

```javascript
dump: $user
dd: $critical_data
```

### 2.7 Error Handling (`try-catch`)
Handle errors to prevent application crashes. ZenoEngine provides specific and actionable error messages (e.g., notifying if an attribute is misspelled or a `do:` block is missing).

```javascript
try: {
  do: {
    // Code that might error
    db.execute: "DELETE FROM users"
  }
  catch: {
    // Detailed error message, e.g.:
    // "validation error: unknown attribute 'tablee' for slot 'db.select'. Allowed attributes: table, where, ..."
    log: "An error occurred: " + $error
  }
}
```

### 2.8 Global Function Definitions (`fn`, `call`)
ZenoEngine supports defining reusable global functions.

```javascript
// 1. Function Definition (Global Scope)
fn: calculate_discount {
  // Function logic
  math.calc: $total * 0.1 { as: $discount }
  log: "Discount calculated: " + $discount
}

// 2. Call Function
var: $total { val: 100000 }
call: calculate_discount
// Output: Discount calculated: 10000
```

### 2.9 Other Utilities
- **Concurrency Control:** `sleep: 1000` (Delay in milliseconds)
- **Timeout:** `ctx.timeout: "5s" { do: { ... } }` (Limit execution time)
- **Include Script:** `include: "src/modules/other.zl"` (Modularize code)

---

## 3. HTTP & Web Server
ZenoEngine has a powerful built-in HTTP server.

### 3.1 Routing
Routing is defined at the root file level.

```javascript
http.get: /hello {
  do: {
    http.ok: { message: "Hello World" }
  }
}

http.post: /submit {
  do: { ... }
}

// Route Grouping
http.group: /api/v1 {
  do: {
    include: "routes/v1.zl"
  }
}
```

### 3.2 Retrieving Request Data
- **Query Param:** `http.query: "page" { as: $page }`
- **Form Data:** `http.form: "email" { as: $email }`
- **File Upload:** See the *File Upload* section.

### 3.3 Response Helpers
These helpers automatically wrap the response in standard JSON: `{ "success": boolean, ... }`.

| Slot | HTTP Status | Purpose |
| :--- | :--- | :--- |
| `http.ok` | 200 | General success |
| `http.created` | 201 | Resource created successfully |
| `http.accepted` | 202 | Request accepted (processing) |
| `http.no_content`| 204 | Success with empty body |
| `http.bad_request`| 400 | Client input error |
| `http.unauthorized`| 401 | Not logged in |
| `http.forbidden` | 403 | No access rights |
| `http.not_found` | 404 | Data not found |
| `http.validation_error`| 422 | Input validation failure |
| `http.server_error`| 500 | Internal system error |

**Manual Response Example:**
```javascript
http.response: 200 {
  type: "text/html" // Default: application/json
  body: "<h1>Custom HTML</h1>"
}
```

### 3.4 Cookies & Redirects
```javascript
// Set Cookie
cookie.set: {
  name: "session_token"
  val: "xyz123"
  age: 3600 // Seconds
}

// Redirect
http.redirect: "/login"
```

---

## 4. Database (SQL)
ZenoEngine supports **MySQL**, **PostgreSQL**, **SQLite**, and **SQL Server**.

### 4.1 Query Builder (Recommended)
Safe and easy way to interact with the database without raw SQL.

**a. Retrieving Data (`db.get`, `db.first`, `db.count`)**
```javascript
db.table: users
db.where: { col: status, val: "active" }
db.where: { col: age, op: ">=", val: 18 } // Supports operators: >, <, >=, <=, LIKE
db.order_by: "created_at DESC"
db.limit: 10
db.offset: 0

// Execution:
db.get: { as: $users }       // List of Maps
// or
db.first: { as: $user }      // Single Map
// or
db.count: { as: $total }     // Integer
```

**b. Data Manipulation (`db.insert`, `db.update`, `db.delete`)**
```javascript
// Insert
db.table: users
db.insert: {
  name: "John"
  email: "john@test.com"
}
// Auto-increment ID is available in $db_last_id

// Update
db.table: users
db.where: { col: id, val: $id }
db.update: {
  status: "banned"
}

// Delete
db.table: users
db.where: { col: id, val: $id }
db.delete: { as: $affected_rows }
```

### 4.2 Raw SQL
For complex queries that the Query Builder cannot handle.

```javascript
// Select
db.select: "SELECT * FROM users WHERE id = ?" {
  val: $id
  as: $result
  first: true // If you want only one row
}

// Execute (Insert/Update/Delete/DDL)
db.execute: "UPDATE users SET accessed_at = NOW()"
```

### 4.3 Database Transactions (ACID)
Ensures data integrity by wrapping multiple operations in a single atomic transaction. If one operation fails, **all changes are reverted (Rollback)**.

```javascript
db.transaction: {
  do: {
    // 1. Deduct Sender Balance
    db.table: accounts
    db.where: { col: id, val: $sender_id }
    db.update: { balance: $new_sender_balance }

    // 2. Add Receiver Balance
    db.table: accounts
    db.where: { col: id, val: $receiver_id }
    db.update: { balance: $new_receiver_balance }
    
    // 3. Record Transaction Log
    db.table: transactions
    db.insert: {
      amount: $amount
      from: $sender_id
      to: $receiver_id
    }
}
// Automatically Rollbacks if any error occurs inside the 'do' block.
```

### 4.4 Multi-Database Support (Database Agnostic)
ZenoEngine is **database agnostic**! You can easily switch between MySQL, SQLite, PostgreSQL, or SQL Server without changing application code.

**a. Database Configuration in `.env`**
```env
# Choose database driver (mysql, sqlite, postgres, sqlserver)
DB_DRIVER=mysql

# For MySQL
DB_HOST=127.0.0.1:3306
DB_USER=root
DB_PASS=password
DB_NAME=my_database

# For SQLite
# DB_DRIVER=sqlite
# DB_NAME=./database.db
```

**b. Multiple Database Connections**
ZenoEngine supports connecting to multiple databases simultaneously. Simply add a `DB_<NAME>_` prefix in `.env`:

```env
# Default Database
DB_DRIVER=mysql
DB_NAME=app_main

# Warehouse Database
DB_WAREHOUSE_HOST=192.168.1.100
DB_WAREHOUSE_NAME=warehouse
```

**c. Using a Specific Database**
Query builder slots support a `db:` parameter to select a connection:

```javascript
// Default DB
db.table: users
db.get: { as: $users }

// Warehouse DB
db.table: analytics {
  db: warehouse
}
db.get: { as: $stats }
```

---

## 5. Input Validation
Validate incoming data easily.

```javascript
validate: $form_data {
  rules: {
    username: "required|min:5"
    email: "required|email"
    age: "numeric|min:18|max:100"
  }
  as: $errors
}

if: $errors_any { // Automatic helper flag
  then: {
    http.validation_error: { errors: $errors }
  }
}
```

---

## 6. Authentication (Auth)
Built-in JWT Authentication system.

### 6.1 Login
Verifies username/password against database hashes (Bcrypt).
```javascript
auth.login: {
  username: $input_email
  password: $input_password
  table: "users"        // Default
  col_user: "email"     // Default
  col_pass: "password"  // Default
  secret: "APP_SECRET"
  as: $token            // JWT Token string
}
```

### 6.2 Middleware (Route Protection)
```javascript
auth.middleware: {
  secret: "APP_SECRET"
  do: {
    // Only runs if token is valid
    auth.user: { as: $current_user }
    http.ok: { message: "Secret Data", user: $current_user }
  }
}
```

---

## 7. Filesystem & Uploads

### 7.1 File Upload
Handles `multipart/form-data` automatically.
```javascript
http.upload: {
  field: "avatar"
  dest: "public/uploads"
  as: $filename
}
```

### 7.2 File Manipulation (IO)
```javascript
// Write File
io.file.write: {
  path: "logs/app.log"
  content: "Server started..."
}

// Read File
io.file.read: "config.json" { as: $content }

// Delete File
io.file.delete: "temp/old.tmp"
```

---

## 8. System Utilities & Math

### 8.1 Strings & Text
```javascript
// Slugify (Title to URL)
text.slugify: "ZenoEngine Guide 2024" { as: $slug }
// Result: "zenoengine-guide-2024"

// Sanitize (Anti XSS)
text.sanitize: "<script>alert('hack')</script>Hello" { as: $clean }
// Result: "Hello"
```

### 8.2 Math
```javascript
// Basic Math (Float)
math.calc: ceil($price * 1.1) { as: $taxed_price }
// Functions: ceil, floor, round, abs, max, min, sqrt, pow

// Financial (Decimal Precision)
money.calc: ($subtotal - $discount) + $tax { as: $total_fix }
```

### 8.3 Date & Time
```javascript
// Get Current Time
date.now: { as: $now, format: "Human" }
// Result: $now (string), $now_obj (time.Time object)

// Format Date
date.format: $user.created_at {
  layout: "dd MMMM yyyy HH:mm"
  as: $pretty_date
}
```

---

## 9. Email & Background Jobs

### 9.1 Sending Email
```javascript
mail.send: $user_email {
  subject: "Welcome!"
  body: "<h1>Hello!</h1>"
  host: "smtp.mailtrap.io"
  port: 587
  user: "smtp_user"
  pass: "smtp_pass"
}
```

### 9.2 Background Workers (Redis/DB)
ZenoEngine supports robust queue-based task processing.

```javascript
// 1. Worker Config (Required in main.zl or boot)
worker.config: ["high", "default", "low"]

// 2. Enqueue Job
job.enqueue: {
  queue: "default"
  payload: {
    task: "send_email"
    email: "user@test.com"
  }
}
```

---

## 10. Templating (Blade)
For rendering dynamic HTML, use `.blade.zl` files in the `views/` folder. Syntax is similar to Laravel Blade.

**Example View:**
`view.blade: "dashboard.blade.zl" { data: $users }`

**Blade Syntax:**
- `{{ $variable }}` : Output (escaped)
- `{!! $html !!}` : Raw Output
- `@if(...) @else @endif`
- `@foreach($items as $item) ... @endforeach`
- `@extends('layout')`
- `@section('content') ... @endsection`
- `@include('partials.header')`

---

## 11. Advanced Logic & Control Flow
Advanced features available for both `.blade.zl` and regular `.zl` files.

### 11.1 Loop & Branching
```javascript
// Switch Case
logic.switch: $status {
  case: "pending" { do: { log: "Waiting" } }
  case: "success" { do: { log: "Success" } }
  default: { do: { log: "Unknown" } }
}

// Foreach with Else (fallback if empty)
logic.forelse: $items {
  as: $item
  do: {
    log: $item.name
  }
  empty: {
    log: "No data"
  }
}
```

---

## 12. Additional Features (JSON, SSE, Cache)

### 12.1 JSON Manipulation
```javascript
// Parse JSON String to Object
json.parse: $json_string { as: $obj }

// Object to JSON String
json.stringify: $obj { as: $json_string }
```

### 12.2 Realtime (Server-Sent Events)
Stream data in realtime to the browser.
```javascript
sse.stream: {
  sse.send: {
    event: "welcome"
    data: "Hello!"
  }
}
```

---

## 13. Testing Framework
Built-in testing framework for unit testing your logic.

```javascript
test: "Check Addition" {
  math.calc: 1 + 1 { as: $result }
  assert.eq: $result { expected: 2 }
}
```

---

## 14. Modern Frontend (Inertia.js)
First-party support for **Inertia.js** (React, Vue, Svelte).

```javascript
inertia.render: {
  component: "Dashboard/Index"
  props: {
    stats: $stats
    user: $auth_user
  }
}
```

---

## 15. Advanced Query Builder
Standardized features equivalent to modern frameworks like Laravel.

```javascript
db.table: users
db.columns: ["id", "name", "email"]
db.join: {
    table: posts
    on: ["users.id", "=", "posts.user_id"]
}
db.where_in: {
    col: "category_id"
    val: [1, 5, 12]
}
db.get: { as: $results }
```

---

## 16. Excel Export
Powerfully generate `.xlsx` files using templates.

```javascript
excel.from_template: "templates/invoice.xlsx" {
  filename: "Invoice.xlsx"
  data: {
    "customer": "John Doe"
    "items": $billing_items
  }
}
```

---

**End of Documentation.** This guide covers standard ZenoEngine features. For more details on technical specifications, see [LANGUAGE_SPECIFICATION.md](LANGUAGE_SPECIFICATION.md).
