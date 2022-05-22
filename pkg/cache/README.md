Here are some libraries mainly about caching, mainly based on in-memory and NoSQL

Memory types are: memory and big_cache
The main ones of NoSQL are: redis

All kinds of libraries only need to implement the interface (driver) defined by the cache.
> The naming of the interface driver here refers to the naming convention of the official mysql interface of Go

## multilevel cache

### L2 cache
The multi-level here mainly refers to the second-level cache: local cache (first-level cache L1) + redis cache (second-level cache L2)
Using a local cache can reduce the network I/O overhead between the application server and redis

> It should be noted that in a system with a small amount of concurrency, the local cache is of little significance, but increases the difficulty of maintenance. But in a high concurrency system,
> Local caching can greatly save bandwidth. But be aware that local caching is not a silver bullet, it will cause data corruption between multiple copies
> Inconsistency will also occupy a large amount of memory, so it is not suitable for saving particularly large data, and the refresh mechanism needs to be strictly considered.

### Expiration

The expiration time of the local cache is at least half smaller than that of the distributed cache, to prevent the local cache from being too long and causing inconsistency of multi-instance data.

## cache problem

The following issues should be noted

- cache penetration
- cache breakdown
- Cache Avalanche

You can refer to the article：[Three major problems of Redis cache](https://mp.weixin.qq.com/s/HjzwefprYSGraU1aJcJ25g)

## Reference
- ristretto：https://github.com/dgraph-io/ristretto (The fastest local cache)
- [Ristretto简介：High performance Go cache](https://www.yuque.com/kshare/2020/ade1d9b5-5925-426a-9566-3a5587af2181)
- bigcache: https://github.com/allegro/bigcache
- freecache: https://github.com/coocood/freecache
- concurrent_map: https://github.com/easierway/concurrent_map
- gocache: https://github.com/eko/gocache (Built-in stores, eg: bigcache,memcache,redis)