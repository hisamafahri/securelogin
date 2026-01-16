# SecureLogin

OAuth 2.0 and OpenID Connect server (almost RFC-Compliant) that provides authentication services using third-party identity providers.

SecureLogin acts as an OAuth/OIDC provider that delegates authentication to external identity providers (Google, GitHub, Microsoft). It allows applications to implement unified "Sign in with Google/GitHub/Microsoft" functionality through a single interface.

## Who is this for?

It's for developers/teams that want to centralize authentication service for their applications.

If your team is building a new application and wants to implement "Sign in with Google/GitHub/Microsoft" without building your own authentication system, you can use this server as a standalone authentication service.

## Features

- OAuth 2.0 Authorization Code Flow with PKCE
- OpenID Connect ID tokens
- Refresh tokens with `offline_access` scope
- JWT-based access and ID tokens
- Multi-provider authentication (Google, GitHub, Microsoft)
- Session management with token revocation
- Standard OIDC discovery endpoints

## Available Commands

- `make dev` - Run with hot reload
- `make format` - Format code
- `make migrate-up` - Apply migrations
- `make migrate-down` - Rollback migrations
- `make migrate-generate` - Generate new migration
- `make migrate-status` - Check migration status
- `make migrate-hash` - Regenerate atlas.sum

## API Endpoints

- `GET /ping` - Health check
- `GET /authorize` - OAuth authorization endpoint
- `GET /signin` - Sign-in page
- `POST /signin/identifier` - Provider selection
- `GET /callback/google` - Google OAuth callback
- `GET /callback/github` - GitHub OAuth callback
- `GET /callback/microsoft` - Microsoft OAuth callback
- `POST /oauth/token` - Token exchange (authorization_code, refresh_token)
- `POST /oauth/revoke` - Token revocation
- `GET /.well-known/openid-configuration` - OIDC discovery
- `GET /.well-known/jwks.json` - JSON Web Key Set
- `GET /oauth/userinfo` - User information (requires access token)

## Configuration

### Environment Variables

- `DATABASE_URL` - PostgreSQL connection string (required)
- `REDIS_URL` - Redis connection string
- `SYSTEM_BASE_URL` - OAuth issuer URL (required)
- `JWT_KEY_STORE_PATH` - JWT signing keys storage path

### Database Setup

Before using the server, configure:

1. Register applications with client_id, client_secret, and redirect_uris
2. Configure authentication providers (Google/GitHub/Microsoft credentials) per application

### JWT Key Pairs

Public-private key pairs for JWT signing are stored in `./data`. The keys are auto-generated on first run. Users can adjust the storage path via the `JWT_KEY_STORE_PATH` environment variable as needed.

## In Roadmap

- Traditional authentication methods (email/phone/username + password)
- Passwordless authentication (magic links, one-time passcodes)
- MFA (Passkeys & TOTP)
- User management (registration, profile updates, account recovery)
- Authorization (role-based access control, permissions)
- OAuth Device Flow
- At rest encryption for client secrets and tokens
- Support for more identity providers (Facebook, Twitter, LinkedIn, etc.)
- Admin interface for managing clients and providers
- JWKS rotation and key management
