# Redis Purger
Simple script to delete specific keys of a Redis cluster.

## Usage
```shell
> go build -o redis-purger .
> ./redis-purger --cluster --hosts="host1.redis.com;host2.redis.com" --key="myKey*"
```
