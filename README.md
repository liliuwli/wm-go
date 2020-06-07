# wm-go
  这是php开发者用workerman思路开发得一款tcp/udp异步服务框架，目前只实现部分功能  
  
+支持自定义协议(默认协议为telnet)
+支持多goruntime工作，提供handle/onconnect/onclose等api，让开发者可以专注实现业务
+未支持守护进程
+未支持日志模块
+未支持异步发起连接
+未支持毫秒级定时器
+未支持分布式构建