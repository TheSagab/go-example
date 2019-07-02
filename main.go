package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

func init() {
	//open a db connection
	//uses MySQL with database named go-web, username root and password root
	var err error
	db, err = gorm.Open("mysql", "root:root@/go-web?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}

	//Migrate the schema
	db.AutoMigrate(&personModel{})
}

type personModel struct {
		gorm.Model
		Name string `json:"name"`
	}

func createPerson(c *gin.Context) {
	//person := personModel{Name: c.PostForm("name")}
	var person personModel
	c.BindJSON(&person)
	db.Save(&person)
	c.JSON(http.StatusCreated, gin.H{
		"status": http.StatusCreated, "message": "Person created successfully!", "resourceId": person.ID})
}

func fetchAllPerson(c *gin.Context) {
	var persons []personModel

	db.Find(&persons)

	if len(persons) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No person found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": persons})
}

func fetchSinglePerson(c *gin.Context) {
	var person personModel
	personID := c.Param("id")

	db.First(&person, personID)

	if person.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No person found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": person})
}

func updatePerson(c *gin.Context) {
	var person personModel
	personID := c.Param("id")

	db.First(&person, personID)

	if person.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No person found!"})
		return
	}

	// todo: use a better way to get person name
	var personFromJson personModel
	c.BindJSON(&personFromJson)

	db.Model(&person).Update("name", personFromJson.Name)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "person updated successfully!"})
}

func deletePerson(c *gin.Context) {
	var person personModel
	personID := c.Param("id")

	db.First(&person, personID)

	if person.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No person found!"})
		return
	}

	db.Delete(&person)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Person deleted successfully!"})
}
func main() {
	router := gin.Default()

	v1 := router.Group("/api/v1/persons")
	{
		v1.POST("/", createPerson)
		v1.GET("/", fetchAllPerson)
		v1.GET("/:id", fetchSinglePerson)
		v1.PUT("/:id", updatePerson)
		v1.DELETE("/:id", deletePerson)
	}
	router.Run()
	// Listen and Server in 0.0.0.0:8080

	router.Run(":8080")
}
