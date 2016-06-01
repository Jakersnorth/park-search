# Lab "8"

For this "lab" you won't be turning anything in, but rather you will be working on learning more about SQL injection and how to prevent it.

First off, if you are seeing this you already cloned down your repo, and pulled in my changes.

Complete the following steps before working:
### Database
1. Open up a connection to your database, using the command you used in lab 7, which is found in your database connection settings under "Psql"
2. Run the following commands in order:
  * `CREATE TABLE usr(id int PRIMARY KEY, name varchar(200), password varchar(200));`
  * `INSERT INTO usr VALUES(1, 'admin', 'fjciS3905FjvosKL23fs'),(2,'me','pass'),(3,'another_user','password12345');`
  * `SELECT * FROM usr;` (This is just to check that you got the data inserted)
3. If this is a lab computer or you haven't worked on this computer before, don't forget to add the file `.env` and place inside of that `DATABASE_URL=` and then your database url.
4. Run `heroku local` and open a new tab in your browser, pointing to `localhost:5000`
5. Try logging in, using `me` and `pass` you should see that you logged in as "me"
6. Try using SQL Injection to bypass the login restriction, logging in as "admin". If you followed the slides you should see "Logged in as another_user". Try looking at the Go code, and thinking about what you did to your query to get this result instead of simply logging in as "me"
7. Retry step 6, except modifying your injection to reliably login as a specific user
8. Try inserting some data into the table, or dropping the whole table (Remember **THIS IS PERMANENT!**)
9. Modify the go code to protect against SQL Injection of this form [Check out the docs for the library we are using.](https://godoc.org/github.com/lib/pq#hdr-Queries) (Hint: The only way they describe how to handle input is using parameterized queries. The function Query can take an arbitrary number of arguments, and order matters)
10. Check to make sure that your change actually prevented an attack.

Nothing to submit for this lab, if you got this far you are well on your way to learning how to secure your database!
