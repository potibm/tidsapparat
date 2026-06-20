# Tidsapparat

## OIDC Authentication

Tidsapparat supports OIDC-based authentication for the admin panel. The backend validates JWT Bearer tokens against a configurable OIDC provider (e.g., Dex, Keycloak, Auth0), and the React Admin frontend handles the login flow via `oidc-client-ts`.

### Configuration

Authentication is configured via the `auth` section in `backend/config/config.yaml` (or environment variables):

```yaml
auth:
  type: "oidc"
  name: "Dex"
  authority: "https://dex.tidsapparat.test/dex"
  client_id: "react-admin-client"
  skip_tls_verify: false
```

| Key | Environment Variable | Description |
|-----|---------------------|-------------|
| `auth.type` | `AUTH_TYPE` | Authentication type. Currently only `oidc` is supported. |
| `auth.name` | `AUTH_NAME` | Display name of the Identity Provider (shown on the login button). |
| `auth.authority` | `AUTH_AUTHORITY` | OIDC Issuer URL (must expose `/.well-known/openid-configuration`). |
| `auth.client_id` | `AUTH_CLIENT_ID` | OIDC Client ID registered with the provider. |
| `auth.skip_tls_verify` | `AUTH_SKIP_TLS_VERIFY` | **Development only.** Disables TLS certificate verification for the OIDC discovery endpoint. Must be `false` in production. |

### Authentication Flow

1. **Frontend**: When `auth.type` is `oidc`, the admin app bootstraps a `UserManager` from `oidc-client-ts` using the authority and client ID from `/api/config`.
2. **Login**: The user clicks "Login with {Name}" on the custom login page. This triggers `signinRedirect()` to the OIDC provider.
3. **Callback**: After successful authentication, the provider redirects back to `https://tidsasapparat.test/auth-callback`. The `authProvider.handleCallback()` exchanges the authorization code for tokens.
4. **API Calls**: The `dataProvider` automatically attaches the access token as a `Bearer` header to every `/api/admin` request.
5. **Backend Validation**: The `AuthMiddleware` verifies the JWT signature and issuer against the configured OIDC provider. Valid tokens inject the `userID` into the Gin context and GORM request context.
6. **Audit Trail**: The GORM audit callbacks automatically populate `created_by`, `modified_by`, and `deleted_by` fields on all models that embed `AuditModel`.

### Default Local Credentials

When using the bundled Dex instance, the default user is:

- **Email**: `admin@example.com`
- **Password**: `password`
- **User ID**: `08a8684b-db88-4b73-90a9-3cd1661f5466`

### Disabling Authentication

To run the admin panel without authentication (not recommended for production), omit the `auth` section from the configuration entirely. The frontend will then skip OIDC bootstrapping and the backend will not enforce Bearer token validation on `/api/admin` routes.
