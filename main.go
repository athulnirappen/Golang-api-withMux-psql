package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/joho/godotenv"
)

// create product
func createProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	//db connect
	db := connect()
	defer db.Close()

	//creating product instance
	product := &Product{
		ID: uuid.New().String(),
	}

	//decode request
	_ = json.NewDecoder(r.Body).Decode(&product)

	// insert into data base
	_, err := db.Model(product).Insert()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//return product

	json.NewEncoder(w).Encode(&product)

}

func getProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	//db connect
	db := connect()
	defer db.Close()

	//create product slice
	var products []Product
	if err := db.Model(&products).Select(); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return

	}

	//getting products 

	json.NewEncoder(w).Encode(products)
}

func getSingleProduct(w http.ResponseWriter,r *http.Request)  {
	w.Header().Set("Content-Type", "application/json")

	//db connect
	db := connect()
	defer db.Close()

	//get id

	params:=mux.Vars(r)
	productId:=params["id"]

	product:=&Product{ID: productId}
	if err:=db.Model(product).WherePK().Select();err!=nil{
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//getting product 

	json.NewEncoder(w).Encode(product)
}

func deleteProduct(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Content-Type", "application/json")

	//db connect
	db := connect()
	defer db.Close()

	params:=mux.Vars(r)
	productId:=params["id"]
 
	//create product instance alternative way
	// product:=&Product{ID: productId}

	// result,err:=db.Model(product).WherePK().Delete()

	// create product instance 
	product:=&Product{}

	result,err:=db.Model(product).Where("id=?",productId).Delete()
	if err!=nil{
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//getting product 

	json.NewEncoder(w).Encode(result)
	
}

func updateProduct(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    // db connect
    db := connect()
    defer db.Close()

    // Get product ID from URL
    params := mux.Vars(r)
    productId := params["id"]

    // Create product instance with the ID
    product := &Product{ID: productId}

    // Decode the incoming request body into the product struct
    if err := json.NewDecoder(r.Body).Decode(product); err != nil {
        log.Println("Error decoding request body:", err)
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request data"})
        return
    }

    // Update the product in the database
    _, err := db.Model(product).WherePK().
        Set("name = ?", product.Name).
        Set("quantity = ?", product.Quantity).
        Set("price = ?", product.Price).
        Set("store = ?", product.Store).
        Update()
    if err != nil {
        log.Println("Error updating product:", err)
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{"error": "Error updating product"})
        return
    }

    // Return the updated product
    json.NewEncoder(w).Encode(product)
}


func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(uuid.New())

	r := mux.NewRouter()
	r.HandleFunc("/api/products", createProduct).Methods("POST")
	r.HandleFunc("/api/products", getProducts).Methods("GET")
	r.HandleFunc("/api/products/{id}", getSingleProduct).Methods("GET")
	r.HandleFunc("/api/products/{id}", updateProduct).Methods("PUT")
	r.HandleFunc("/api/products/{id}", deleteProduct).Methods("DELETE")
	// db:=connect()
	// defer db.Close()

	// product := Product{ID: uuid.New(),Name: "My product",Quantity: 23,Price: 34.5}
	// fmt.Println(product)

	//server create

	log.Fatal(http.ListenAndServe(":8000", r))
}
