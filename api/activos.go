// inicio de programa
package main

//importar las herramientas que el backend necesita

import (
	//paquete de go para trabajar con bases de datos
	"database/sql"
	//sirve para convertir JSON -> structs de GO
	"encoding/json"
	//imprime mensajes de consola en el servidor
	"log"
	//corazon del servidor web / crea rutas, maneja peticiones, envia respuestas, levanta servidor
	"net/http"
	//este es el driver de MySQL para GO
	_ "github.com/go-sql-driver/mysql"
)

// el struct es un molde o plantilla que debe contener los datos que se van a recibir o enviar
// sirve para recibir JSON
type Activo struct {
	NombreResponsable string `json:"nombre_responsable"`
	Departamento      string `json:"departamento"`
	Tipo              string `json:"tipo"`
	Marca             string `json:"marca"`
	Modelo            string `json:"modelo"`
	Serie             string `json:"serie"`
	Fecha             string `json:"fecha"`
	Observacion       string `json:"observacion"`
}

// Middleware CORS
func enableCors(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func main() {

	//CONEXIÓN A BD sql.Open Maneja dos valores no uno por eso es db y err
	//db = puntero de conexión para ejecutar insetar y hacer SELECT, UPDATE, DELETE
	db, err := sql.Open("mysql", "root:7896@tcp(127.0.0.1:3306)/inventario")
	//Si hubo error al abrir conexión muestra error
	if err != nil {
		log.Fatal(err)
	}
	//cierra base de datos
	defer db.Close()

	//aqui llega la peticion HTTP a la ruta /api/activos
	// w http.ResponseWriter OBJETO PARA RESPONDER AL FRONTEND
	// r parametro que presenta la petición HTTP que envió el frontend
	http.HandleFunc("/api/activos", func(w http.ResponseWriter, r *http.Request) {
		enableCors(w)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		//http.MethodPost es constante de Go, Go tiene constantes para medotos HTTP
		if r.Method != http.MethodPost {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}

		//al crear un struct llamado Activo, se puede crear una variable activo de tipo Activo
		var activo Activo

		//convertir el JSON del frontend a un struct usable en go
		// json.NewDecoder(r.Body) -> crea un lector de JSON
		// Decode(&activo) -> convierte el JSON y llena los campos del struct que creamos Activo
		// &activo se crea el puntero hacia el struct Activo para poder modificarlo
		err := json.NewDecoder(r.Body).Decode(&activo)
		if err != nil {
			http.Error(w, "JSON INVALIDO", http.StatusBadRequest)
			return
		}

		//validar si todos los campos bienen bien

		if activo.NombreResponsable == "" ||
			activo.Departamento == "" ||
			activo.Tipo == "" ||
			activo.Marca == "" ||
			activo.Modelo == "" ||
			activo.Serie == "" ||
			activo.Fecha == "" ||
			activo.Observacion == "" {

			http.Error(w, "Todos los campos son obligatorios", http.StatusBadRequest)
			return
		}

		//ejecutar un comando que no devuelve filas INSERT, UPDATE, DELETE, CREATE
		//Devuelve sql.Result y error (si algo salió mal)
		// _, es un identificador en blanco, dn.Exec() nesecita dos parametros pero con _ ignoramos el primer parametro
		_, err = db.Exec(`
			INSERT INTO inventario
			(responsable, departamento, tipo, marca, modelo, serie, fecha_registro, observacion)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		`,
			activo.NombreResponsable,
			activo.Departamento,
			activo.Tipo,
			activo.Marca,
			activo.Modelo,
			activo.Serie,
			activo.Fecha,
			activo.Observacion,
		)

		if err != nil {
			http.Error(w, "Error al insertar en Mysql", http.StatusInternalServerError)
			return
		}

		//w.Header devuelve un mapa de headers HTTP (tipo de contenido, servidor, fecha)
		//.Set("Content-Type", "application/json") le dice al navegador que la respuesta que se enviará es JSON
		w.Header().Set("Content-Type", "application/json")
		//json.NewEncoder -> crea un Encoder JSON que va a escribir al objeto w que es l http.ResponseWriter
		//Encode() -> convierte lo que pases a JSON, lo escribe en el objeto w, agrega salto de linea y devuelve error si algo falla
		//map{string}string{} -> se crea un mapa en Go la llave es un string y es valor es un string, para crear un objeto JSON sin usar struct
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "ok",
			"message": "Activo registrado correctamente",
		})
	})

	//imprimir en consola que el servidor esta funcionando con log.
	log.Println("Servidor corriendo en http://localhost:8080")
	//aqui es donde inicia el servidor HTTP y se queda escuchando peticiones
	http.ListenAndServe(":8080", nil)

}
