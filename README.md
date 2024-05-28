1. Описание проекта:
Веб-сервер + оболочка (фронтенд), которые реализуют функциональность планировщика задач.

2. Выполнены все (кроме месяцев) задачи со звездочкой:
- Реализована возможность определять извне порт при запуске сервера
- Реализована возможность определять путь к файлу базы данных через переменную окружения
- Правила повторения задач:
    - Недели
    - Месяцы (не сделано)
- Рабочая строка поиска
- Аутентификация
- Создание докер образа

3. Инструкция по запуску кода локально: 
    - Примеры .env: 
        TODO_PORT=7540
        TODO_DB_FILE=../scheduler.db
        TODO_PASSWORD=123

    - Адрес, который следует указывать в браузере:
        http://localhost:7540/

4. Параметры в tests/settings.go следует использовать:
var Port = 7540
var DBFile = "../scheduler.db"
var FullNextDate = false 
var Search = true
var Token = *инструкция по созданию токена*
    - Для создания токена нужно явно указать пароль с помощью команды 
    "export TODO_PASSWORD=123" в консоли.
    - Авторизоваться на сайте http://localhost:7540/login.html
    - Достать JWT-токен из логов консоли и вставить в settings.go (тест)

5. Docker образ: 
Команды bash для сборки и запуска DOCKER образа:
    docker build --tag todofinal:v2 .
    docker run -d -p 7540:7540 todofinal:v2

6. Файлы:
    - В директории `tests` находятся тесты для проверки API, которое должно быть реализовано в веб-сервере.
    - Директория `web` содержит файлы фронтенда.
    - Директория `internal` содержит поддиректории с логикой веб-сервера:
        - `database` - логика работы с базой данных
        - `date` - реализация правила повторения задач
        - `handlers` - логика обработки HTTP-запросов
        - `middleware` - механизм middleware
        - `models`- модели для работы с базой данных