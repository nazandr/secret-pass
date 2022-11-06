# Secret pass

Secret-pass is a simple, self-hosted service that allows you to securely get a secret from someone.
You can try in on [Secret-pass](https://nniel.site)

## Parameters

| Command line | Environment variable | Default | Description                 |
| ------------ | -------------------- | ------- | --------------------------- |
| dbg          | DEBUG                | false   | Enable debug mode           |
| lifespan, l  | LIFESPAN             |         | Secret lifespan, _required_ |
| port, p      | PORT                 | :8080   | Service port                |