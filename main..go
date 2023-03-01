package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type Point struct {
	X float64
	Y float64
}

// Estructuras para las distancias y mensajes recibidos por el cliente (nave)
type satelliteRec struct {
	Name     string   `json:"name"`
	Distance float64  `json:"distance"`
	Message  []string `json:"message"`
}

type satelliteRecsplit struct {
	Distance float64  `json:"distance"`
	Message  []string `json:"message"`
}
type position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}
type respuestaCliente struct {
	Position position `json:"position"`
	Message  string   `json:"message"`
}

type salellitesRecAll struct {
	Satellites []satelliteRec `json:"satellites"`
}

type satellitesRecBody []satelliteRec

//Estructuras para las posiciones conocidas de cada satelite

type satellitePos struct {
	Name       string `json:"name"`
	position_x float64
	position_y float64
}
type salellitesPosAll []satellitePos

var satellitesPosition = salellitesPosAll{
	{
		Name:       "Kenobi",
		position_x: -500,
		position_y: -200,
	},
	{
		Name:       "Skywalker",
		position_x: 100,
		position_y: -100,
	},
	{
		Name:       "Sato",
		position_x: 500,
		position_y: 100,
	},
}

var dist1, dist2, dist3 float64
var msgS1, msgS2, msgS3 []string
var respuestaCli = respuestaCliente{}

//funcion para determinar las cordenadas de la nave

/*func GetLocation(distances ...float64) (x, y float64) {

	var x1, x2, y2, x3, y1, y3 float64
	//Se extraen las posiciones conocidas de los satelites
	x1 = satellitesPosition[0].position_x
	y1 = satellitesPosition[0].position_y
	x2 = satellitesPosition[1].position_x
	y2 = satellitesPosition[1].position_y
	x3 = satellitesPosition[2].position_x
	y3 = satellitesPosition[2].position_y

	distancias := []float64{}

	//Se leen las distancias de los satelites a las naves
	for _, t := range distances {
		//fmt.Println(t)
		distancias = append(distancias, t)
	}
	d1 := distancias[0]
	d2 := distancias[1]
	d3 := distancias[2]

	//Formula de la trilateracion en 2 Dimensiones
	x = (math.Pow(d1, 2) - math.Pow(d2, 2) + math.Pow(x2, 2) - math.Pow(x1, 2)) / (2 * (x2 - x1))
	y = (math.Pow(d1, 2)-math.Pow(d3, 2)+math.Pow(x3, 2)-math.Pow(x1, 2)+math.Pow(y3, 2)-math.Pow(y1, 2))/(2*(y3-y1)) - ((x3-x1)/(y3-y1))*x

	d13 := math.Sqrt(math.Pow(x1-x3, 2) + math.Pow(y1-y3, 2))
	d23 := math.Sqrt(math.Pow(x2-x3, 2) + math.Pow(y2-y3, 2))

	//Se ultiliza una validacion para verificar que las distancias sean consistentes a la ubicacion de cada satelite
	if math.Abs(d13-d23) > d1+d3 || d13 > d1+d3 || d23 > d2+d3 || d13 == 0 || d23 == 0 {
		fmt.Printf("Pailos")
	} else {
		fmt.Printf("Location: (%.0f, %.0f)\n", x, y)
	}
	x = math.Round(x)
	y = math.Round(y)
	return x, y

}*/

func GetLocation(distances ...float64) (x, y float64) {

	var x1, x2, y2, x3, y1, y3 float64
	//Se extraen las posiciones conocidas de los satelites
	x1 = satellitesPosition[0].position_x
	y1 = satellitesPosition[0].position_y
	x2 = satellitesPosition[1].position_x
	y2 = satellitesPosition[1].position_y
	x3 = satellitesPosition[2].position_x
	y3 = satellitesPosition[2].position_y

	A := Point{x1, y1}
	B := Point{x2, y2}
	C := Point{x3, y3}

	distancias := []float64{}

	//Se leen las distancias de los satelites a las naves
	for _, t := range distances {
		//fmt.Println(t)
		distancias = append(distancias, t)
	}
	dA := distancias[0]
	dB := distancias[1]
	dC := distancias[2]
	/*
		dA := 538.52
		dB := 141.42
		dC := 509.9*/

	ex := divide(subtract(B, A), distance(subtract(B, A)))
	i := dot(ex, subtract(C, A))
	a := subtract(subtract(C, A), multiply(ex, i))
	ey := divide(a, distance(a))
	//ez := cross(ex, ey)
	d := distance(subtract(B, A))
	j := dot(ey, subtract(C, A))

	x = (math.Pow(dA, 2) - math.Pow(dB, 2) + math.Pow(d, 2)) / (2 * d)
	y = ((math.Pow(dA, 2) - math.Pow(dC, 2) + math.Pow(i, 2) + math.Pow(j, 2)) / (2 * j)) - ((i / j) * x)

	trilatX := A.X + (ex.X * x) + (ey.X * y)
	trilatY := A.Y + (ex.Y * x) + (ey.Y * y)

	fmt.Println(math.Round(trilatX), math.Round(trilatY))

	x = math.Round(trilatX)
	y = math.Round(trilatY)
	return x, y

}

func GetMessage(messages ...[]string) string {
	// Encontrar la longitud máxima de señal
	maxLen := 0
	for _, message := range messages {
		//fmt.Println(message)
		if len(message) > maxLen {
			maxLen = len(message)
		}
		//fmt.Prntln(maxLen)
	}

	// Construir una matriz de palabras para cada señal
	wordsMatrix := make([][]string, len(messages))
	for i, message := range messages {
		wordsMatrix[i] = make([]string, maxLen)
		for j := 0; j < maxLen; j++ {
			if j < len(message) {
				wordsMatrix[i][j] = message[j]
			} else {
				wordsMatrix[i][j] = ""
			}
		}
	}

	// Construir el mensaje secreto a partir de la matriz de palabras
	var secretMessage strings.Builder
	for j := 0; j < maxLen; j++ {
		wordCount := make(map[string]int)
		for i := range wordsMatrix {
			word := wordsMatrix[i][j]
			if word != "" {
				wordCount[word]++
			}
		}
		//fmt.Println(wordCount)
		var mostFrequentWord string
		var highestCount int
		for word, count := range wordCount {
			if count > highestCount {
				mostFrequentWord = word
				highestCount = count
			}
		}

		secretMessage.WriteString(mostFrequentWord)
		secretMessage.WriteRune(' ')
	}

	return strings.TrimSpace(secretMessage.String())
}

// Funcion que se ejecuta desde el cliente
func GetLocMess(w http.ResponseWriter, r *http.Request) {
	var satellite salellitesRecAll
	var satellitesSplit = satellitesRecBody{}
	reqBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Fprintf(w, "Datos de entarda erroneos")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)

	}

	json.Unmarshal(reqBody, &satellite)
	fmt.Println(satellite, "Satelitte")
	//Split de cada diccionario de datos relacionado a cada satelite
	for _, sat := range satellite.Satellites {
		fmt.Println(sat)
		satellitesSplit = append(satellitesSplit, sat)
	}

	d1 := satellitesSplit[0].Distance
	d2 := satellitesSplit[1].Distance
	d3 := satellitesSplit[2].Distance

	x, y := GetLocation(d1, d2, d3)
	fmt.Println("La cordenada X es: ", x, " La cordenada Y es: ", y)

	msg1 := satellitesSplit[0].Message
	msg2 := satellitesSplit[1].Message
	msg3 := satellitesSplit[2].Message

	msgRes := GetMessage(msg1, msg2, msg3)
	fmt.Println(GetMessage(msg1, msg2, msg3))

	respuestaCli.Position.X = x
	respuestaCli.Position.Y = y
	respuestaCli.Message = msgRes
	//fmt.Println(satellitesSplit)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(respuestaCli)
}

// Funcion que se ejecuta desde el cliente
func PostLocMessSplit(w http.ResponseWriter, r *http.Request) {
	var satellite satelliteRecsplit
	reqBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Fprintf(w, "Datos de entarda erroneos")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)

	}

	json.Unmarshal(reqBody, &satellite)

	//Split de cada diccionario de datos relacionado a cada satelite
	vars := mux.Vars(r)
	name := vars["satellite_name"]

	switch name {
	case "kenobi":
		dist1 = satellite.Distance
		msgS1 = satellite.Message
	case "skywalker":
		dist2 = satellite.Distance
		msgS2 = satellite.Message
	case "sato":
		dist3 = satellite.Distance
		msgS3 = satellite.Message
	default:
		w.WriteHeader(http.StatusForbidden)
	}

	//fmt.Println(satellitesSplit)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(satellite)
}

// Funcion que se ejecuta desde el cliente
func GetLocMessSplit(w http.ResponseWriter, r *http.Request) {

	//Split de cada diccionario de datos relacionado a cada satelite

	//fmt.Println(satellitesSplit)
	w.Header().Set("Content-Type", "application/json")
	//Validamos que se hayan llenado los valores de las distancias para poder utilizar el metodo GET
	if dist1 == 0 || dist2 == 0 || dist3 == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode("No hay informacion suficiente")
	} else {
		fmt.Println(dist1)
		fmt.Println(dist2)
		fmt.Println(dist3)
		x, y := GetLocation(dist1, dist2, dist3)
		fmt.Println("La cordenada X es: ", x, " La cordenada Y es: ", y)
		msgRes := GetMessage(msgS1, msgS2, msgS3)
		fmt.Println(GetMessage(msgS1, msgS2, msgS3))

		respuestaCli.Position.X = x
		respuestaCli.Position.Y = y
		respuestaCli.Message = msgRes
		//fmt.Println(satellitesSplit)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(respuestaCli)

	}

}

func GetLocMessSplitDelete(w http.ResponseWriter, r *http.Request) {

	//Split de cada diccionario de datos relacionado a cada satelite

	dist1 = 0
	dist2 = 0
	dist3 = 0
	msgS1 = nil
	msgS2 = nil
	msgS3 = nil
	//fmt.Println(satellitesSplit)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Los datos fueron eliminados")

}

func main() {

	//Rutas
	router := mux.NewRouter().StrictSlash(true)
	//router.HandleFunc("/topsecret/", GetLoc)
	router.HandleFunc("/topsecret/", GetLocMess).Methods("POST")
	router.HandleFunc("/topsecret_split/{satellite_name}", PostLocMessSplit).Methods("POST")
	router.HandleFunc("/topsecret_split/", GetLocMessSplit).Methods("GET")
	router.HandleFunc("/topsecret_split/", GetLocMessSplitDelete).Methods("DELETE")
	//Server

	log.Fatal(http.ListenAndServe(":8080", router))
}

func subtract(a, b Point) Point {
	return Point{a.X - b.X, a.Y - b.Y}
}

func dot(a, b Point) float64 {
	return (a.X * b.X) + (a.Y * b.Y)
}

func distance(a Point) float64 {
	return math.Sqrt(math.Pow(a.X, 2) + math.Pow(a.Y, 2))
}

func divide(a Point, b float64) Point {
	return Point{a.X / b, a.Y / b}
}

func cross(a, b Point) Point {
	return Point{(a.Y * b.X) - (a.X * b.Y), (a.X * b.Y) - (a.Y * b.X)}
}

func multiply(a Point, b float64) Point {
	return Point{a.X * b, a.Y * b}
}
