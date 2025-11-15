# Implementación del protocolo HTTP/1.1 con una conexión TCP.
*El proyecto no implementa inovaciones sino es más un proyecto para aprender como funciona la conexión por TCP y el protocolo HTTP*
1. Incluye la lectura por bytes a traves de un io.Reader que incluye conexiónes (como tcp), pero tamién implementa cualquier otro que de stream de bytes.
2. Implementa el parser que lee y parsea la request line, headers y body.
3. Implementa http responses de forma normal, y por chunks (para archivos pesados)
4. Implementa un server mock con las funcionalidades de un servidor

# Para ejecutar
1. Instala go de la forma que más te guste [forma mas sencilla](https://webinstall.dev/golang/)
2. `go run ./cmd/httpserver` Esto abrira el server en el puerto 42069 😎
3. Actualmente nomas están las páginas /yourproblem, /myproblem, /video y la default que sería cualquier otra.
4. si quieres ver que funcionen los videos, puedes descargar y poner tu video en la carpeta de assets.
5. Además, para probar las responses por chunked encoding puedes ir a /httpbin que hace un proxy a https://httpbin.org/ que es una página web para hacer pruebas de conexiones http. Aquí yo use el comando
`curl --raw -v http://localhost:42069/httpbin/stream/100` Para probar streams de 100 respuestas json

# Referencia
El proyecto fue realizada de forma guíada en [bootdev](https://www.boot.dev/lessons/b0cebf37-7151-48db-ad8a-0f9399f94c58).
Donde se describe las caracteríticas del protocolo, el uso de stream de bytes, las diferencias de tcp, udp y http y los puntos a codear de forma escrita (sin convertirse en un curso de copy paste).
