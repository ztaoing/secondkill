# 秒杀业务系统

###### sk-app：是秒杀业务系统，主要是接收用户的秒杀请求，处理用户黑白名单等，然后将秒杀请求通过redis发送给sk-core系统

# 流程介绍
普通用户进行秒杀时，首先与秒杀业务系统进行交互，秒杀业务系统主要负责对请求进行限流、用户黑白名单过滤、并发限制和用户数据签名校验。

秒杀业务系统的工作流程如下：
秒杀活动和秒杀商品的信息存储在zookeeper中，并且可以使用zookeeper的watch机制实时更新信息。

1. 从zookeeper中加载秒杀活动数据到内存中
2. 监听zookeeper中的数据变化，实时更新缓存在内存中秒杀活动数据
3. 从redis中加载黑白名单数据到内存中
4. 设置白名单
5. 对用户请求进行黑名单限制
6. 对用户请求进行流量限制、秒级限制、分级限制。
7. 将用户数据进行签名校验、检验参数的合法性
8. 将用户请求通过redis传递给业务核心系统进行处理
9. 接收秒杀核心系统返回的秒杀处理结果，并返回给用户

