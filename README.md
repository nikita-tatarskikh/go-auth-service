# go-auth-service
https://www.notion.so/Test-task-Junior-BackDev-215fcddafff2425a8ca7e515e21527e7
1) Для работы необходим образ mongodb доступный на mongodb://localhost:27017
2) Приложене по дефолту запускается на  127.0.0.1:8080
3) Эндопинт 127.0.0.1:8080/sign-up: принимает post запрос, content-type: x-www-form-urlencoded, ожидаемый key:guid, по запросу возвращает пару токенов access и refresh.
4) 127.0.0.1:8080/refresh - принмает post запрос, content-type: application-json, возвращает пару новых токенов.

