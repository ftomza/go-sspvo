![report go](https://goreportcard.com/badge/github.com/ftomza/go-sspvo)
![Build and Test](https://github.com/ftomza/go-sspvo/workflows/Build%20and%20Test/badge.svg)
![codecov](https://codecov.io/gh/ftomza/go-sspvo/branch/master/graph/badge.svg)

# go-sspvo
Module for working with **Суперсервис "Поступление в ВУЗ онлайн"**
## Назначение
Модуль предназначен для подготовки сообщений с дальнейшей отправкой в сервис. Подготовка включает в себя добавление основных полей, таких как *ОГРН* *КПП*, формирования подписанного токена для отправки данных, валидация приходящего токена. Подпись и валидация проводится согласно требованиям сервиса, а именно использования инфраструктуры открытого ключа ГОСТ. Версия Golang 1.14

## Установка
```bash
go get "github.com/ftomza/go-sspvo"
```

## Использование
Пример получения справочника уровней бюджета:
```go
restyClient := resty.New()  
restyClient.SetHostURL("http://localhost:7777")  
restyClient.SetDebug(false)  
sspvoClient, err := client.NewRestyClient(restyClient,  
  client.SetAPIBase("/api"),  
  client.SetOGRN("test"),  
  client.SetKPP("test"),  
)  
if err != nil {  
   log.Fatal(err)  
}  
  
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)  
defer cancel()  
  
data, err := sspvoClient.Send(ctx, message.NewCLSMessage(message.CLSLevelBudget)).Data()  
if err != nil {  
   log.Fatal(err, string(data))  
}  
  
log.Print(string(data))

// Output:
//  
// <?xml version="1.0" encoding="UTF-8"?>  
// <LevelBudget>  
// <Budget><ID>1</ID><Code></Code><Name>Федеральный</Name><Actual>true</Actual></Budget>  
// <Budget><ID>2</ID><Code></Code><Name>Региональный</Name><Actual>true</Actual></Budget>  
// <Budget><ID>3</ID><Code></Code><Name>Муниципальный</Name><Actual>true</Actual></Budget>  
// </LevelBudget>
```

Пример получения количества токенов ожидающий получения и подтверждения.
```go
restyClient := resty.New()  
restyClient.SetHostURL("http://localhost:7777")  
restyClient.SetDebug(false)  
sspvoClient, err := client.NewRestyClient(restyClient,  
  client.SetAPIBase("/api"),  
  client.SetOGRN("test"),  
  client.SetKPP("test"),  
)  
if err != nil {  
   log.Fatal(err)  
}  
  
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)  
defer cancel()  
  
data, err := sspvoClient.Send(ctx, message.NewInfoAllMessage()).Data()  
if err != nil {  
   log.Fatal(err, string(data))  
}  
  
log.Print(string(data))

// Output:
//  
// {"Messages": 2}
```
Пример отправки данных, получения информации и ее подтверждения.
```go
var (  
   key = `  
-----BEGIN PRIVATE KEY-----  
...
-----END PRIVATE KEY-----  
`  
  cert = `  
-----BEGIN CERTIFICATE-----  
...
-----END CERTIFICATE-----  
`  
)  
 
restyClient := resty.New()  
restyClient.SetHostURL("http://localhost:7777")  
sspvoClient, err := client.NewRestyClient(restyClient,  
	client.SetAPIBase("/api"),  
	client.SetOGRN("test"),  
	client.SetKPP("test"),  
)  
if err != nil {  
	log.Fatal(err)  
}  

gostCrypto, _ := crypto.NewGostCrypto(crypto.SetCert(cert), crypto.SetKey(key))  

ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)  

data, err := sspvoClient.Send(ctx, message.NewActionMessage(  
	gostCrypto,  
	message.ActionAdd,  
	message.DatatypeSubdivisionOrg,  
	[]byte(`  
<?xml version="1.0" encoding="utf-8"?>  
<PackageData>  
	<SubdivisionOrg> 
		<UID>TEST69</UID> 
		<Name>Подвал</Name> 
	</SubdivisionOrg>
</PackageData>  
	`))).Data()  
cancel()  
if err != nil {  
	log.Fatal(err, string(data))  
}  

log.Println("Add SubdivisionOrg:", string(data))  

newMessage := struct {  
	IDJWT string `json:"IDJWT"`  
}{}  

err = json.Unmarshal(data, &newMessage)  
if err != nil {  
	log.Fatal(err, string(data))  
}  

idJWT, _ := strconv.Atoi(newMessage.IDJWT)  
ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)  
data, err = sspvoClient.Send(ctx, message.NewInfoMessage(gostCrypto, idJWT)).Data()  
cancel()  
if err != nil {  
	log.Fatal(err, string(data))  
}  

log.Println("Info:", string(data))  

ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)  
data, err = sspvoClient.Send(ctx, message.NewConfirmMessage(gostCrypto, idJWT)).Data()  
cancel()  
if err != nil {  
	log.Fatal(err, string(data))  
}  

log.Println("Confirm:", string(data))   

// Output:
//
// 2020/10/01 21:09:06 Add SubdivisionOrg: {"IDJWT":"1405161"}
// 2020/10/01 21:09:07 Info: {"ResponseToken":".."}
// 2020/10/01 21:09:08 Confirm: {"IDJWT":"1405161","Result":"true"}
```

## Содержание
#### Крипто модуль `GostCrypto` из пакет `crypto`
Данный крипто модуль поддерживает инфраструктуру открытых ключей ГОСТ Р 34.10-2001 и хэш функцию ГОСТ Р 34.11-2012 Стрибог (Streebog).
Конструктор:
- `NewGostCrypto(opts ...Option) (sspvo.Crypto, error)`.

Поддерживаемы опции:
- `SetCert(cert string) Option` - задать сертификат с открытым ключом в формате *PEM*
- `SetKey(key string) Option` - необязательный, задать закрытый ключ в формате *PEM*

Если будет передан только сертификат, то крипто модуль будет поддерживать только проверку подписи.
#### Виды Сообщений, пакет `message`
Все сообщен возвращают данные в виде массива байт.
##### Простые сообщения, `Message`
##### `CLSMessage`
Используется для получения справочников из сервиса, например, уровни бюджета, олимпиады, достижения и так далее. Результат выполнения будет XML ответ содержащий перечисление запрашиваемого справочника:
```XML
<?xml version="1.0" encoding="UTF-8"?>  
<LevelBudget>  
	<Budget><ID>1</ID><Code></Code><Name>Федеральный</Name><Actual>true</Actual></Budget>  
	<Budget><ID>2</ID><Code></Code><Name>Региональный</Name><Actual>true</Actual></Budget>  
	<Budget><ID>3</ID><Code></Code><Name>Муниципальный</Name><Actual>true</Actual></Budget>  
</LevelBudget>
```
Конструктор:
- `NewCLSMessage(cls CLS) *CLSMessage`.

Параметры:
- `cls CLS` - принимает значение допустимого *класса* справочника, *классы* перечислены в виде констант начинающихся на `CLS*`

##### `InfoAllMessage`
Простое сообщение для получения количества токенов ожидающих обработку(получения информации по ним и подтверждения полученной информации). Результатом выполнения будет JSON ответ содержащий количество токенов:
```json
{"Messages": 2}
``` 
Конструктор:
- `NewInfoAllMessage() *InfoAllMessage`.

##### Крипто сообщения, `SignMessage`
Сообщение с поддержкой подписи отправляемых сообщений и проверкой подписи входящего сообщения.
##### `ActionMessage`
Используется для отправки данных в сервис, например, создание конкурсной группы, направлений и так далее. Результатом выполнения будет JSON ответ содержащий номер токена:
```json
{"IDJWT":"1405161"}
```
Конструктор:
- `NewActionMessage(crypto sspvo.Crypto, action Action, datatype Datatype, data []byte) *ActionMessage`.

Параметры:
- `crypto sspvo.Crypto` - принимает значение крипто модуля.
- `action Action` - принимает значение допустимого *действия* сервиса, *действия* перечислены в виде констант начинающихся на `Action*`.
- `datatype Datatype` - принимает значение допустимого *типа данных* сервиса, *типы данных* перечислены в виде констант начинающихся на `Datatype*`
- `data []byte` - принимает массив байтов для передачи в сообщении полезной нагрузки *Payload*, может содержать пустое или `nil` значение.

##### `ConfirmMessage`
Используется для отправки подтверждения полученной информации по номеру токена. Результатом выполнения будет JSON ответ содержащий результат подтверждения:
```json
{"IDJWT":"1405161","Result":"true"}
```
Конструктор:
- `NewConfirmMessage(crypto sspvo.Crypto, idJWT int) *ConfirmMessage`.

Параметры:
- `crypto sspvo.Crypto` - принимает значение крипто модуля.
- `idJWT int` - принимает значение идентификатора токена.

##### `InfoMessage`
Используется для получения информации по номеру токена. Результатом выполнения будет JSON ответ содержащий токен:
```json
{"ResponseToken": "[headers in base64].[payload in base64].[sign in base64]"}
```
Конструктор:
- `NewInfoMessage(crypto sspvo.Crypto, idJWT int) *InfoMessage`.

Параметры:
- `crypto sspvo.Crypto` - принимает значение крипто модуля.
- `idJWT int` - принимает значение идентификатора токена.

Если в качестве параметра идентификатора токена передать `0`, то будет возвращен первый ожидающий обработки токен.

#### Клиент отправки сообщений, пакет `client`
В данном решении для отправки данных сообщения используется REST клиент  [resty.v2](https://github.com/go-resty/resty).
Конструктор:
- `NewRestyClient(rest *resty.Client, opts ...Option) (sspvo.Client, error)`

Параметры:
- `rest *resty.Client` - REST клиент

Опции:
- `SetAPIBase(apiBase string) Option` - задать базовый путь сервиса
- `SetOGRN(ogrn string) Option` - задать ОГРН для аутентификации на сервисе
- `SetKPP(kpp string) Option` - задать КПП для аутентификации на сервисе
