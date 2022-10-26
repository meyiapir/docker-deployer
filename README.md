# Deployer

***
**`VERSION: DEV-0.6`**

## Сборка и запуск приложения:

1) Установить Go: 1.18.2
2) В терминале перейти в папку `deployer`
3) Запустить команду(указать аргументы при необходимости):
    ```shell
    go build main.go
    ```
4) В папке deployer появиться исполняемый файл - Готово!

## Конфигурация:
В папке `deployer\сonfigs` файл `config.json`
```json
{
   "registry": {
      "address": "url",
      "port": 5000
   },
   "delete_after_push": {
      "enabled": true
   }
}
```
Параметр `delete_after_push` - удалять ли файлы после отправки в репозиторий
Параметр `registry` - адрес репозитория
