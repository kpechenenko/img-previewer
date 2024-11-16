# Первьювер изображений

Веб-сервис, позволяющий создавать превью (уменьшенные копии) изображений из сети в формате jpeg.
Финальный проект по курсу Otus "GO разработчик".

## Запуск сервиса

Конфигурация выполняется через переменные окружения.
По умолчанию приложение запускается на localhost на 8080 порту.

Без docker

```bash
make run
```

С docker

```bash
make run-img
```

## Использование сервиса

После запуска сервиса перейти
по [http://127.0.0.1:8080/preview/100/300/raw.githubusercontent.com/OtusGolang/final_project/master/examples/image-previewer/_gopher_original_1024x504.jpg](http://127.0.0.1:8080/preview/100/300/raw.githubusercontent.com/OtusGolang/final_project/master/examples/image-previewer/_gopher_original_1024x504.jpg)

* http://127.0.0.1:8080/ - адрес приложения,
* `preview/` - эндпоинт для создания превью,
* `100` и `300` - размер превью,
* `raw.githubusercontent.com/OtusGolang/final_project/master/examples/image-previewer/_gopher_original_1024x504.jpg` -
  ссылка на исходное изображение.

## Запуск тестов

```bash
make test
```

## Запуск линтера

```bash
make lint
```