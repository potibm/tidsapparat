# Tidsapparat

![tidsapparat logo](doc/tidsapparat.svg)

> _Tidsapparart_ is a Danish term for time device (even a time machine).

It is a editor for timetables at [demoparties](https://en.wikipedia.org/wiki/Demoscene#Parties). Can be used in conjunction with [billedapparat](https://github.com/potibm/billedapparat).

## Tooling

- [Go](https://go.dev)
  - [Gin Web Framework](https://gin-gonic.com)
  - [GORM](https://gorm.io)
  - [Cobra](https://cobra.dev) & [Viper](https://github.com/spf13/viper)
- [React](https://react.dev)
  - [Vite](https://vitejs.dev/)
  - [React Admin](https://marmelab.com/react-admin/)
  - [Flowbite React](https://flowbite-react.com) & [Tailwind CSS](https://tailwindcss.com)
- [SQLite](https://www.sqlite.org)
- Observability
  - [Sentry](https://sentry.io)
  - [OpenTelemetry](https://opentelemetry.io)
- Development & Ops
  - [mise](https://mise.jdx.dev/)
  - [Docker](https://www.docker.com)

## Quickstart

We use `mise` to automatically manage all tool versions (Go, Node, etc.) and project tasks.

```bash
# 1. Install mise (if not already installed)
curl https://mise.run | sh

# 2. Setup the project (installs dependencies and starts infra)
mise run setup

# 3. Start the development server (hot-reload for backend & frontend)
mise run dev
```

## Documentation

_todo_
