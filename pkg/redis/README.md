
Unit testing can use https://github.com/alicebob/miniredis, you can start a local mock redis

- [Mock Redis in unit tests](https://medium.com/@elliotchance/mocking-redis-in-unit-tests-in-go-28aff285b98)

## case

- [Redis distributed locks are useless to understand, and a big failure has occurred...](https://mp.weixin.qq.com/s/BO-gly5iGLVmuG5B_FIpoQ)
- [After reading this three major problems of Redis caching](https://mp.weixin.qq.com/s/HjzwefprYSGraU1aJcJ25g)

## Redis Optimization direction

### Parameter optimization

MaxIdle is set to a high point, which can ensure that there are enough connections to obtain redis in the case of burst
traffic, and there is no need to establish a connection in high traffic conditions

**go-redis parameter optimization**

```yaml
  min_idle_conn: 30               
  dial_timeout: "1s"
  read_timeout: "500ms"
  write_timeout: "500ms"
  pool_size: 500
  pool_timeout: "60s"
```

**redisgo parameter optimization**

```yaml
maxIdle = 30
maxActive = 500
dialTimeout = "1s"
readTimeout = "500ms"
writeTimeout = "500ms"
idleTimeout = "60s"
```

### use optimization

- Add redis slave library
- Increase the batch data of redis slave libraries, and pull data concurrently with goroutines according to the number 
  of redis slave libraries
- Extensive use of pipeline directives for bulk data
- Reduce key fields
- The value storage codec of redis is changed to msgpack

## Pipeline
- https://redis.io/topics/pipelining
- [Pipeline batch compatible with go redis cluster](http://xiaorui.cc/archives/5557)
- https://www.tizi365.com/archives/309.html