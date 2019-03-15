# Simple TODO app built in GoLang using PostgreSQL for data storage.
***

## About This Application
1. GoLang application is hosted on **AWS EC2 Instance**.
2. PostgreSQL database is hosted on **AWS RDS**.
3. API documentation is currently hosted on **Swagger Hub** (https://app.swaggerhub.com/apis/karolispx/golang-rh-todo/1.0.0-oas3/) and can be used for making all API calls that are available in the application. More information contained in links below. 
4. **docs** folder in the repository also contains HTML generated Swagger docs and Postman collection in JSON format.

***

## Features & Endpoints
**Complete Swagger Documentation:**


### [User Authentication](https://github.com/karolispx/golang-rh-todo/wiki/3.-User-Authentication)
1. [Registration](https://github.com/karolispx/golang-rh-todo/wiki/3.-User-Authentication#1-registration)
2. [Login](https://github.com/karolispx/golang-rh-todo/wiki/3.-User-Authentication#2-login)
### [Tasks](https://github.com/karolispx/golang-rh-todo/wiki/4.-Tasks)
1. [Get tasks](https://github.com/karolispx/golang-rh-todo/wiki/4.-Tasks#1-get-tasks)
2. [Create a task](https://github.com/karolispx/golang-rh-todo/wiki/4.-Tasks#2-create-a-task)
3. [Delete all tasks](https://github.com/karolispx/golang-rh-todo/wiki/4.-Tasks#3-delete-all-tasks)
4. [Get a specific task](https://github.com/karolispx/golang-rh-todo/wiki/4.-Tasks#4-get-a-specific-task)
5. [Update a specific task](https://github.com/karolispx/golang-rh-todo/wiki/4.-Tasks#5-update-a-specific-task)
6. [Delete a specific task](https://github.com/karolispx/golang-rh-todo/wiki/4.-Tasks#6-delete-a-specific-task)
7. [Watch a specific task](https://github.com/karolispx/golang-rh-todo/wiki/4.-Tasks#7-watch-a-specific-task)
8. [Unwatch a specific task](https://github.com/karolispx/golang-rh-todo/wiki/4.-Tasks#8-unwatch-a-specific-task)

***

## Project Setup
* [Project Setup Locally](https://github.com/karolispx/golang-rh-todo/wiki/1.-Project-Setup-Locally)
* [Project Setup On A Server](https://github.com/karolispx/golang-rh-todo/wiki/2.-Project-Setup-On-A-Server)

***

## Database Tables
```
CREATE TABLE users (
    userid SERIAL PRIMARY KEY,
    email_address text NOT NULL,
    password text NOT NULL,
    last_action text NOT NULL,
    date_created text NOT NULL
);


CREATE TABLE tasks (
    taskid SERIAL PRIMARY KEY,
    userid integer NOT NULL,    
    task text NOT NULL,
    watching text NOT NULL,
    date_created text NOT NULL,
    date_updated text NOT NULL
);
```
