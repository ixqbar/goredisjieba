<?php

$redis_handle = new Redis();
$redis_handle->connect('127.0.0.1', 6479);
$redis_handle->select(0);

$result = $redis_handle->rawCommand('cutforsearch', '我来到北京清华大学', 1);
print_r($result);
