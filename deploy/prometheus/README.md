
By starting prometheus to collect various running indicators of the service, it is convenient for monitoring.
start a prometheus

```bash
prometheus --config.file=prometheus.yml
```

You can check the operation of the service through http://localhost:8090/ping

## Reference

- https://github.com/yolossn/Prometheus-Basics