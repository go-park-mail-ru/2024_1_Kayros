# Homework 3

По нашему мнению, основная сущность нашего приложения - это заказы. Мы проведем нагрузочное тестирование: откроем страницу ресторана 1000000 раз и добавим 1000000 раз еду в заказ (для этого уберем ограничение на беке). Для проведения нагрузочного тестирования будем использовать [wrk](https://github.com/wg/wrk).

### Открытие страницы ресторана 

Для выполнения этого нагрузочного теста напишем скрипт на Lua:

`
    wrk.method = "GET" 
    wrk.headers["Content-Type"] = "application/json"
    wrk.headers["XCSRF-Token"] = "380869aba9e630b7625b3226164f387552185e9e69f9fdff62ae008febcf9d23&ff361a24-621c-4854-b416-fcbe6b1122f2!ZE8T8yEW"
    wrk.headers["Cookie"] = "session_id=b09ccf30-d099-480f-aa8c-71a04a049d44;csrf_token=c13ca394c9e3f18b1c75ec35e6fa8dd90db3aeac6f96fd08b290d5e3001ae58e.b09ccf30-d099-480f-aa8c-71a04a049d44!CYhR5ej0;request = function()
        return wrk.format(nil, "/api/v1/restaurants/4", nil, nil)
    end
`
### Скрипт для запуска

` wrk -t4 -c100 -d600s --script=open_rest.lua --latency https://resto-go.online `

### Результат 
`
    Running 10m test @ https://resto-go.online
    4 threads and 100 connections
    Thread Stats   Avg      Stdev     Max   +/- Stdev
        Latency    69.60ms  141.22ms   2.00s    94.07%
        Req/Sec   655.23    347.09     1.52k    62.81%
    Latency Distribution
        50%   31.66ms
        75%   53.22ms
        90%  123.71ms
        99%  815.50ms
    1461609 requests in 10.00m, 540.83MB read
    Socket errors: connect 0, read 0, write 0, timeout 638
    Requests/sec:   2435.85
    Transfer/sec:      0.90MB
`

### Добавление еды в заказ

Для выполнения этого нагрузочного теста напишем скрипт на Lua:

`
wrk.method = "PUT"
wrk.headers["Content-Type"] = "application/json"
wrk.headers["XCSRF-Token"] = "44f96d469525c8f5f62a664ec1b04514ec427ba407eb3a697e9c0efced84cd57.dcfde8c5-2041-4724-a033-50bde404ceab!S26JTfE6"
wrk.headers["Cookie"] = "session_id=dcfde8c5-2041-4724-a033-50bde404ceab; csrf_token=44f96d469525c8f5f62a664ec1b04514ec427ba407eb3a697e9c0efced84cd57.dcfde8c5-2041-4724-a033-50bde404ceab!S26JTfE6;"

local count = 1
local food_id = 72

request = function()
    count = count + 1

    local body = string.format('{"food_id": %d, "count": %d}', food_id, count)
    return wrk.format(nil, "/api/v1/order/food/update_count", nil, body)
end
`

### Скрипт для запуска

`wrk -t4 -c100 -d500s --script=add_food.lua --latency https://resto-go.online`

### Результат 
`
Running 8m test @ https://resto-go.online
  4 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   519.54ms  307.90ms   2.00s    76.01%
    Req/Sec    51.63     32.09   252.00     76.29%
  Latency Distribution
     50%  439.02ms
     75%  655.63ms
     90%  936.65ms
     99%    1.56s 
  99017 requests in 8.33m, 55.50MB read
  Socket errors: connect 0, read 0, write 0, timeout 364
Requests/sec:    198.02
Transfer/sec:    113.65KB
`

### Примечание
Если тесты не выдают ошибку на все запросы, то вы можете авторизоваться в [нашем сервисе](https://resto-go.online) и поменять session_id и csrf_token в Cookie, не забыв про заголовок `XCSRF-Token`.

### Вывод
Для нашего уровня нагрузки RPC == 2500 на чтение очень даже неплохо. Запись, конечно, похуже RPC == 200. Можно увеличить это числа путем увеличения мощности тачки, на которой запущен сервис, а также путем настроки СУБД. 

При такой нагрузке мы добились max-open-connections = 10. Если было бы больше открытых соединений, то мы бы обработали больше запросов в секунду, но поскольку у нас не такая нагрузка на сервис, достаточно того, что имеем. Запросы в БД оптимизированы, поэтому единственная видимая точка роста это - работа с connection pool и number of goroutines.

