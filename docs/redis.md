# Redis with Docker Compose

Multienv needs a running Redis with [JSON commands](requirements.md#redis).
This example uses the [official Redis image](https://hub.docker.com/_/redis) and Docker Compose v2.
All files below live outside the checkout, for example in `~/.config/multienv/redis`.

## Configuration files

```sh
mkdir -p ~/.config/multienv/redis
cd ~/.config/multienv/redis
```

Create `compose.yaml`:

```yaml
services:
  redis:
    image: redis:8.2-alpine
    command: [redis-server, /usr/local/etc/redis/redis.conf]
    ports:
      - "127.0.0.1:6379:6379"
    volumes:
      - ./redis.conf:/usr/local/etc/redis/redis.conf:ro
      - redis-data:/data
    restart: unless-stopped

volumes:
  redis-data:
```

Create `redis.conf`, replacing the example password with your own secret:

```text
bind 0.0.0.0
protected-mode yes
port 6379
dir /data
appendonly yes
appendfsync everysec
requirepass replace-with-a-long-random-secret
loadmodule /usr/local/lib/redis/modules/rejson.so
```

The module path belongs to the image above. Explicitly loading it enables the JSON commands used by Multienv;
Redis configuration and module loading are described in the
[Redis configuration reference](https://redis.io/docs/latest/operate/oss_and_stack/management/config/).
The container listens on its network interface, while the published host port is restricted to loopback.
The named volume preserves data when the container is replaced.

## Start and check

From the directory containing `compose.yaml`:

```sh
docker compose up -d
docker compose logs --tail 30 redis
docker compose exec redis redis-cli --askpass PING
docker compose exec redis redis-cli --askpass COMMAND INFO JSON.GET JSON.SET JSON.DEL
```

Enter the password from `redis.conf`. Expect `PONG` and nonempty descriptions for all three JSON commands.

## Connect Multienv

Create `~/.config/multienv/redis-credentials.yaml` using the same password:

```yaml
username: default
password: replace-with-a-long-random-secret
```

`requirepass` sets the password of Redis's default user; see the
[Redis authentication reference](https://redis.io/docs/latest/operate/oss_and_stack/management/security/acl/).

Set the connection in `server.yaml`:

```yaml
storage:
  address: 127.0.0.1
  port: 6379
  db: 0
  credentials_path: ~/.config/multienv/redis-credentials.yaml
```

Start the [server](server.md#start), then run `multienv-agent doctor global` from a configured host.

## Stop

```sh
docker compose down
```

This keeps the named volume. `docker compose down -v` also deletes the stored Redis data.
