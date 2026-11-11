# Aggregator Blog CLI Tool
Hi! Thanks for checking my project out :) 
I made this purely for educational reasons, so please be sure to check out the below warning. 

>[!warning]
> This documentation is probably incomplete and may lead to increasing frustration and laptop defenestration. 

<3
## Prerequisites
You must have Go and Postgres (v15+) installed .
## Setting up the Database
> [!important]
> For the sake of documentation and project submission purposes, I am including these steps here, but I will try to figure out how to streamline this process later down the line. I apologize for the inconvenience.

> [!caution]
> oh yeah, and the below steps only works for Linux. Who the f is steve jobs? and don't even get me started with microsoft.

> [!tip]
> damn these markdown boxes are fun with their own colors lol

1. Start the server:
`sudo service postgresql start`
2. Enter the shell:
`sudo -u postgres psql`
3. Create the database:
`CREATE DATABASE gator;`
4. Connect to the database (still in the shell):
`\c gator`
5. Create a password:
`ALTER USER <user> PASSWORD <password>`
6. Exit the shell
7. Test your connection string:
`psql 'postgres://<user>:<password>@localhost:5432/gator'`
8. Create a `.gatorconfig.json` in your home directory:
9. Add the connection string to your gatorconfig.json in your home but with the extra string `?sslmode=disable`. Your json config should look like:
```json
{
	"db_url": "postgres://<user>:<password>@localhost:5432/gator?sslmode=disable"
}
```
10. Save the file. Now your database is set up! Phew.

## Install the Project
Now we get to the meat of it all.

After cloning this repository, run `go install .` from the root of the project.

Run `gator-blog` and if an output like `not enough arguments` occurs, then it's installed!
## Now What? Commands!
Cool! Let's run some commands :)
#### register
`gator-blog register <kimi-no-nawa>`\
This adds a user to your database.
Try adding a few users :)
#### users
`gator-blog users`\
This lists the users
You will see which user is currently logged in.
#### login
`gator-blog login <name>`\
This will login to the specified user.
#### addfeed
`gator-blog addfeed <name-of-feed> <url>`\
This will add the specified feed

>[!note]
Example. If you want to add the RSS feed of my blog (:o), you wold do something like
`gator-blog addfeed "Miranda's Blog" "https://mirmeep.github.io/blog/posts/index.xml"`

ain't that neato?

Add a few more feeds!
#### feeds
`gator-blog feeds`\
This will list the feeds of all users.
#### follow
`gator-blog follow <url>`\
This allows the currently logged in user to follow a feed that is already in the database.
#### following
`gator-blog following`\
This lists the feeds the currently logged in user is following.
#### unfollow
`gator-blog unfollow <url>`\
Un-follows a feed that the currently logged in user is following.
### agg
`gator-blog agg <time-string>`\
This aggregates blogs! The `<time-string>` can be something like `5s`
#### browse
`gator-blog browse <optional-limit>`\
After aggregating, this allows the currently logged in user to see their feeds! 

