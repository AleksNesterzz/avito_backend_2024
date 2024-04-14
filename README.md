# Что использовал
- БД: PostgreSQL (драйвер pgx)
- Роутер: chi
- Тестирование: testing
- Чтение конфига: cleanenv
- Развертиывание приложения: Docker
- Тестирование API: Postman
# Сборка и запуск проекта  
`docker build -t go-app . && docker-compose up --build go-app`
Приложение будет готово к использованию после вывода в консоль строки "Postgres is up - executing command"
# Для запуска тестов  
 `docker-compose run go-app go test ./apitest`
# Какие вопросы были
1. Каким образом реализовать отправку данных с 5-минутной задержкой, если у пользователя нет флага last_revision?
   Сделал локальный кэш, который через горутину обновляется каждые 5 минут из БД (в идеале использовать Redis)
# ПРИМЕРЫ ЗАПРОСОВ
Запросы в Postman:
1. Post-запрос на добавление сегмента
   ![image](https://github.com/AleksNesterzz/avito_backend_task_2023/assets/109950730/225e7c12-7db8-49db-879d-bb9e951940e5)   
   Добавим еще сегментов : AVITO_PERFORMANCE_VAS, AVITO_DISCOUNT_30 и AVITO_DISCOUNT_50
2. Delete-запрос на удаление сегмента
   ![image](https://github.com/AleksNesterzz/avito_backend_task_2023/assets/109950730/9822fcf7-e25c-48c6-bfdf-957468fcc762)
3. Patch-запрос на обновление сегментов пользователя
  ![image](https://github.com/AleksNesterzz/avito_backend_task_2023/assets/109950730/d979edbe-bf6f-466b-9fe8-de1a1859a3b2)
4. Get-запрос для получения сегментов пользователя
   ![image](https://github.com/AleksNesterzz/avito_backend_task_2023/assets/109950730/b52bd9db-2d47-45aa-9467-dd6483a3ac0b)
   При введении ID пользователя, у которого нет сегментов, то будет выведено segments : null
5. Get-запрос для получения истории ссылки на файл истории пользователя
   ![image](https://github.com/AleksNesterzz/avito_backend_task_2023/assets/109950730/22ea93d7-6990-4301-8b60-487d23d99776)   
   Вот что находится в этом файле:   
   ![image](https://github.com/AleksNesterzz/avito_backend_task_2023/assets/109950730/b384f0eb-fa6e-49f8-8aaf-8a96d73ab153)   
   Сегмент AVITO_END был добавлен и удален за кадром до создания сегмента AVITO_VOICE_MESSAGES.

   
   

   

