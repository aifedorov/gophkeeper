// Package api provides the gRPC API layer for the GophKeeper server.
//
// This package contains the gRPC service definitions and generated code
// for handling client requests. It includes:
//
//   - Protocol Buffer definitions (.proto files) for service contracts
//   - Generated Go code from protobuf definitions
//   - Service implementations for authentication, credentials, and binary file management
//
// # Services
//
// ## AuthService
//
// Handles user authentication operations. No authentication required.
//
// Methods:
//   - Register: Create a new user account with encrypted credentials.
//     Returns JWT access token, encryption key (base64), and user ID.
//     Error codes: AlreadyExists (login taken), InvalidArgument (invalid credentials)
//   - Login: Authenticate an existing user and issue JWT tokens.
//     Returns JWT access token, encryption key (base64), and user ID.
//     Error codes: Unauthenticated (invalid credentials), InvalidArgument (invalid input)
//
// Request/Response:
//   - RegisterRequest: login (string), password (string)
//   - RegisterResponse: access_token (string), encryption_key (bytes), user_id (string)
//   - LoginRequest: login (string), password (string)
//   - LoginResponse: access_token (string), encryption_key (bytes), user_id (string)
//
// ## CredentialService
//
// Manages user credentials (login/password pairs) with end-to-end encryption.
// All operations require authentication via JWT token in metadata.
//
// Methods:
//   - Create: Store a new credential with encryption.
//     Returns the ID of the newly created credential.
//     Error codes: AlreadyExists (name already exists), InvalidArgument (missing/invalid fields)
//   - List: Retrieve all credentials for the authenticated user.
//     Returns a list of decrypted credentials.
//     Error codes: Unauthenticated (invalid token)
//   - Update: Modify an existing credential.
//     Returns success status.
//     Error codes: NotFound (credential not found), InvalidArgument (missing/invalid fields)
//   - Delete: Delete a credential for the authenticated user.
//     Returns success status.
//     Error codes: NotFound (credential not found)
//
// Request/Response:
//   - CreateRequest: name (string, required), login (string, required), password (string, required), notes (string, optional)
//   - CreateResponse: id (string)
//   - ListRequest: (empty)
//   - ListResponse: credentials (repeated ListItem with id, name, login, password, notes)
//   - UpdateRequest: id (string, required), name (string, required), login (string, required), password (string, required), notes (string, optional)
//   - UpdateResponse: success (bool)
//   - DeleteRequest: id (string, required)
//   - DeleteResponse: success (bool)
//
// ## BinaryService
//
// Manages binary file storage with streaming support for large files.
// All operations require authentication via JWT token in metadata.
// Files are encrypted using AES-256-GCM before storage.
//
// Methods:
//   - Upload: Upload a binary file using client streaming.
//     First message contains metadata (name, size, notes), subsequent messages contain file chunks.
//     Returns the file ID upon successful upload.
//     Error codes: InvalidArgument (invalid metadata), Internal (upload failure)
//   - List: Retrieve metadata for all files belonging to the authenticated user.
//     Returns a list of file metadata (id, name, size, notes, uploaded_at).
//     Error codes: Unauthenticated (invalid token)
//   - Download: Download a binary file using server streaming.
//     First message contains file metadata, subsequent messages contain encrypted file chunks.
//     Error codes: NotFound (file not found), Unauthenticated (invalid token)
//   - Delete: Delete a file for the authenticated user.
//     Returns success status.
//     Error codes: NotFound (file not found)
//
// Request/Response:
//   - UploadRequest (stream): oneof { Metadata (name, size, notes) | chunk (bytes) }
//   - UploadResponse: file_id (string)
//   - ListRequest: (empty)
//   - ListResponse: files (repeated MetadataResponse with id, name, size, notes, uploaded_at)
//   - DownloadRequest: file_id (string, required)
//   - DownloadResponse (stream): oneof { Metadata (name, size, notes) | chunk (bytes) }
//   - DeleteRequest: file_id (string, required)
//   - DeleteResponse: (empty)
//
// # Authentication
//
// All operations except AuthService.Register and AuthService.Login require authentication.
//
// Authentication is performed via JWT token passed in gRPC metadata:
//
//	md := metadata.Pairs("access_token", "<jwt-token>")
//	ctx := metadata.NewOutgoingContext(ctx, md)
//
// The JWT token is validated by the authentication interceptor, which extracts:
//   - User ID: used to scope all operations to the authenticated user
//   - Encryption key: used to encrypt/decrypt user data (stored in token claims)
//
// # Security
//
// Security features:
//   - End-to-end encryption: All sensitive data (passwords, credentials, files) are encrypted
//     using AES-256-GCM before storage
//   - Password hashing: User passwords are hashed using bcrypt before storage
//   - JWT tokens: Stateless authentication using signed JWT tokens
//   - TLS: gRPC connections should use TLS in production
//   - User isolation: All operations are scoped to the authenticated user
//
// Encryption:
//   - Credentials: login, password, and notes are encrypted individually
//   - Binary files: Files are encrypted in chunks during upload/download
//   - Encryption key: Derived from user password using Argon2, stored in JWT token
//
// # Error Handling
//
// The API uses standard gRPC status codes:
//   - OK: Successful operation
//   - InvalidArgument: Invalid request parameters (e.g., empty required fields, invalid format)
//   - Unauthenticated: Missing or invalid JWT token, or invalid credentials
//   - AlreadyExists: Resource already exists (e.g., duplicate login or credential name)
//   - NotFound: Resource not found (e.g., credential or file doesn't exist)
//   - Internal: Unexpected server error
//
// Error messages are returned in the gRPC status details and should be checked
// using status.FromError().
//
// # Code Generation
//
// The gRPC code is generated using buf:
//
//	cd internal/server/api/grpc
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
//	    │   ├── credential/v1/  # Credential service proto
//	    │   └── binary/v1/      # Binary file service proto
//	    ├── gen/                # Generated Go code
//	    │   ├── auth/v1/        # Generated auth service code
//	    │   ├── credential/v1/ # Generated credential service code
//	    │   └── binary/v1/      # Generated binary service code
//	    ├── buf.yaml            # Buf configuration
//	    └── buf.gen.yaml        # Buf generation config
//
// # Example Usage
//
// Complete client example:
//
//	// 1. Connect to server
//	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(creds))
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer conn.Close()
//
//	// 2. Register new user
//	authClient := authv1.NewAuthServiceClient(conn)
//	registerResp, err := authClient.Register(ctx, &authv1.RegisterRequest{
//	    Login:    "user@example.com",
//	    Password: "secure-password",
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// 3. Use access token for authenticated requests
//	md := metadata.Pairs("access_token", registerResp.GetAccessToken())
//	ctx = metadata.NewOutgoingContext(ctx, md)
//
//	// 4. Create credential
//	credClient := credentialv1.NewCredentialServiceClient(conn)
//	credResp, err := credClient.Create(ctx, &credentialv1.CreateRequest{
//	    Name:     "Gmail",
//	    Login:    "user@gmail.com",
//	    Password: "gmail-password",
//	    Notes:    "Personal email account",
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// 5. List all credentials
//	listResp, err := credClient.List(ctx, &credentialv1.ListRequest{})
//	if err != nil {
//	    log.Fatal(err)
//	}
//	for _, cred := range listResp.GetCredentials() {
//	    fmt.Printf("ID: %s, Name: %s, Login: %s\n", cred.GetId(), cred.GetName(), cred.GetLogin())
//	}
//
//	// 6. Upload binary file
//	binaryClient := binaryv1.NewBinaryServiceClient(conn)
//	uploadStream, err := binaryClient.Upload(ctx)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Send metadata first
//	err = uploadStream.Send(&binaryv1.UploadRequest{
//	    Data: &binaryv1.UploadRequest_File{
//	        File: &binaryv1.UploadRequest_Metadata{
//	            Name:  "document.pdf",
//	            Size:  1024,
//	            Notes: "Important document",
//	        },
//	    },
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Send file chunks
//	file, _ := os.Open("document.pdf")
//	buf := make([]byte, 64*1024) // 64KB chunks
//	for {
//	    n, err := file.Read(buf)
//	    if err == io.EOF {
//	        break
//	    }
//	    uploadStream.Send(&binaryv1.UploadRequest{
//	        Data: &binaryv1.UploadRequest_Chunk{
//	            Chunk: buf[:n],
//	        },
//	    })
//	}
//
//	uploadResp, err := uploadStream.CloseAndRecv()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Uploaded file ID: %s\n", uploadResp.GetFileId())
//
// # Streaming
//
// BinaryService uses gRPC streaming for efficient file transfer:
//
//   - Upload: Client streaming - client sends metadata first, then file chunks
//   - Download: Server streaming - server sends metadata first, then file chunks
//
// Recommended chunk size: 64KB - 1MB for optimal performance.
// The server buffers up to 1MB per request.
//
// # Rate Limiting and Limits
//
// Current limits (may vary by deployment):
//   - Maximum file size: 10GB (defined in binary service)
//   - Chunk size: Recommended 64KB - 1MB
//   - Buffer size: 1MB per streaming request
//
// # Versioning
//
// All services are versioned (v1). Future breaking changes will be introduced
// in new versions (v2, v3, etc.) while maintaining backward compatibility.
package api
