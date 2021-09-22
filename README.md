# go-auth-service
https://www.notion.so/Test-task-Junior-BackDev-215fcddafff2425a8ca7e515e21527e7
1) Для работы необходим образ mongodb доступный на mongodb://localhost:27017
2) Приложене по дефолту запускается на  127.0.0.1:8080
3) Эндопинт 127.0.0.1:8080/sign-up: принимает post запрос, content-type: x-www-form-urlencoded, ожидаемый key:guid, по запросу возвращает пару токенов access и refresh.
![image](https://user-images.githubusercontent.com/34633194/134423165-175ccbbd-42bb-48ec-8f9a-dc3f3c7ecd9c.png)
5) 127.0.0.1:8080/refresh - принмает post запрос, content-type: application-json, возвращает пару новых токенов.
![image](https://user-images.githubusercontent.com/34633194/134423345-ff9ea302-cad7-46df-9d79-d39282eb727d.png)

- Примечания
1) Отловлены всевозможные ошибки.
2) secretKey оставлен в коде, т.к это не prod app.
- TO-DO:
1) Для удобности чтения и работы с кодом, неплохо было бы сделать рефакторинг - убрать повторяющийся код в отдельные функции.
2) Задать структуру проекта и разгрузить main.go

