# Firebase Authentication API Documentation

This API provides Firebase-based authentication for user registration and login. The authentication system is designed to work with Firebase Client SDK on the frontend.

## Authentication Endpoints

### 1. Get Authentication Info
**GET** `/auth/info`

Returns Firebase project configuration needed for client-side authentication.

**Response:**
```json
{
  "project_id": "your-firebase-project-id",
  "message": "Use Firebase Client SDK for authentication. After authentication, send the ID token to /auth/login or /auth/verify endpoints."
}
```

### 2. Register New User
**POST** `/auth/register`

Creates a new user account in Firebase Auth and saves user data to the database.

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123",
  "display_name": "John Doe" // optional
}
```

**Response:**
```json
{
  "user": {
    "id": 1,
    "firebase_uid": "firebase-uid-string",
    "email": "user@example.com",
    "display_name": "John Doe",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  },
  "token": "firebase-custom-token"
}
```

### 3. Login with ID Token
**POST** `/auth/login`

Authenticates a user using a Firebase ID token obtained from client-side authentication.

**Request Body:**
```json
{
  "id_token": "firebase-id-token-from-client-side-auth"
}
```

**Response:**
```json
{
  "user": {
    "id": 1,
    "firebase_uid": "firebase-uid-string",
    "email": "user@example.com",
    "display_name": "John Doe",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  },
  "token": "firebase-custom-token"
}
```

### 4. Simple Login (Returns Instructions)
**POST** `/auth/simple-login`

Accepts email/password credentials and returns instructions for client-side authentication.

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "message": "Please use Firebase Client SDK to authenticate with these credentials, then send the ID token to /auth/login endpoint",
  "email": "user@example.com",
  "next_step": "Use Firebase signInWithEmailAndPassword() and send the resulting ID token to /auth/login"
}
```

### 5. Verify Token
**POST** `/auth/verify`

Verifies a Firebase ID token and returns user information. Creates user in database if not exists.

**Request Body:**
```json
{
  "id_token": "firebase-id-token"
}
```

**Response:**
```json
{
  "user": {
    "id": 1,
    "firebase_uid": "firebase-uid-string",
    "email": "user@example.com",
    "display_name": "John Doe",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  },
  "uid": "firebase-uid-string"
}
```

## Authentication Flow

### For Frontend Applications:

1. **Registration:**
   - Call `/auth/register` with email, password, and optional display name
   - Store the returned token for future API calls

2. **Login:**
   - **Option A (Recommended):** Use Firebase Client SDK to authenticate, then call `/auth/login` with the ID token
   - **Option B:** Call `/auth/simple-login` to get instructions, then follow the client-side authentication flow

3. **Token Verification:**
   - Use `/auth/verify` endpoint to verify tokens and get user information

### Client-Side Implementation Example (JavaScript):

```javascript
// Initialize Firebase (you'll need your Firebase config)
import { initializeApp } from 'firebase/app';
import { getAuth, signInWithEmailAndPassword, createUserWithEmailAndPassword } from 'firebase/auth';

const firebaseConfig = {
  // Your Firebase configuration
};

const app = initializeApp(firebaseConfig);
const auth = getAuth(app);

// Login flow
async function loginUser(email, password) {
  try {
    // Authenticate with Firebase
    const userCredential = await signInWithEmailAndPassword(auth, email, password);
    const idToken = await userCredential.user.getIdToken();

    // Send ID token to your backend
    const response = await fetch('/auth/login', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ id_token: idToken }),
    });

    const data = await response.json();
    return data;
  } catch (error) {
    console.error('Login failed:', error);
    throw error;
  }
}
```

## Error Responses

All endpoints return error responses in the following format:

```json
{
  "error": "Error message description"
}
```

Common HTTP status codes:
- `400` - Bad Request (invalid input)
- `401` - Unauthorized (invalid credentials/token)
- `404` - Not Found (user not found)
- `500` - Internal Server Error

## Environment Variables Required

- `FIREBASE_SERVICE_ACCOUNT_KEY` - Path to Firebase service account JSON file
- `FIREBASE_PROJECT_ID` - Firebase project ID (optional, for auth info endpoint)

## Protected API Routes

After authentication, use the token in the Authorization header for protected routes:

```
Authorization: Bearer <your-firebase-id-token>
```

Protected routes include:
- `/api/articles`
- `/api/sources`
- `/api/alerts`
- `/api/monitor/trigger`