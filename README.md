# Тестовое окружение для запуска нагрузочных тестов

## Содержимое

`service` - тестовый сервис

`k8s` - манифесты k8s для развертывания сервиса в кластере

`config` - конфигурация нагрузочных тестов для Яндекс.Танк

`config/res` - результаты (сырые данные) тестов


## Запуск сервиса

1. Из папки сервиса при наличии интерпретатора языка go:

    ```bash
    $ go run ./cmd/service/main.go
    ```

2. Из docker'а:
    ```bash
    $ docker build . -t test-service
    $ docker run -p 8080:8080 test-service
    ```

Далее к сервису можно обратиться по адресу `http://localhost:8080/_info` 


## MicroK8s

Установка зависит от используемой ОС и основная документация: https://microk8s.io/docs/install-alternatives

Нужны буду дополнения: `dashboard, ingress, prometheus, registry`
```
$ microk8s enable dns dashboard ingress prometheus registry
```

### Обновление кода в registry 

>Здесь и далее `<ip-адрес>` - это ip-адрес microk8s

```
service/$ docker build . -t <ip-адрес>:32000/test-service:v1
service/$ docker push <ip-адрес>:32000/test-service
```

### Создание namespace и обновление манифестов

```
$ microk8s kubectl create namespace services
```

> `multipass transfer` выполняется только для MacOS X (надо уточнить для вашей системы)

Деплоймент:
```
k8s/$ multipass transfer ./test-service-deployment.yaml microk8s-vm:
k8s/$ microk8s kubectl apply -f test-service-deployment.yaml -n services
```

Сервис:
```
k8s/$ multipass transfer ./test-service-service.yaml microk8s-vm:
k8s/$ microk8s kubectl apply -f test-service-service.yaml -n services
```

Ингресс:
```
k8s/$ multipass transfer ./test-service-ingress.yaml microk8s-vm:
k8s/$ microk8s kubectl apply -f test-service-ingress.yaml -n services
```

### Дашборд

```
$ microk8s dashboard-proxy
```

и далее по адресу https://<ip-адрес>:10443


## Запуск тестов yandex.tank

```
$ docker run -v $(pwd)/config:/var/loadtest --rm -it direvius/yandex-tank -c /var/loadtest/test.yaml
```
