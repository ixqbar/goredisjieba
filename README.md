### version 0.0.2

### config.xml
```
<?xml version="1.0" encoding="UTF-8" ?>
<config>
    <address>0.0.0.0:6379</address>
    <db>0</db>
    <dict>/data/dict</dict>
</config>
```
* 其中db为dict定义目录下的一个子目录，所有字典存在该子目录下,当使用select db时可切换分词字典

### usage
```
make linux

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
127.0.0.1:6379>extract 我来到北京清华大学  20
清华大学
来到
北京
127.0.0.1:6379>
```
* dependent on https://github.com/yanyiwu/gojieba 

更多疑问请+qq群 233415606
