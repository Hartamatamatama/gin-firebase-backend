# Testing Endpoint


### Health Check menggunakan Backend URL (http://localhost:8081)
- Endpoint
```bash
GET

{{BACKEND_BASE_URL}}/v1/health

```

### Example

![](screenshots/health_check.png)

---

### Login untuk mendapatkan idToken


- Endpoint
```bash
POST

https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key={{FIREBASE_API_KEY}}
```

### Example

![](screenshots/login.png)

---

### Verifiy Token menggunakan idToken untuk mendapatkan accessToken.
#### Note: Pastikan bahwa akun sudah mendapatkan pesan verifikasi via email ketika registrasi (daftar akun)
#### Setelah tahap ini, data user akan masuk ke database
- Endpoint
```bash
POST

{{BACKEND_BASE_URL}}/v1/auth/verify-token
```


### Example
![](screenshots/check_verification_via_backend_url.png)
![](screenshots/users_table.png)

---

### Get All Product menggunakan accessToken
- Endpoint
```bash
GET

{{BACKEND_BASE_URL}}/v1/products
```

### Example

![](screenshots/get_all_products.png)

---

### Get Product By id menggunakan accessToken
- Endpoint
```bash
GET

{{BACKEND_BASE_URL}}/v1/products/4
```

### Example

![](screenshots/get_products_by_id.png)

---

# Seeder

- Command Running Seed
```bash
go run seeds/seed.go
```

### Setelah di run

![](screenshots/products_table.png)

---

# Flow (Sequence Diagram)

### Create User

![](screenshots/SD_Create_User.png)

### Code
```bash
sequenceDiagram
participant Client
participant Middleware
participant Handler
participant Service
participant Repository
participant DB
Client->>Middleware: HTTP Request
Middleware->>Handler: forward request
Handler->>Service: CreateUser()
Service->>Repository: SaveUser()
Repository->>DB: INSERT USER
DB-->>Repository: OK
Repository-->>Service: user saved
Service-->>Handler: user response
Handler-->>Client: JSON Response
```

# Get All Product

![](screenshots/SD_Get_All_Product.png)

### Code

```bash
sequenceDiagram
participant Client
participant Middleware
participant Handler
participant Service
participant Repository
participant DB

Client->>Middleware: HTTP Request (Get v1/products)
Middleware->>Handler: forward request
Handler->>Service: GetAll
Service->>Repository: FindAll
Repository->>DB: SELECT * FROM products
DB-->>Repository: product rows
Repository-->>Service: products list
Service-->>Handler: products response
Handler-->>Client: JSON Response
```

# Get Product by Id

![](screenshots/SD_Get_Product_by_Id.png)

### Code

```bash
sequenceDiagram
participant Client
participant Middleware
participant Handler
participant Service
participant Repository
participant DB

Client->>Middleware: HTTP Request (Get v1/products/id)
Middleware->>Handler: forward request
Handler->>Service: GetById
Service->>Repository: FindById
Repository->>DB: SELECT * FROM products where id=?
DB-->>Repository: product rows
Repository-->>Service: products list
Service-->>Handler: products response
Handler-->>Client: JSON Response
```