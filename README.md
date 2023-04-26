### JWT Auth API

**Based on**
- Gin
- Golang-jwt

#### API

- POST ```/auth/login``` Generate a jwt token based on login request

    - Login Request Body<br>
    ```json
    {
        "username": "xxx",
        "password": "xxx"
    }
    ```