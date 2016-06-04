// by setting package as main, Go will compile this as an executable file.
// Any other package turns this into a library
package main

// These are your imports / libraries / frameworks
import (
	// this is Go's built-in sql library
	"database/sql"
	"log"
	"net/http"
	"os"
	"fmt"
	//"strconv"

	// this allows us to run our web server
	"github.com/gin-gonic/gin"
	// this lets us connect to Postgres DB's
	_ "github.com/lib/pq"
)

var (
	// this is the pointer to the database we will be working with
	// this is a "global" variable (sorta kinda, but you can use it as such)
	db *sql.DB
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	var errd error
	// here we want to open a connection to the database using an environemnt variable.
	// This isn't the best technique, but it is the simplest one for heroku
	db, errd = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if errd != nil {
		log.Fatalf("Error opening database: %q", errd)
	}
	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("html/*")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	router.GET("/ping", func(c *gin.Context) {
		ping := db.Ping()
		if ping != nil {
			// our site can't handle http status codes, but I'll still put them in cause why not
			c.JSON(http.StatusOK, gin.H{"error": "true", "message": "db was not created. Check your DATABASE_URL"})
		} else {
			c.JSON(http.StatusOK, gin.H{"error": "false", "message": "db created"})
		}
	})

	router.GET("/query1", func(c *gin.Context) {
		table := "<table class='table'><thead><tr>"
		rows, err := db.Query("SELECT name AS Parks FROM park NATURAL JOIN favorite NATURAL JOIN usr WHERE usr.username = 'rdiaz0'")
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		// foreach loop over rows.Columns, using value
		cols, _ := rows.Columns()
		if len(cols) == 0 {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		for _, value := range cols {
			table += "<th class='text-center'>" + value + "</th>"
		}
		table += "</thead><tbody>"
		var name string
		for rows.Next() {
			rows.Scan(&name)
			// can't combine ints and strings in Go. Use strconv.Itoa(int) instead
			table += "<tr><td>" + name + "</td></tr>"
		}
		// finally, close out the body and table
		table += "</tbody></table>"
		c.Data(http.StatusOK, "text/html", []byte(table))
	})

	router.GET("/query2", func(c *gin.Context) {
		table := "<table class='table'><thead><tr>"
		
		rows, err := db.Query("SELECT trailName AS Trails FROM trail JOIN favorite ON trail.trailId = favorite.trailId JOIN usr ON favorite.username = usr.username WHERE usr.username = 'rdiaz0'")
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		// foreach loop over rows.Columns, using value
		cols, _ := rows.Columns()
		if len(cols) == 0 {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		for _, value := range cols {
			table += "<th class='text-center'>" + value + "</th>"
		}
		table += "</thead><tbody>"
		// columns
		var name string
		for rows.Next() {
			rows.Scan(&name)
			table += "<tr><td>" + name + "</td></tr>" 
		}
		table += "</tbody></table>"
		c.Data(http.StatusOK, "text/html", []byte(table))
	})

	router.POST("/submit", func(c *gin.Context) {
		description := c.PostForm("description")
		rating := c.PostForm("rating")
		warnings := c.PostForm("warnings")

		stmt, err := db.Prepare("INSERT INTO review (username, parkId, postdate, detail, rating, warning) VALUES ('rdiaz0', 3, CURRENT_DATE, $1, $2, $3);")
		res, err := stmt.Exec(description, rating, warnings)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		fmt.Println(res)

	})

	router.POST("/search", func(c *gin.Context) {
		table := "<table class='table' id='currRes'><thead><tr>"
		search := c.PostForm("searchTerm")

		searchTerm := "%" + search + "%"

		rows, err := db.Query("SELECT park.name AS results FROM park WHERE park.name LIKE $1", searchTerm)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		cols, _ := rows.Columns()
		if len(cols) == 0 {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		for _, value := range cols {
			table += "<th class='text-center'>" + value + "</th>"
		}
		table += "</thead><tbody>"
		// columnss
		var name string
		for rows.Next() {
			rows.Scan(&name)
			table += "<tr><td>" + name + "</td></tr>" 
		}
		table += "</tbody></table>"
		c.Data(http.StatusOK, "text/html", []byte(table))
	})

	// NO code should go after this line. it won't ever reach that point
	router.Run(":" + port)
}

/*
Example of processing a GET request

// this will run whenever someone goes to last-first-lab7.herokuapp.com/EXAMPLE
router.GET("/EXAMPLE", func(c *gin.Context) {
    // process stuff
    // run queries
    // do math
    //decide what to return
    c.JSON(http.StatusOK, gin.H{
        "key": "value"
        }) // this returns a JSON file to the requestor
    // look at https://godoc.org/github.com/gin-gonic/gin to find other return types. JSON will be the most useful for this
})

*/
