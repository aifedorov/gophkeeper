// Package api provides the API layer for the GophKeeper binary.
//
// This package contains the gRPC service definitions and generated code
// for handling client requests. It includes:
//
//   - Protocol Buffer definitions (.proto files) for service contracts
//   - Generated Go code from protobuf definitions
//   - Service interface definitions for authentication and credential management
//
// # Services
//
// ## AuthService
//
// Handles user authentication operations:
//   - Register: Create new user accounts with encrypted credentials
//   - login: Authenticate existing users and issue JWT tokens
//
// ## CredentialService
//
// Manages user credentials with end-to-end encryption:
//   - Create: Store new credentials with encryption
//   - List: Retrieve and decrypt a specific credential
//   - Update: Modify existing credentials
//   - Delete: Soft delete credentials
//   - List: Retrieve all user credentials
//
// # Security
//
// All credential operations require:
//   - Valid JWT token in request metadata (key: "access_token")
//   - User authentication via AuthService
//   - End-to-end encryption of sensitive data (passwords, login, metadata)
//
// # Code Generation
//
// The gRPC code is generated using buf:
//
//	cd internal/binary/api/grpc
//	buf generate
//
// Or use the Makefile:
//
//	make proto-gen
//
// # Directory Structure
//
//	api/
//	└── grpc/
//	    ├── proto/              # Protocol Buffer definitions
//	    │   ├── auth/v1/        # Authentication service proto
//	    │   └── credential/v1/  # Credential service proto
//	    ├── gen/                # Generated Go code
//	    │   ├── auth/v1/        # Generated auth service code
//	    │   └── credential/v1/  # Generated credential service code
//	    ├── buf.yaml            # Buf configuration
//	    └── buf.gen.yaml        # Buf generation config
//
// # Error Handling
//
// The API layer uses standard gRPC status codes:
//   - OK: Successful operation
//   - InvalidArgument: Invalid request parameters
//   - Unauthenticated: Missing or invalid JWT token
//   - AlreadyExists: Resource already exists (e.g., duplicate login)
//   - NotFound: Resource not found
//   - Internal: Unexpected binary error
//
// # Example Usage
//
// Client code example:
//
//	// Connect to binary
//	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(...))
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer conn.Close()
//
//	// Create auth client
//	authClient := authv1.NewAuthServiceClient(conn)
//
//	// Register new user
//	resp, err := authClient.Register(ctx, &authv1.RegisterRequest{
//	    login:    "user@example.com",
//	    password: "secure-password",
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Use access token for authenticated requests
//	md := metadata.Pairs("access_token", resp.GetAccessToken())
//	ctx = metadata.NewOutgoingContext(ctx, md)
//
//	// Create credential client
//	credClient := credentialv1.NewCredentialServiceClient(conn)
//
//	// Create new credential
//	credResp, err := credClient.Create(ctx, &credentialv1.CreateRequest{
//	    Name:     "Gmail",
//	    login:    "user@gmail.com",
//	    password: "gmail-password",
//	    Metadata: "Personal email account",
//	})
package api
