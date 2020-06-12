# GoNotes

GoNotes is a self hostable alternative to sites like Hastebin / Pastebin. It is no where near as feature rich nor supported as Hastebin nor Pastebin. This project was an easy way for me to learn about web servers in Golang.

## Installation

Set up a postgresql server, change the appropriate database connection variables in `/src/models/database.go` and `/database_cleaner/database_cleaner.go` *(Or adapt this project to use environment / config file variables)*

The database schema can be found in `/data/database.sql`

Clone this repo and run
```bash
cd src/
go build main.go
```
Then run the appropriate output file.

You'll need to download and move the appropiate [ACE Editor](https://github.com/ajaxorg/ace-builds/) and [Bootstrap](https://getbootstrap.com/docs/4.5/getting-started/download/) releases to the `/src/public/css` and `/src/public/js` folders

## Database cleaner

The database_cleaner is intended to be run as a cron job, or just wheneveri you need to use it, it will remove all posts that haven't been viewed in the last 15days and, when the project supports it, all *expired* posts.

Remember to build it using, also remember to change the database connection variables *(Or adapt this project to use environment / config file variables)* to the same as you have for the web app.
```bash
cd database_cleaner/
go build database_cleaner.go
```

## Contributing and Issues

Pull requests are always welcome.
This was only a small project for me to learn the basics of a web server in Golang so forgive me if my responses to issues / PRs are delayed or nonexistent.
