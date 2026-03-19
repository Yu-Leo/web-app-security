# Web App Security

Система обнаружения аномалий в трафике веб-приложений.


## Деплой

```bash
helm install service-echo ./infrastructure/service-echo

helm install envoy ./infrastructure/envoy
```

## Makefile команды

```bash
make up            # поднять все компоненты в docker + применить миграции
make down          # удалить все контейнеры и volume БД
make stop          # остановить контейнеры без удаления БД
make add-mock-data # создать мок-данные через backend API

make db-up         # поднять только postgres в docker и накатить миграции
make db-stop       # остановить postgres
make db-down       # удалить контейнер postgres и volume БД

make logs          # смотреть логи всех контейнеров
make ps            # список контейнеров
make build         # собрать образы
```
