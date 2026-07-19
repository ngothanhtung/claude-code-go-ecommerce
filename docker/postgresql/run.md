# How to run PostgreSQL Docker container

To run a PostgreSQL Docker container, you can use the following command:

```bash
docker-compose up -d
```

If you previously started the container with a different major PostgreSQL version, remove the old volume first so the new image can initialize a fresh data directory.
