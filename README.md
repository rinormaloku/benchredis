## Run benchmark for redis

This repo tests whether Redis is able to handle up to 1 million secrets
and still return sub milliseconds responses.

Run redis 
```bash
docker run -d --name perf-redis -p 6379:6379 redis
```

Run tests
```
go build -o benchredis
./benchredis | tee results.log
```
Stop and remove redis
```bash
docker rm -f perf-redis
```

## Results

**Even at 1m api keys vault is able to serve requests in sub milliseconds.**

Measurements:

```bash
Average response time: 23.139000 microseconds, errors: 0
Total errors adding 1000000 keys was: 0

Memory: 593MiB
```