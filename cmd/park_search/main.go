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
	"strconv"

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
		// put your query here
		rows, err := db.Query("SELECT person.firstName, song.title, song.length FROM favorite JOIN person ON person.personId = favorite.personId JOIN song ON song.songId = favorite.songId WHERE song.length = (SELECT max(song.length) FROM song)") // <--- EDIT THIS LINE
		if err != nil {
			// careful about returning errors to the user!
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
		// once you've added all the columns in, close the header
		table += "</thead><tbody>"
		// declare all your RETURNED columns here
		var name string      // <--- EDIT THESE LINES
		var title string //<--- ^^^^
		var length int
		for rows.Next() {
			// assign each of them, in order, to the parameters of rows.Scan.
			// preface each variable with &
			rows.Scan(&name, &title, &length) // <--- EDIT THIS LINE
			// can't combine ints and strings in Go. Use strconv.Itoa(int) instead
			table += "<tr><td>" + name + "</td><td>" + title + "</td><td>" + strconv.Itoa(length) + "</td></tr>" // <--- EDIT THIS LINE
		}
		// finally, close out the body and table
		table += "</tbody></table>"
		c.Data(http.StatusOK, "text/html", []byte(table))
	})

	router.GET("/query2", func(c *gin.Context) {
		table := "<table class='table'><thead><tr>"
		// put your query here
		rows, err := db.Query("SELECT person.firstName, person.favoriteGenre, album.genre FROM person JOIN favorite ON person.personId = favorite.personId JOIN song ON favorite.songId = song.songId JOIN album ON song.albumId = album.albumId WHERE album.genre <> person.favoriteGenre") // <--- EDIT THIS LINE
		if err != nil {
			// careful about returning errors to the user!
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
		// once you've added all the columns in, close the header
		table += "</thead><tbody>"
		// columns
		var name string
		var favoriteGenre string
		var songGenre string
		for rows.Next() {
			rows.Scan(&name, &favoriteGenre, &songGenre) // put columns here prefaced with &
			table += "<tr><td>" + name + "</td><td>" + favoriteGenre + "</td><td>" + songGenre + "</td></tr>" // <--- EDIT THIS LINE
		}
		// finally, close out the body and table
		table += "</tbody></table>"
		c.Data(http.StatusOK, "text/html", []byte(table))
	})

	router.GET("/query3", func(c *gin.Context) {
		table := "<table class='table'><thead><tr>"
		// put your query here
		rows, err := db.Query("SELECT artist.firstName, artist.lastName FROM artist WHERE artist.artistID IN (SELECT album.artistId FROM album WHERE genre = 'Country')") // <--- EDIT THIS LINE
		if err != nil {
			// careful about returning errors to the user!
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
		// once you've added all the columns in, close the header
		table += "</thead><tbody>"
		// columns
		var fName string
		var lName string
		for rows.Next() {
			rows.Scan(&fName, &lName) // put columns here prefaced with &
			table += "<tr><td>" + fName + "</td><td>" + lName + "</td></tr>" // <--- EDIT THIS LINE
		}
		// finally, close out the body and table
		table += "</tbody></table>"
		c.Data(http.StatusOK, "text/html", []byte(table))
	})

	router.POST("/login", func(c *gin.Context) {
		// this is meant for SQL injection examples ONLY.
		// Don't copy this for use in an actual environment, even if you do stop SQL injection
		username := c.PostForm("username")
		password := c.PostForm("password")

		rows, err := db.Query("SELECT usr.name FROM usr WHERE usr.name = $1 AND usr.password = $2;", username, password)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		cols, _ := rows.Columns()
		if len(cols) == 0 {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		rowCount := 0
		var resultUser string
		for rows.Next() {
			rows.Scan(&resultUser)
			rowCount++
		}
		// quick way to check if the user logged in properly
		if rowCount == 0 {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		// instead of HTML, we are going to return a JSON file
		c.JSON(http.StatusOK, gin.H{"username": resultUser})
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