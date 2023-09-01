#utiliza la imagen de golang
FROM golang:latest

#define el directorio de trabajo
WORKDIR /go/src/app

#inicializamos el modulo
RUN go mod init

# Instala las dependencias
RUN go install github.com/codegangsta/gin@latest

# Copia el contenido del directorio actual en el directorio de trabajo
COPY ./app .

# Expone el puerto 8080 y 6061
EXPOSE 6061 8080 