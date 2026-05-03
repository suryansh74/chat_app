# Graph Report - .  (2026-05-04)

## Corpus Check
- Corpus is ~5,005 words - fits in a single context window. You may not need a graph.

## Summary
- 173 nodes · 339 edges · 15 communities detected
- Extraction: 66% EXTRACTED · 34% INFERRED · 0% AMBIGUOUS · INFERRED: 115 edges (avg confidence: 0.8)
- Token cost: 150 input · 520 output

## Community Hubs (Navigation)
- [[_COMMUNITY_Auth Service Tests|Auth Service Tests]]
- [[_COMMUNITY_HTTP Handlers & Routes|HTTP Handlers & Routes]]
- [[_COMMUNITY_Server Initialization|Server Initialization]]
- [[_COMMUNITY_Auth Middleware|Auth Middleware]]
- [[_COMMUNITY_Error Handling|Error Handling]]
- [[_COMMUNITY_User Repository|User Repository]]
- [[_COMMUNITY_Handler Tests|Handler Tests]]
- [[_COMMUNITY_Email Verification|Email Verification]]
- [[_COMMUNITY_Validation Errors|Validation Errors]]
- [[_COMMUNITY_Domain Models|Domain Models]]
- [[_COMMUNITY_Token Payload|Token Payload]]
- [[_COMMUNITY_Service Ports|Service Ports]]
- [[_COMMUNITY_Token Models|Token Models]]
- [[_COMMUNITY_Repository Port|Repository Port]]
- [[_COMMUNITY_Project Notes|Project Notes]]

## God Nodes (most connected - your core abstractions)
1. `NewAuthService()` - 14 edges
2. `WriteJSON()` - 12 edges
3. `NewInMemoryUserRepository()` - 11 edges
4. `AuthMiddleware()` - 10 edges
5. `NewAuthHandler()` - 10 edges
6. `NewServer()` - 9 edges
7. `AuthService` - 7 edges
8. `InMemoryUserRepository` - 7 edges
9. `MockAuthService` - 7 edges
10. `main()` - 7 edges

## Surprising Connections (you probably didn't know these)
- `AuthMiddleware()` --calls--> `WriteJSON()`  [INFERRED]
  /home/cmd/Desktop/code/go/chat_app/shared/middleware/auth_middleware.go → /home/cmd/Desktop/code/go/chat_app/shared/helper/response.go
- `GuestMiddleware()` --calls--> `WriteJSON()`  [INFERRED]
  /home/cmd/Desktop/code/go/chat_app/shared/middleware/auth_middleware.go → /home/cmd/Desktop/code/go/chat_app/shared/helper/response.go
- `NewServer()` --calls--> `NewEmailVerificationService()`  [INFERRED]
  /home/cmd/Desktop/code/go/chat_app/server/server.go → /home/cmd/Desktop/code/go/chat_app/internal/auth/services/email_verification_service.go
- `NewServer()` --calls--> `NewAuthService()`  [INFERRED]
  /home/cmd/Desktop/code/go/chat_app/server/server.go → /home/cmd/Desktop/code/go/chat_app/internal/auth/services/auth_service.go
- `NewServer()` --calls--> `NewInMemoryUserRepository()`  [INFERRED]
  /home/cmd/Desktop/code/go/chat_app/server/server.go → /home/cmd/Desktop/code/go/chat_app/internal/auth/repositories/inmem_user_repo.go

## Communities

### Community 0 - "Auth Service Tests"
Cohesion: 0.23
Nodes (13): NewAuthService(), TestLogin_InvalidEmail(), TestLogin_Success(), TestLogin_UserNotFound(), TestLogin_WrongPassword(), TestValidateRegisterInput_EmptyFields(), TestValidateRegisterInput_InvalidEmail(), TestValidateRegisterInput_NameTooShort() (+5 more)

### Community 1 - "HTTP Handlers & Routes"
Cohesion: 0.13
Nodes (7): translateServiceError(), isNumeric(), NewEmailVerificationHandler(), AuthHandler, EmailVerificationHandler, WriteJSON(), checkHealth()

### Community 2 - "Server Initialization"
Cohesion: 0.13
Nodes (9): Config, LoadConfig(), Sender, Init(), Sync(), main(), NewSender(), NewServer() (+1 more)

### Community 3 - "Auth Middleware"
Cohesion: 0.23
Nodes (11): AuthMiddleware(), GuestMiddleware(), decodeJSON(), TestAuthMiddleware_EmptyCookieValue(), TestAuthMiddleware_ExpiredToken(), TestAuthMiddleware_InvalidToken(), TestAuthMiddleware_MissingCookie(), TestAuthMiddleware_UserContextKey() (+3 more)

### Community 4 - "Error Handling"
Cohesion: 0.18
Nodes (5): AppError, IsEmailAlreadyExists(), NewEmailAlreadyExists(), AuthService, ValidationError

### Community 5 - "User Repository"
Cohesion: 0.31
Nodes (9): NewInMemoryUserRepository(), TestInMemoryUserRepository_CreateUser_AddsUserToRepository(), TestInMemoryUserRepository_CreateUser_ReturnsErrorForDuplicateEmail(), TestInMemoryUserRepository_EmailExists_ReturnsFalseForNonExistentEmail(), TestInMemoryUserRepository_EmailExists_ReturnsTrueForExistingEmail(), TestInMemoryUserRepository_GetUserByEmail_ReturnsErrorForNonExistent(), TestInMemoryUserRepository_GetUserByEmail_ReturnsUser(), TestInMemoryUserRepository_SeededUsers_ArePresent() (+1 more)

### Community 6 - "Handler Tests"
Cohesion: 0.31
Nodes (11): NewAuthHandler(), TestAuthHandler_Login_InvalidCredentials(), TestAuthHandler_Login_Success(), TestAuthHandler_Login_ValidationError(), TestAuthHandler_Logout_Success(), TestAuthHandler_Register_EmailExists(), TestAuthHandler_Register_Success(), TestAuthHandler_Register_ValidationError() (+3 more)

### Community 7 - "Email Verification"
Cohesion: 0.31
Nodes (2): NewEmailVerificationService(), EmailVerificationService

### Community 8 - "Validation Errors"
Cohesion: 0.67
Nodes (4): cleanFieldName(), translateTag(), TranslateValidationErrors(), ValidationError

### Community 9 - "Domain Models"
Cohesion: 0.53
Nodes (4): LoginInput, RegisterInput, User, ValidationError

### Community 10 - "Token Payload"
Cohesion: 0.5
Nodes (2): NewPayload(), Payload

### Community 11 - "Service Ports"
Cohesion: 0.67
Nodes (2): AuthServicePort, EmailVerificationServicePort

### Community 12 - "Token Models"
Cohesion: 0.67
Nodes (1): TokenUser

### Community 13 - "Repository Port"
Cohesion: 0.67
Nodes (1): UserRepository

### Community 14 - "Project Notes"
Cohesion: 0.67
Nodes (3): Future Change: OTP Logic from DB to Cache, Phase 1: Authentication Package, Phase 1.1: Rate Limiting Setup

## Knowledge Gaps
- **2 isolated node(s):** `Phase 1.1: Rate Limiting Setup`, `Future Change: OTP Logic from DB to Cache`
  These have ≤1 connection - possible missing edges or undocumented components.
- **Thin community `Email Verification`** (10 nodes): `NewEmailVerificationService()`, `email_verification_service.go`, `email_verification_service.go`, `.GetUserByEmail()`, `.UpdateUser()`, `EmailVerificationService`, `.GenerateOTP()`, `.IsVerified()`, `.SendOTP()`, `.VerifyOTP()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Token Payload`** (5 nodes): `payload.go`, `NewPayload()`, `payload.go`, `Payload`, `.Valid()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Service Ports`** (4 nodes): `port.go`, `port.go`, `AuthServicePort`, `EmailVerificationServicePort`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Token Models`** (3 nodes): `models.go`, `models.go`, `TokenUser`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Repository Port`** (3 nodes): `port.go`, `port.go`, `UserRepository`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `NewServer()` connect `Server Initialization` to `Auth Service Tests`, `HTTP Handlers & Routes`, `User Repository`, `Handler Tests`, `Email Verification`?**
  _High betweenness centrality (0.213) - this node is a cross-community bridge._
- **Why does `WriteJSON()` connect `HTTP Handlers & Routes` to `Auth Service Tests`, `Auth Middleware`?**
  _High betweenness centrality (0.142) - this node is a cross-community bridge._
- **Why does `AuthMiddleware()` connect `Auth Middleware` to `HTTP Handlers & Routes`, `Server Initialization`?**
  _High betweenness centrality (0.133) - this node is a cross-community bridge._
- **Are the 12 inferred relationships involving `NewAuthService()` (e.g. with `TestValidateRegisterInput_NameTooShort()` and `TestValidateRegisterInput_InvalidEmail()`) actually correct?**
  _`NewAuthService()` has 12 INFERRED edges - model-reasoned connections that need verification._
- **Are the 10 inferred relationships involving `WriteJSON()` (e.g. with `AuthMiddleware()` and `GuestMiddleware()`) actually correct?**
  _`WriteJSON()` has 10 INFERRED edges - model-reasoned connections that need verification._
- **Are the 8 inferred relationships involving `NewInMemoryUserRepository()` (e.g. with `TestInMemoryUserRepository_EmailExists_ReturnsFalseForNonExistentEmail()` and `TestInMemoryUserRepository_EmailExists_ReturnsTrueForExistingEmail()`) actually correct?**
  _`NewInMemoryUserRepository()` has 8 INFERRED edges - model-reasoned connections that need verification._
- **Are the 8 inferred relationships involving `AuthMiddleware()` (e.g. with `TestAuthMiddleware_MissingCookie()` and `TestAuthMiddleware_ValidToken()`) actually correct?**
  _`AuthMiddleware()` has 8 INFERRED edges - model-reasoned connections that need verification._