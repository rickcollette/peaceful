# Peaceful examples

## hello.go
This is a simple router example when you execute:

```bash
go run hello.go
```
 then:

```bash
curl http://localhost:8080/hello
```

you should see:

```json
{
    "ok"
}
```

## jawh-example.go
The `jaht-example.go` file provides a simple example of how to use the `jaht` package for JWT authentication in a Go web application. It includes routes to generate a JWT token and a protected route that requires a valid JWT token to access.

### Application Structure

The application creates a new router and adds a middleware to validate JWT tokens on all incoming requests. It includes two primary routes:

1. **Generate Token Route (`/generate-token`)**: This route generates a JWT token for a given user ID and sends it back in the response.

2. **Protected Route (`/protected`)**: This is a protected route that can only be accessed with a valid JWT token.

### Running the Application

To run the application, navigate to the directory containing the `jaht-example.go` file and execute the following command:

```sh
go run jaht-example.go
```

The application will start, and the server will listen on port 8080.

### Testing with cURL

#### Generating a JWT Token

You can generate a JWT token by making a GET request to the `/generate-token` endpoint. Use the following cURL command:

```sh
curl http://localhost:8080/generate-token
```

The response will contain a JWT token:

```json
{
    "token": "your-jwt-token"
}
```

#### Accessing the Protected Route

To access the protected route, include the JWT token in the Authorization header of your request. Replace `"your-jwt-token"` with the actual token received from the `/generate-token` endpoint. Use the following cURL command:

```sh
curl -H "Authorization: your-jwt-token" http://localhost:8080/protected
```

If the token is valid, you'll receive a response from the protected route:

```json
{
    "message": "Welcome to the protected route!"
}
```

If the token is invalid or expired, you'll receive a 401 Unauthorized error.

## Note

Ensure to replace `"your-secret-key"` with your actual secret key used for signing the JWT tokens, and secure it appropriately in a production environment.
