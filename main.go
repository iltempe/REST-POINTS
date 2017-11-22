//realized analizing and change this experiment https://gist.github.com/EtienneR/ed522e3d31bc69a9dec3335e639fcf60
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

type Pois struct {
	Id  int    `gorm:"AUTO_INCREMENT" form:"id" json:"id"`
	lat string `gorm:"not null" form:"lat" json:"lat"`
	lon string `gorm:"not null" form:"lon" json:"lon"`
}

func InitDb() *gorm.DB {
	// Openning file
	db, err := gorm.Open("sqlite3", "./data.db")
	// Display SQL queries
	db.LogMode(true)

	// Error
	if err != nil {
		panic(err)
	}
	// Creating the table
	if !db.HasTable(&Pois{}) {
		db.CreateTable(&Pois{})
		db.Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(&Pois{})
	}

	return db
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.Next()
	}
}

func main() {
	r := gin.Default()

	r.Use(Cors())

	v1 := r.Group("api/v1")
	{
		v1.POST("/pois", PostPoi)
		v1.GET("/pois", GetPois)
		v1.GET("/pois/:id", GetPoi)
		v1.PUT("/pois/:id", UpdatePoi)
		v1.DELETE("/pois/:id", DeletePoi)
		v1.OPTIONS("/pois", OptionsPoi)     // POST
		v1.OPTIONS("/pois/:id", OptionsPoi) // PUT, DELETE
	}

	r.Run(":8500")
}

func PostPoi(c *gin.Context) {
	db := InitDb()
	defer db.Close()

	var poi Pois
	c.Bind(&poi)

	if poi.lat != "" && poi.lon != "" {
		db.Create(&poi)
		// Display error
		c.JSON(201, gin.H{"success": poi})
	} else {
		// Display error
		c.JSON(422, gin.H{"error": "Fields are empty"})
		//curl -i -X POST -H "Content-Type: application/json" -d "{ \"lat\": \"42\", \"lon\": \"54\" }" http://localhost:8500/api/v1/pois
	}
}

func GetPois(c *gin.Context) {
	// Connection to the database
	db := InitDb()
	// Close connection database
	defer db.Close()

	var pois []Pois
	// SELECT * FROM users
	db.Find(&pois)

	// Display JSON result
	c.JSON(200, pois)

	// curl -i http://localhost:8080/api/v1/users
}

func GetPoi(c *gin.Context) {
	// Connection to the database
	db := InitDb()
	// Close connection database
	defer db.Close()

	id := c.Params.ByName("id")
	var poi Pois
	// SELECT * FROM users WHERE id = 1;
	db.First(&poi, id)

	if poi.Id != 0 {
		// Display JSON result
		c.JSON(200, poi)
	} else {
		// Display JSON error
		c.JSON(404, gin.H{"error": "Poi not found"})
	}

}

func UpdatePoi(c *gin.Context) {
	// Connection to the database
	db := InitDb()
	// Close connection database
	defer db.Close()

	// Get id
	id := c.Params.ByName("id")
	var poi Pois
	// SELECT * FROM users WHERE id = 1;
	db.First(&poi, id)

	if poi.lat != "" && poi.lon != "" {

		if poi.Id != 0 {
			var newPoi Pois
			c.Bind(&newPoi)

			result := Pois{
				Id:  poi.Id,
				lat: newPoi.lat,
				lon: newPoi.lon,
			}

			// UPDATE users SET
			db.Save(&result)
			// Display modified data in JSON message "success"
			c.JSON(200, gin.H{"success": result})
		} else {
			// Display JSON error
			c.JSON(404, gin.H{"error": "Poi not found"})
		}

	} else {
		// Display JSON error
		c.JSON(422, gin.H{"error": "Fields are empty"})
	}
}

func DeletePoi(c *gin.Context) {
	// Connection to the database
	db := InitDb()
	// Close connection database
	defer db.Close()

	// Get id
	id := c.Params.ByName("id")
	var poi Pois
	// SELECT * FROM users WHERE id = 1;
	db.First(&poi, id)

	if poi.Id != 0 {
		// DELETE FROM users WHERE id = user.Id
		db.Delete(&poi)
		// Display JSON result
		c.JSON(200, gin.H{"success": "Poi #" + id + " deleted"})
	} else {
		// Display JSON error
		c.JSON(404, gin.H{"error": "Poi not found"})
	}
}

func OptionsPoi(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Methods", "DELETE,POST, PUT")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	c.Next()
}
