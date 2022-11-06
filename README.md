# Secret pass

Secret-pass is a simple, self-hosted service that allows you to securely get a secret from someone.

## Parameters

| Command line | Environment variable | Default | Description                 |
| ------------ | -------------------- | ------- | --------------------------- |
| dbg          | DEBUG                | false   | Enable debug mode           |
| lifespan, ls | LIFESPAN             |         | Secret lifespan, _required_ |
| port, p      | PORT                 | :8080   | Service port                |