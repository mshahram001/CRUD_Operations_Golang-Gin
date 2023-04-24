package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func db_creation() {
	db, err := sql.Open("mysql", "root:india@123@tcp(127.0.0.1:3306)/")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	//  Create Database
	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS http_gin")
	if err != nil {
		panic(err.Error())
	}
}

type info struct {
	ID   int
	Name string
	Age  int
	City string
}

// Get
func get_user(c *gin.Context) {
	db, err := sql.Open("mysql", "root:india@123@tcp(127.0.0.1:3306)/http_gin")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// Read Data
	results, err := db.Query("SELECT * FROM database1")
	if err != nil {
		panic(err.Error())
	}
	defer results.Close()

	var output interface{}
	for results.Next() {
		var id int
		var name string
		var age int
		var city string

		err = results.Scan(&id, &name, &age, &city)
		if err != nil {
			panic(err.Error())
		}
		output = fmt.Sprintf("%d  %s  %d  %s ", id, name, age, city)
		c.IndentedJSON(http.StatusOK, output)
	}
}

// POST
func post_user(c *gin.Context) {
	db, err := sql.Open("mysql", "root:india@123@tcp(127.0.0.1:3306)/http_gin")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	var data info
	err = c.BindJSON(&data)
	if err != nil {
		return
	}
	c.IndentedJSON(http.StatusCreated, data)
	query_data := fmt.Sprintf(`INSERT INTO database1 VALUES(%d,"%s",%d,"%s")`, data.ID, data.Name, data.Age, data.City)
	fmt.Println(query_data)

	//insert data
	insert, err := db.Query(query_data)
	if err != nil {
		panic(err.Error())
	}

	defer insert.Close()
	fmt.Println("Yay, values added!")
}

// Put
func update_user(c *gin.Context) {

	db, err := sql.Open("mysql", "root:india@123@tcp(127.0.0.1:3306)/http_gin")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// Get user input
	var data info
	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update data
	query_data := fmt.Sprintf("UPDATE database1 SET Name='%s', Age=%d, City='%s' WHERE ID=%d", data.Name, data.Age, data.City, data.ID)
	if _, err := db.Exec(query_data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

// Delete
func delete_user(c *gin.Context) {
	db, err := sql.Open("mysql", "root:india@123@tcp(127.0.0.1:3306)/http_gin")
	if err != nil {
		panic(err.Error())
	}
	var record info
	if err := c.BindJSON(&record); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("Record ID", record.ID)
	defer db.Close()
	query_data := fmt.Sprintf("DELETE FROM database1 WHERE ID = %d", record.ID)
	if _, err := db.Exec(query_data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func main() {
	db_creation()
	// DB connectivity
	db, err := sql.Open("mysql", "root:india@123@tcp(127.0.0.1:3306)/http_gin")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	// Create table
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS database1(ID INT NOT NULL AUTO_INCREMENT, Name VARCHAR(20),Age INT, City VARCHAR(20), PRIMARY KEY (id) );")
	if err != nil {
		panic(err.Error())
	}
	router := gin.Default()
	router.POST("/post", post_user)
	router.GET("/get", get_user)
	router.PUT("/update", update_user)
	router.DELETE("/delete", delete_user)
	router.Run("localhost:8080")
}
