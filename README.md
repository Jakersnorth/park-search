Group 1: Park Search Website

To see the db schema and example data inserts, look in the databaseInfo folder

We created three stored procedures for the database. Two of them are used for inserting new entries in location and favorite while ensuring that the new entry only references one of the two optional foreign keys of parkId and trailId. These two procedures guarantee that there will never be an entry that connects to both a park and a trail.

The last stored procedure is for adding new parks to the database with a location included as well. This procedure helps reduce the number of necesary insert statements by filling in both the park table and the location table.

The view we made creates a table of all parks along with the most recent review added for that park. The purpose of this view is for when displaying search results of parks so that users can view a relevant review about the park without going to the park page.

Website: http://group1-park-search.herokuapp.com/
