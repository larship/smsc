Пакет предназначен для отправки сообщений через сервис SMS Центр (https://smsc.ru)

## Подключение

```
go get github.com/larship/smsc
```

## Инициализация

```
smscClient := smsc.New("login", "password")
```

## Отправка СМС

```
resp, err := smscClient.SendSms("phone", "text")
```
