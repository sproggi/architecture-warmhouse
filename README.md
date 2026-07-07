# Project_template

Это шаблон для решения проектной работы. Структура этого файла повторяет структуру заданий. Заполняйте его по мере работы над решением.

# Задание 1. Анализ и планирование

<aside>

"Тёплый дом" — это небольшая компания, которая организует удалённое управление отоплением в доме.
Архитектура приложения представляет из себя монолит на Go с СУБД Postgres.
Код синхронный.
Управление идёт от сервера к датчику.
Данные о температуре также получаются через запрос от сервера к датчику.
Установка и подключение новых датчиков выполняется только с выездом специалиста по подключению.

</aside>

### 1. Описание функциональности монолитного приложения

**Управление отоплением:**

- Пользователи могут удалённо управлять отоплением в своих домах.
- Система поддерживает получение данных о температуре датчиков, установленных в домах.

**Мониторинг температуры:**

- Пользователи могут, просматривать текущую температуру через веб-интерфейс.
- Система поддерживает получение данных о температуре датчиков, установленных в домах.

### 2. Анализ архитектуры монолитного приложения

- Язык программирования: Go
- СУБД: PostrgeSQL
- Все компоненты находятся в рамках одного приложения
- Взаимодействие между компонентами: Cинхронное

### 3. Определение доменов и границы контекстов

- Управление пользователями
- Управление устройствами
- Мониторинг устройств
- Платежи

### **4. Проблемы монолитного решения**

- Синхронное взаимодейтсвие компонентов.
- Любые фичи или устранения багов требуют остановки всего приложения.
- Ограниченная масштабируемость.

### 5. Визуализация контекста системы — диаграмма С4

```markdown
[Диаграмма контекста](https://github.com/sproggi/architecture-warmhouse/blob/warmhouse/own_c4/current_smart_home_uml.puml)
```

# Задание 2. Проектирование микросервисной архитектуры

В этом задании вам нужно предоставить только диаграммы в модели C4. Мы не просим вас отдельно описывать получившиеся микросервисы и то, как вы определили взаимодействия между компонентами To-Be системы. Если вы правильно подготовите диаграммы C4, они и так это покажут.

**Диаграмма контейнеров (Containers)**

```markdown
[Диаграмма контейнеров](https://github.com/sproggi/architecture-warmhouse/blob/warmhouse/own_c4/smart_home_containers.puml)
```

**Диаграмма компонентов (Components)**

```markdown
[Компонент API_GW](https://github.com/sproggi/architecture-warmhouse/blob/warmhouse/own_c4/api_gateway.puml)
[Компонент sensors](https://github.com/sproggi/architecture-warmhouse/blob/warmhouse/own_c4/sensors.puml)
[Компонент monitoring](https://github.com/sproggi/architecture-warmhouse/blob/warmhouse/own_c4/monitoring.puml)
[Компонент database](https://github.com/sproggi/architecture-warmhouse/blob/warmhouse/own_c4/database.puml)
[Компонент queue](https://github.com/sproggi/architecture-warmhouse/blob/warmhouse/own_c4/queue.puml)
[Компонент cache](https://github.com/sproggi/architecture-warmhouse/blob/warmhouse/own_c4/cache.puml)
```

**Диаграмма кода (Code)**

```markdown
[Диаграмма кода sensors](https://github.com/sproggi/architecture-warmhouse/blob/warmhouse/own_c4/sensors_code.puml)
[Диаграмма классов для sensors](https://github.com/sproggi/architecture-warmhouse/blob/warmhouse/own_c4/sensors_classes.puml)
```


# Задание 3. Разработка ER-диаграммы

```markdown
[Диаграмма ER](https://github.com/sproggi/architecture-warmhouse/blob/warmhouse/own_c4/smart_home_erd.puml)
```

# Задание 4. Создание и документирование API

### 1. Тип API

В проекте используется OpenAPI.
Более поддерживаемый и знакомый функционал.

### 2. Документация API

```markdown
[Open API Docs](https://github.com/sproggi/architecture-warmhouse/blob/warmhouse/openapi_doc.yaml)
```

# Задание 5. Работа с docker и docker-compose

Перейдите в apps.

Там находится приложение-монолит для работы с датчиками температуры. В README.md описано как запустить решение.

Вам нужно:

1) сделать простое приложение temperature-api на любом удобном для вас языке программирования, которое при запросе /temperature?location= будет отдавать рандомное значение температуры.

Locations - название комнаты, sensorId - идентификатор названия комнаты

```
	// If no location is provided, use a default based on sensor ID
	if location == "" {
		switch sensorID {
		case "1":
			location = "Living Room"
		case "2":
			location = "Bedroom"
		case "3":
			location = "Kitchen"
		default:
			location = "Unknown"
		}
	}

	// If no sensor ID is provided, generate one based on location
	if sensorID == "" {
		switch location {
		case "Living Room":
			sensorID = "1"
		case "Bedroom":
			sensorID = "2"
		case "Kitchen":
			sensorID = "3"
		default:
			sensorID = "0"
		}
	}
```

2) Приложение следует упаковать в Docker и добавить в docker-compose. Порт по умолчанию должен быть 8081

3) Кроме того для smart_home приложения требуется база данных - добавьте в docker-compose файл настройки для запуска postgres с указанием скрипта инициализации ./smart_home/init.sql

Для проверки можно использовать Postman коллекцию smarthome-api.postman_collection.json и вызвать:

- Create Sensor
- Get All Sensors

Должно при каждом вызове отображаться разное значение температуры

Ревьюер будет проверять точно так же.


