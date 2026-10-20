# Imagen Go Lang v1.27 Para Backend
FROM golang:1.27

# Ruta de Trabajo
WORKDIR /usr/share/providencia

# Copia de Dependencias "Go Lang"
COPY ../go.mod ../go.sum ./

# Comando "Descargar dependecias"
RUN go mod download

# Copia Directorios y Ficheros de la aplicacion
COPY ../services ./services
COPY ../config ./config
COPY ../assets ./assets
COPY ../utils ./utils
COPY ../main.go ./

# Crear directorio DB y PRIVATE
RUN mkdir -p db private

# Crear binario "Go Lang"
RUN go build -o providencia main.go

# Se ejecuta el binario
CMD ["./providencia"]
