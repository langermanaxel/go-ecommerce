# Go Ecommerce

Este es un proyecto de ecommerce desarrollado en Go utilizando el framework Gin.

## Requisitos

- Go 1.24.1 o superior

## Instalación

1. Clona el repositorio:

   ```sh
   git clone https://github.com/langermanaxel/go-ecommerce.git
   cd go-ecommerce
   ```

2. Instala las dependencias:

   ```sh
   go mod tidy
   ```

## Ejecución

Para ejecutar el servidor, utiliza el siguiente comando:

```sh
go run main.go
```

## Dependencias

Este proyecto utiliza las siguientes dependencias:

- [gin-gonic/gin](https://github.com/gin-gonic/gin) v1.10.0
- [go-playground/validator](https://github.com/go-playground/validator) v10.20.0
- [golang-jwt/jwt](https://github.com/golang-jwt/jwt) v3.2.2
- [mongo-driver](https://go.mongodb.org/mongo-driver) v1.17.3
- [x/crypto](https://golang.org/x/crypto) v0.26.0

## Estructura del Proyecto

```sh
go-ecommerce/
├── main.go
├── go.mod
├── go.sum
└── ...
```

## Contribuir

1. Haz un fork del proyecto
2. Crea una nueva rama (`git checkout -b feature/nueva-funcionalidad`)
3. Realiza tus cambios y haz commit (`git commit -am 'Agrega nueva funcionalidad'`)
4. Sube tus cambios a tu fork (`git push origin feature/nueva-funcionalidad`)
5. Abre un Pull Request
