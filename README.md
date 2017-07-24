### version 0.0.1

### usage
```
make

./bin/goRedisJieba_linux --config=config.xml
```

### command
```
redis-cli --raw
127.0.0.1:6379> tag 我来到北京清华大学
我/r
来到/v
北京/ns
清华大学/nt
127.0.0.1:6379> cut 我来到北京清华大学 0
我
来到
北京
清华大学
127.0.0.1:6379> cut 我来到北京清华大学 1
我
来到
北京
清华大学
127.0.0.1:6379> cutforsearch 我来到北京清华大学 1
我
来到
北京
清华
华大
大学
清华大学
127.0.0.1:6379> cutforsearch 我来到北京清华大学 0
我
来到
北京
清华
华大
大学
清华大学
127.0.0.1:6379>
```

更多疑问请+qq群 233415606