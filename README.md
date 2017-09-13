# Хакатон по курсу golang

Пишем биржу :)

Биржа состоит из 3-х компонентов:
- Клиент
- Брокер
- Биржа

## Биржа

Биржа - центральная точка всей системы. Она сводит между собой различных продавцов и покупателей из разных брокеров, результатом сделок которых является изменение цены.

В нашем хакатоне будут использоваться исторические данные реальных торгов по фьючерсным контрактам на бирже РТС - https://cloud.mail.ru/public/7UCu/9SuVemxkb

Если на биржу ставится заявка на покупку или продажу, то она ставится в очередь и когда цена доходит до неё и хватает объёма - заявка исполняется, брокеру уходит соответствующее уведомление. Если не хватает объёма, то заявка исполняется частичсно, брокеру так же уходит уведомление.
Если несколько участников поставилос заявку на одинаковый уровеньт цены, то исполняются в порядке добавления.

В связи с тем что наши данные исторические - предполагем, что мы не можем выступить инициатором сделки (т.е. сдвинуть цену, купив по рынку), можем приобрести только если другая сторона выступила инциатором.

Помимо исполнения сделок брокер транслирует цену инструментов всем подключенным брокерам. Список транслируемых инструментов берётся из конфига, сами цены - из файла.
В связи с тем что цены синтетические - мы не смотрим на цены, а просто начинаем таранслировать ту цену что есть. Инфомрация о изменении цены отправляется каждую секунду. Если в секунду ( под конфигом) произошло больше чем 1 сделка - они аггрегируются. Брокеру отправляется OHLCV ( open, high, low, close, volume ), где:
- open - цена открытия интервала (первая сделка)
- high - максимальная цена в интервале
- low - минимальная цена в интервале
- close - цена закрытия интервала ( последняя сделка )
- volume - количество проторгованных контрактов

Формат обмена данными с брокером - protobuf через GRPC

```protobuf
syntax = "proto3";

message OHLCV {
  int64 ID = 1;
  int32 Time = 2;
  int32 Interval = 3;
  float Open = 4;
  float High = 5;
  float Low = 6;
  float Close = 7;
  int32 Time = 8;
  string Ticker = 9;
}

message Deal {
    int32 BrokerID = 1;
    int32 ClientID = 2;
    string Ticker = 3;
    int32 Amount = 4;
    bool Partial = 5;
    int32 Time = 6;
    float Price = 7;
}

```

## Брокер

Брокер - это организация, которая предсотавляет своим клиентам доступ на биржу.
У неё есть список клиентов, которые могут взаимодействовать посредством неё с биржей, так же она хранит количество их позиицй и историю сделок.

Брокер аггрегирует внутри себя информацию от биржи по ценовым данным, позволяя клиенту посмотреть историю. По-умолчанию, хранится история за последнеи 5 минут (300 секунд).

Брокер предоставляет клиентам JSON-апи (REST или JSON-RPC), через который им доступныы следующие возможности:
- посмотреть свои позиции и баланс
- посмотрить историю своих сделок
- отправить на биржу заявку на покупку или продажу тикера
- отменить ранее отправленную заявку
- посмотреть последнюю истории торгов

```sql
CREATE TABLE `clients` (
    `id` int NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `login` varchar(300) NOT NULL,
    `password` varchar(300) NOT NULL,
    `balance` int NOT NULL
);

INSERT INTO `clients` (`id`, `login`,  `password`, `balance`) 
    VALUES (1, 'Vasily', '123456', 200000),
    VALUES (2, 'Ivan', 'qwerty', 200000),
    VALUES (3, 'Olga', '1qaz2wsx', 200000);

CREATE TABLE `positions` (
    `id` int NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `user_id` int NOT NULL,
    `ticker` varchar(300) NOT NULL,
    `amount` int NOT NULL,
    KEY user_id(user_id)
);

INSERT INTO `clients` (`user_id`, `ticker`, `amount`) 
    VALUES (1, 'SiM7', '123456', 200000),
    VALUES (1, 'RIM7', '123456', 200000),
    VALUES (2, 'RIM7', 'qwerty', 200000);
    
CREATE TABLE `orders_history` (
    `id` int NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `time` int NOT NULL,
    `user_id` int,
    `ticker` varchar(300) NOT NULL,
    `amount` int NOT NULL,
    KEY user_id(user_id)
);
    
```

## Клиент

Клиент - это любой пользователь АПИ брокера. Это может быть биржевой терминал, веб-сайт, торговый робот.

Мы делаем простой биржевой терминал, который нам позволяет покупать-продавать и смотреть свою историю.

https://s.mail.ru/Evbb/F2jm79n8m


## Организация работы

Обсуждаем проект, делимся на команды, пишем код. 3 компонента, 3 команды. Работам короткими этапами по 45 минут, стараясь в конце этапа иметь логически завершенный кусок.

Каждая команда мерджит свой компонент в мастер, если считает что всё будет рабоать.

Проект собираем через gb, чтобы не иметь проблем с GOPATH.

По возможности - пишем тесты.