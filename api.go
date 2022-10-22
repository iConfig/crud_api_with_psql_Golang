package main 

import (
	"fmt"
	"time"
	"net/http"
	"encoding/json"
	"database/sql"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// database connection details as constant 
const (
	DB_host = "localhost"
	DB_port = 5432
	DB_user = "postgres"
	DB_pass = "postgres"
	DB_name = "mydata"

)

// error check function
func checkErr(err error){
	if err != nil {
		panic(err)
	}
}

// setup database
func DBsetup() *sql.DB {
	DBconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
							DB_host,DB_port,DB_user,DB_pass,DB_name)

	DB , err := sql.Open("postgres", DBconn)

	checkErr(err)
	
	return DB
}

// messages function 
func Messages(message string){
	fmt.Println(" ")
	fmt.Println(message)
	fmt.Println(" ")
}

// Center struct
type Center struct {
	CenterID  int `json:"id"`
	CenterName string `json:"name"`
	CenterLocation string `json:"location"`
	CenterContact string `json:"contact"`
}

// JsonResponse struct
type JsonResponse struct{
	Type string `json:"type"`
	Data []Center `json:"data"`
	Message string `json:"message"`
}



// router function 
func Handler(){
	router := mux.NewRouter()

	// GET all centers router
	router.HandleFunc("/centers", GetCenters).Methods("GET")
	
	// ADD center
	router.HandleFunc("/addcenters", AddCenters).Methods("POST")

	// DELETE a center 
	router.HandleFunc("/delcenter/{centerId}", DeleteCenter).Methods("DELETE")

	// DELETE all centers 
	router.HandleFunc("/delcenters", DeleteCenters).Methods("DELETE") 

	// server 

	currentTime := time.Now()

	fmt.Println(currentTime,": API server running on 127.0.0.1:8090 .......\n ")
	err := http.ListenAndServe("127.0.0.1:8090", router)

	checkErr(err)

}

// Get all centers
func GetCenters(w http.ResponseWriter, r *http.Request){
	DB := DBsetup()

	Messages("Getting all centers...")

	// Get all centers from table 
	rows , err := DB.Query(`SELECT * FROM centers`)

	checkErr(err)

	var centers []Center

	for rows.Next(){

		var center_id int 
		var name string
		var location string 
		var contact string

		err = rows.Scan(&center_id, &name , &location , &contact)

		checkErr(err)

		centers = append(centers, Center{CenterID: center_id, CenterName: name,
						CenterLocation:location, CenterContact:contact})
	} 

	var response = JsonResponse{Type:"success", Data:centers}

	json.NewEncoder(w).Encode(response)

}

// Add to centers function 
func AddCenters(w http.ResponseWriter, r *http.Request){

	CenterName := r.FormValue("name")
	CenterLocation := r.FormValue("location")
	CenterContact := r.FormValue("contact")

	var response = JsonResponse{}

	if CenterName == "" || CenterLocation == "" || CenterContact == "" {
		response = JsonResponse{Type: "error", Message:"You are missing a required field"}
	} else {
		DB := DBsetup()

		Messages("Adding centers into Database")

		fmt.Println("Adding new center: " + CenterName + " located at " + CenterLocation + " with contact number " + CenterContact)
		
		var lastInsertID int 
		err := DB.QueryRow(`INSERT INTO centers (name,location,contact) VALUES ($1,$2,$3) returning center_id;`,CenterName,CenterLocation,CenterContact).Scan(&lastInsertID)
		checkErr(err)

		response = JsonResponse{Type: "success", Message:"New center added successfully!"}
	}
	json.NewEncoder(w).Encode(response)

}

//Delete a center function
func DeleteCenter(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)

	centerID := vars["id"]

	var response = JsonResponse{}

	if centerID == "" {
		response = JsonResponse{Type:"error", Message:"You did not provide an ID! "}		
	} else {
		DB := DBsetup()

		Messages("Deleting center from Database...")
		_, err := DB.Exec("DELETE FROM centers WHERE center_id = $1", centerID)
		checkErr(err)

		response = JsonResponse{Type:"success", Message:"center has been deleted successfully!"}		
	}
	json.NewEncoder(w).Encode(response)
}

// DELETE all centers 
func DeleteCenters(w http.ResponseWriter, r *http.Request){
	DB := DBsetup()

	var response = JsonResponse{}

	Messages("Deleting all centers...")

	_, err := DB.Exec("DELETE * FROM centers")
	checkErr(err)

	Messages("All centers deleted successfully! ")

	response = JsonResponse{Type:"success", Message:"All centers deleted successfully "}
	json.NewEncoder(w).Encode(response)
}


func main(){
	Handler()
}