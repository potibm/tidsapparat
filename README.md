# Tidsapparat

![tidsapparat logo](doc/tidsapparat.svg)

> _Tidsapparart_ is a Danish term for time device (even a time machine).

It is a editor for timetables at [demoparties](https://en.wikipedia.org/wiki/Demoscene#Parties). Can be used in conjunction with [billedapparat](https://github.com/potibm/billedapparat).

## Tooling

- **Backend**
  - [Go](https://go.dev)
  - [Gin Web Framework](https://gin-gonic.com)
  - [GORM](https://gorm.io)
  - [Cobra](https://cobra.dev) & [Viper](https://github.com/spf13/viper)
- **Frontend**
  - [React](https://react.dev)
  - [Vite](https://vitejs.dev/)
  - [React Admin](https://marmelab.com/react-admin/)
- **Database**
  - [SQLite](https://www.sqlite.org)
- **Identity & Local Infrastructure**
  - [Dex](https://dexidp.io/) (OIDC Provider)
  - [Traefik](https://traefik.io/) (Local Edge Router)
  - [mkcert](https://github.com/FiloSottile/mkcert) & dnsmasq (Local TLS & `.test` Routing)
- **Observability**
  - [Sentry](https://sentry.io)
  - [OpenTelemetry](https://opentelemetry.io)
- **Development & Ops**
  - [mise](https://mise.jdx.dev/)
  - [Docker](https://www.docker.com)

## Quickstart

We use `mise` to automatically manage all tool versions (Go, Node, etc.) and project tasks.

```bash
# 1. Install mise (if not already installed)
curl https://mise.run | sh

# 2. Setup local infrastructure
# This generates local certificates, configures dnsmasq/resolver, and updates /etc/hosts (Linux)
mise run infra:prepare

# 3. Start local services (Traefik, Dex, etc.)
mise run infra:up

# 3. Start the development server (hot-reload for backend & frontend)
overmind s --timeout 10
```

## Local Environment

Once the stack is running, you can access the applications via:

- Tidsapparat Admin: https://tidsapparat.test
- Dex IdP: https://dex.tidsapparat.test

## Documentation

Please refer to the [configuration](doc/configuration.md) and [sync](doc/sync.md) documentation. Further information will be added soon.
