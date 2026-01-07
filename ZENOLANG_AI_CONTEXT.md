# ZenoLang AI Context & Syntax Reference

> **SYSTEM PROMPT / CONTEXT FOR AI AGENTS**
> Use this file as the SOURCE OF TRUTH when writing ZenoLang code.
> ZenoLang is a slot-based configuration language running on ZenoEngine (Go).

---

## 1. Syntax Rules (STRICT)

1.  **Structure**: Tree-based structure using indentation (2 spaces) and braces `{}`.
2.  **Slots**: Format is `slot_name: value { children }`.
    *   Example: `log: "Hello"`
    *   Example: `if: $condition { ... }`
3.  **Variables**: Always prefixed with `$` (e.g., `$user_id`).
    *   **Assignment**: `var: $name { val: "value" }`
    *   **Access**: `$user.name` (Dot notation), `$items.0` (Index).
    *   **Interpolation**: `"Hello " + $name` (String concatenation).
4.  **Quoting**:
    *   Strings *should* be quoted if they contain spaces or special chars: `"valid string"`.
    *   Keys in maps DO NOT need quotes: `name: "Budi"`.
5.  **Comments**: Use `//` for single line comments.

---

## 2. API Reference (Slots)

### 2.1 Database (SQL)
**Driver Agnostic**: MySQL, PostgreSQL, SQLite, SQL Server.

| Slot | Signature / Children | Description |
| :--- | :--- | :--- |
| `db.table` | `value: "table_name"`, `db: "conn_name"` | Set active table. |
| `db.get` | `as: $var` | Fetch multiple rows (List of Maps). |
| `db.first` | `as: $var` | Fetch single row (Map). |
| `db.count` | `as: $var` | Count rows (Int). |
| `db.insert` | Child keys as columns. | Insert data. Returns `$db_last_id`. |
| `db.update` | Child keys as columns. | Update data. Requires `db.where`. |
| `db.delete` | `as: $count` (optional) | Delete data. Requires `db.where`. |
| `db.where` | `col: "name"`, `val: $val`, `op: "="` | Add WHERE clause. |
| `db.where_in` | `col: "id"`, `val: [1, 2]` | WHERE IN clause. |
| `db.where_not_in`| `col: "id"`, `val: [1, 2]` | WHERE NOT IN clause. |
| `db.where_null` | `value: "col_name"` | WHERE col IS NULL. |
| `db.where_not_null`| `value: "col_name"` | WHERE col IS NOT NULL. |
| `db.join` | `table: "other"`, `on: ["t1.id", "=", "t2.fk"]` | INNER JOIN. |
| `db.left_join` | (Same as join) | LEFT JOIN. |
| `db.transaction`| `do: { ... }` | Atomic transaction. Auto-rollback on error. |
| `db.select` | `value: "SQL"`, `val: $p1`, `as: $res` | Raw SQL Select (User input binding safe). |
| `db.execute` | `value: "SQL"` | Raw SQL Execute (User input binding safe). |

### 2.2 HTTP Server & Client

| Slot | Signature / Children | Description |
| :--- | :--- | :--- |
| `http.get`, `post`| `value: "/path"`, `do: { ... }` | Define Route. |
| `http.query` | `value: "param"`, `as: $var` | Get Query Param (`?id=1`). |
| `http.form` | `value: "field"`, `as: $var` | Get Form/Body data. |
| `http.response` | `status: 200`, `data: $json` | Send JSON response. |
| `http.ok`, `created`| `value: { ... }` | Send 200/201 JSON response. |
| `http.redirect` | `value: "/url"` | Redirect user. |
| `http.upload` | `field: "file"`, `dest: "path"`, `as: $name` | Handle File Upload. |

### 2.3 Logic & Flow Control

| Slot | Signature / Children | Description |
| :--- | :--- | :--- |
| `if` | `value: "$a > $b"`, `then: {}`, `else: {}` | Conditional. Ops: `==, !=, >, <, >=, <=`. |
| `switch` | `value: $var`, `case: "val" { ... }` | Switch case. |
| `while` / `loop` | `value: "cond"`, `do: {}` | While loop. |
| `for` | `value: $list`, `as: $item`, `do: {}` | Foreach loop. |
| `try` | `do: {}`, `catch: {}` | Error handling. Error msg in `$error`. |
| `fn` | `value: name`, `children...` | Define global function. |
| `call` | `value: name` | Call global function. |

### 2.4 Utils & Security

| Slot | Signature | Description |
| :--- | :--- |
| `log` | `value: "msg"` | Print to console. |
| `var` | `val: $val` | Set variable. |
| `sleep` | `value: ms` | Sleep for N milliseconds. |
| `coalesce` | `val: $a`, `default: "b"`, `as: $r` | Null coalescing. |
| `is_null` | `val: $a`, `as: $bool` | Check if null. |
| `cast.to_int` | `val: $v`, `as: $i` | Cast to Integer. |
| `crypto.hash` | `val: $pass`, `as: $hash` | Bcrypt hash. |
| `crypto.verify` | `hash: $h`, `text: $p`, `as: $ok` | Verify Bcrypt. |
| `crypto.verify_aspnet`| `hash: $h`, `password: $p`, `as: $ok`| Verify ASP.NET Identity hash. |
| `arrays.length` | `value: $arr`, `as: $len` | Get array length. |

#### 2.5 Authentication (JWT Security) 🔐
**Note**: ZenoEngine uses **Stateless JWT** (JSON Web Tokens). Do NOT use traditional session cookies.

| Slot | Signature | Description |
| :--- | :--- |
| `auth.login`| `id: $id`, `as: $token` | Generate JWT for user ID. |
| `auth.user` | `as: $user` | Get current authenticated user (from token). |
| `auth.check`| `as: $bool` | Check if request has valid token. |
| `auth.middleware`| `do: { ... }` | Protect route (Requires Bearer Token). |
| `sec.csrf_token`| `as: $token` | Get CSRF token (for forms). |

| `sec.csrf_token`| `as: $token` | Get CSRF token (for forms). |

#### 2.6 View Rendering (Blade Engine) 🖥️
ZenoLang supports a Blade-compatible template engine (Laravel-style).

| Slot | Signature | Description |
| :--- | :--- | :--- |
| `view.blade` | `value: "file.html"`, `data_key: $val` | Render HTML template in `/views`. |

**Blade Syntax Support**:
- `{{ $var }}` (Escaped Output)
- `!! $var !!` (Raw Output)
- Directives: `@if`, `@else`, `@foreach`, `@switch`.
- Auth: `@auth`, `@guest`, `@can('ability')`.
- Components: `<x-alert type="error" />`.
- Layouts: `@extends('layout')`, `@section('content')`, `@yield('content')`.

#### 2.7 Modern Web (Inertia & Excel) 🆕

| Slot | Description |
| :--- | :--- |
| `inertia.render` | `component: "Page"`, `props: { ... }` | Render Inertia Page. |
| `excel.from_template` | `filename: "f.xlsx"`, `data: { ... }`, `images: {}`, `formulas: {}` | Export Excel. |

---

## 3. High-Confidence Code Patterns

### 3.1 Standard CRUD (Create)
```javascript
http.post: /users {
  do: {
    // 1. Validation
    validate: $form {
      rules: {
        name: "required"
        email: "required|email"
      }
      as: $errors
    }
    
    // 2. Error Check
    if: $errors_any {
      then: { http.validation_error: { errors: $errors } }
      else: {
        // 3. Database Insert
        db.table: users
        db.insert: {
          name: $form.name
          email: $form.email
          created_at: date.now
        }
        
        http.created: { 
            message: "User created"
            id: $db_last_id 
        }
      }
    }
  }
}
```

### 3.2 Authentication (Login)
```javascript
http.post: /login {
  do: {
    // 1. Find User
    db.table: users
    db.where: { col: email, val: $form.email }
    db.first: { as: $user }
    
    // 2. Verify Logic
    var: $is_valid { val: false }
    
    if: $user {
      then: {
        crypto.verify: {
          hash: $user.password
          text: $form.password
          as: $is_valid
        }
      }
    }
    
    // 3. Result
    if: $is_valid {
      then: {
        auth.login: { id: $user.id, as: $token }
        http.ok: { token: $token }
      }
      else: {
        http.unauthorized: { message: "Invalid credentials" }
      }
    }
  }
}
```

### 3.3 Background Job
```javascript
// Enqueue
job.enqueue: {
  queue: "default"
  payload: {
    action: "process_report"
    user_id: 101
  }
}
```

### 3.4 Legacy Migration (ASP.NET Core Identity)
Use this pattern to authenticate users migrated from ASP.NET Core database (Identity V3) without resetting passwords.

```javascript
http.post: /login_migrated {
  do: {
    // 1. Get User
    db.table: AspNetUsers
    db.where: { col: Email, val: $form.email }
    db.first: { as: $user }

    // 2. Verify Legacy Hash
    var: $is_valid { val: false }
    
    if: $user {
      then: {
        // Check using ASP.NET Core Identity V3 Verifier
        crypto.verify_aspnet: {
          hash: $user.PasswordHash
          password: $form.password
          as: $is_valid
        }
      }
    }

    // 3. Login or Fail
    if: $is_valid {
      then: {
        // Issue new Zeno Token
        auth.login: { id: $user.Id, as: $token }
        http.ok: { token: $token }
      }
      else: {
        http.unauthorized: { message: "Invalid credentials" }
      }
    }
  }
}
    }
  }
}
```

---

## 4. Appendix: Detailed Reference

### 4.1 Validator Rules (`validator.validate`)
Used in `rules: { field: "rule|rule:param" }`

| Rule | Example | Description |
| :--- | :--- | :--- |
| `required` | `name: "required"` | Field must not be empty. |
| `email` | `email: "email"` | Must be valid email format. |
| `numeric` | `age: "numeric"` | Must be a number. |
| `min:X` | `age: "min:18"`, `pass: "min:8"` | Min value (numeric) or min length (string). |
| `max:X` | `age: "max:100"` | Max value (numeric) or max length (string). |

### 4.2 Date & Time (`date.*`)
Standard layout: `RFC3339` or `"Human"` (`02 Jan 2006 15:04`).

| Slot | Signature |
| :--- | :--- |
| `date.now` | `layout: "Human"`, `as: $now` |
| `date.format`| `val: $date_obj`, `layout: "YYYY-MM-DD"`, `as: $str` |
| `date.parse` | `val: "2024-01-01"`, `as: $date_obj` |
| `date.add` | `val: $date`, `add: "1h"`, `as: $new_date` |

### 4.3 Filesystem (`io.*`)
| Slot | Signature |
| :--- | :--- |
| `io.file.read` | `path: "file.txt"`, `as: $content` |
| `io.file.write`| `path: "file.txt"`, `content: $data`, `mode: 0644` |
| `io.file.delete`| `path: "file.txt"` |
| `io.dir.create` | `path: "folder/sub"` |
