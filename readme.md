# 用户管理接口文档

## 概述
此文档描述了用户管理系统的两个主要接口：用户登录和用户注册。这些接口基于Gin框架开发，并与MySQL数据库和Redis进行交互。每个接口都采用JSON格式的请求和响应。

## 接口一：用户登录

### URL
`POST /login`

### 请求参数
```json
{
    "username": "string",  // 用户名，必填
    "password": "string"   // 密码，必填
}
```

### 响应参数
- 成功响应
```json
{
    "message": "Login successful",
    "session_token": "string"  // 会话令牌
}
```

- 失败响应
```json
{
    "error": "Invalid request"         // 请求格式错误
}
或
```json
{
    "error": "Invalid username or password"  // 用户名或密码错误
}
或
```json
{
    "error": "Failed to create session"  // 创建会话失败
}
```

### 功能描述
此接口用于用户登录。请求参数包含用户名和密码。系统会验证用户名和密码是否正确，如果正确则生成JWT会话令牌并返回。如果用户名或密码不正确，则返回错误信息。

### 示例请求
```sh
curl -X POST "http://localhost:8080/login" -H "Content-Type: application/json" -d '{"username": "testuser", "password": "testpassword"}'
```

### 示例响应
```json
{
    "message": "Login successful",
    "session_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

## 接口二：用户注册

### URL
`POST /signup`

### 请求参数
```json
{
    "username": "string",  // 用户名，必填
    "password": "string",  // 密码，必填
    "email": "string",      // 邮箱，必填
    "captchaId": "string",  // 验证码ID，必填
    "value": "string"  // 验证码，必填
}
```

### 响应参数
- 成功响应
```json
{
    "message": "User signup successfully",
    "session_token": "string"  // 会话令牌
}
```

- 失败响应
```json
{
    "error": "Invalid request"         // 请求格式错误
}
或
```json
{
    "error": "Username already exists" // 用户名已存在
}
或
```json
{
    "error": "Failed to create user"  // 创建用户失败
}
或
```json
{
    "error": "Failed to create session"  // 创建会话失败
}
```

### 功能描述
此接口用于用户注册。请求参数包含用户名、密码和邮箱。系统会检查用户名是否已经存在，如果不存在则创建新用户，并生成JWT会话令牌并返回。如果用户名已存在或创建用户失败，则返回错误信息。

### 示例请求
```sh
curl -X POST "http://localhost:8080/signup" -H "Content-Type: application/json" -d '{"username": "newuser", "password": "newpassword", "email": "newuser@example.com"}'
```

### 示例响应
```json
{
    "message": "User signup successfully",
    "session_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

## 错误代码
- `Invalid request`: 请求格式错误，通常是缺少必填字段或字段类型错误。
- `Invalid username or password`: 用户名或密码错误。
- `Username already exists`: 用户名已存在。
- `Failed to create user`: 创建用户失败，通常是数据库操作失败。
- `Failed to create session`: 创建会话失败，通常是JWT生成失败。
- `Invalid captcha`: 验证码错误

---


以下是关于验证码获取的接口文档。该文档描述了生成验证码 ID 和获取验证码图片的 API 接口，包括请求方法、URL、参数、响应格式以及示例。

## 接口三：生成验证码 ID

### 请求

- **方法**：`GET`
- **URL**：`/captcha`

### 描述

此接口用于生成一个新的验证码 ID。

### 响应

- **状态码**：
  - `200 OK`：成功生成验证码 ID
  - `500 Internal Server Error`：服务器内部错误

- **响应体**：
  - **成功**：
    ```json
    {
      "captchaId": "some-valid-captcha-id"
    }
    ```
  - **失败**：
    ```json
    {
      "error": "Failed to generate captcha image"
    }
    ```

### 示例

#### 请求示例

```
GET /captcha
```

#### 响应示例

```json
{
  "captchaId": "abc123"
}
```

---

## 接口四: 获取验证码图片

### 请求

- **方法**：`GET`
- **URL**：`/captcha/:captchaId`

### 描述

此接口用于根据给定的验证码 ID 获取对应的验证码图片。

### 参数

- **路径参数**：
  - `captchaId` (string)：验证码 ID，必须由 `Createcaptchaid` 接口生成。

### 响应

- **状态码**：
  - `200 OK`：成功返回验证码图片
  - `500 Internal Server Error`：服务器内部错误，未能生成验证码图片

- **响应体**：
  - **成功**：返回验证码图片，内容类型为 `image/png`。
  - **失败**：
    ```json
    {
      "error": "Failed to generate captcha image"
    }
    ```

### 示例

#### 请求示例

```
GET /captcha/abc123
```

#### 响应示例

- **成功的响应**：将返回一个 PNG 图片，浏览器会直接显示该图片。

- **失败的响应**：

```json
{
  "error": "Failed to generate captcha image"
}
```

---

## 注意事项

1. 在调用 `GET /captcha/:captchaId` 接口之前，务必先调用 `GET /captcha` 接口以生成验证码 ID。

---

以上就是接口文档的示例。如果需要添加更多细节或其他接口，请告诉我！