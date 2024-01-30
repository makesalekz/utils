# utils

## api

Proto models

## v1

### config

Consul & Vault config readers

### dialer

Deprecated gRPC dialer

### jwt

Working with JWT claims

### log

JSON logger

### middlewares/auth

JWT server & client middlewares

#### middlewares/metrics

Custom metrics middlewares

### nats

QueueManager to handle simple nats queues

## v2

### auth

Getters for ActorId and TenantId from metadata

### v2/dialer

New gRPC dialer – 1 connection for each s2s link. Uses s2s token & metadata

### v2/middlewares/auth

Refactored JWT server & client middlewares:

- JWT Server – unchanged
- JWT Client – sends s2s token instead of user ID token
- BFF Meta Server – appends actor ID & tenant ID as metadata to context, to send them in s2s calls
