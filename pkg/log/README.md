# illustrate

This logger mainly encapsulates the zap library, which is easy to use. Of course, other logging libraries, such as `logrus`, can also be used according to the interface.

## log function

- Record log information to log file
- Log Slicing - Ability to slice based on log file size or time interval
- Support different log levels (eg: info, debug, warn, error, fatal)
- Support classification and output to different log files by log level
- Ability to print basic information such as calling file/function name and line number, log time, IP, etc.

## Instructions

```go
log.Info("user_id is 1")
log.Warn("user is not exist")
log.Error("params error")

log.Warnf("params is empty")
...
```

## in principle

Try not to print out the log in model, repository, and service. It is better to use `errors.Wrapf` to return errors and messages to the upper layer, and then handle errors in the handler layer.
That is, it is printed out through the log.

The advantage of this is that it avoids printing the same log in multiple places, making troubleshooting easier.

## Reference

 - Log base library zap: https://github.com/uber-go/zap
 - 日Log Splitting Library - By Time：https://github.com/lestrrat-go/file-rotatelogs
 - Log Splitting Library - By Size：https://github.com/natefinch/lumberjack 
 - [Depth | See how to implement high-performance Go components from the Go high-performance logging library zap](https://mp.weixin.qq.com/s/i0bMh_gLLrdnhAEWlF-xDw)
 - [Logger interface for GO with zap and logrus implementation](https://www.mountedthoughts.com/golang-logger-interface/)
 - https://github.com/wzyonggege/logger
 - https://wisp888.github.io/golang-iris-%E5%AD%A6%E4%B9%A0%E4%BA%94-zap%E6%97%A5%E5%BF%97.html
 - https://www.mountedthoughts.com/golang-logger-interface/
