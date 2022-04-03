# Mutants

## Project layout
Ver https://github.com/golang-standards/project-layout

## Configuración
* `maxStreak: 4` Cantidad de repeticiones de una misma letra para que cuente como una secuencia mutante
* `maxCount: 1` Cantidad máxima de secuencias mutantes para que el ADN sea considerado humano

## Solucion planteada (Parte 1)
Se usa un slice 2D auxiliar para calcular la cantidad de rachas que se completan en determinado punto. Esta estructura se recorre de arriba a abajo y de izquierda a derecha.
Para cada punto, se toman en cuenta los puntos que conozco (izquierda, izquierda-arriba, arriba, derecha-arriba) y se calcula en base a las rachas de cada uno de esos puntos en todas las direcciones. Por ejemplo: si estamos procesando el punto (i,j) y quiero saber si obtendré una racha viniendo desde arriba, se toma el "upperStreak" de (i-1,j) y se comparan. Esto genera que para cada punto se tenga las rachas que tiene en ese momento para todas los puntos actualmente conocidos.
En el momento que el algoritmo encuentre más de 1 patrón, devuelve true, en el caso contrario, termina de recorrer toda la entrada y retorna false.
El orden del algoritmo es NxN

## Solucion planteada (Parte 2 y 3)
Se usa Redis como base de datos y un bloqueo optimista para cumplir con el requerimiento de tráfico teorico que recibirían los endpoints. Se utiliza una clave para guardar el hash md5 en base 64 del ADN y su tipo. Un "status" se utiliza como un hash para guardar la cantidad de humanos y mutantes procesados, por cada post recibido se suma 1 si no existe en la base, en caso contrario no se hace nada y se evita recalcular el algoritmo retornando el tipo de ADN que hay en la base.

## Make commands
* `make build`   # Genera challenge y lo deja en el directorio bin
* `make run`     # Ejecuta la aplicacion challenge
* `make test`    # Ejecuta los test y muestra la cobertura por paquete

## Requerimientos
* Golang
* Redis (configuracion del servidor en configs/challenge.yaml)
* Make

## Cloud endpoints
* POST /mutant Mutante

```
curl -v --location --request POST 'http://34.176.170.125:9091/mutant' \
--header 'Content-Type: application/json' \
--data-raw '{
    "dna": ["ATGCGA", "CAGTGC", "TTATGT", "AGAAGG", "CCCCTA", "TCACTG"]
}'
```

* POST /mutant Humano
```
curl -v --location --request POST 'http://34.176.170.125:9091/mutant' \
--header 'Content-Type: application/json' \
--data-raw '{
    "dna": ["ATGCGA", "CAGTGC", "TTATTT", "AGACGG", "GCGTCA", "TCACTG"]
}'
```

* GET /stats
```
curl --location --request GET 'http://34.176.170.125:9091/stats'
```
