# Simple TODO app built in GoLang using PostgreSQL for data storage.
***

## About This Application:
1. GoLang application is hosted on **Openshift** (there's also instructions on how to host in on **AWS EC2 Instance**).
2. PostgreSQL database is hosted on **AWS RDS**.
3. API documentation is currently hosted on **Swagger Hub** (https://app.swaggerhub.com/apis/karolispx/golang-rh-todo/1.0.0-oas3/) and can be used for making all API calls that are available in the application. More information contained in links below. 
4. **docs** folder in the repository also contains HTML generated Swagger docs and Postman collection in JSON format (for AWS EC2 Instance calls).

***

## Features & Endpoints:
* [User Authentication](https://github.com/karolispx/golang-rh-todo/wiki/Feature:-User-Authentication)
* [Tasks](https://github.com/karolispx/golang-rh-todo/wiki/Feature:-Tasks)

***

## Installation:
* [Local Environment](https://github.com/karolispx/golang-rh-todo/wiki/Installation:-Local-Environment)
* [Openshift](https://github.com/karolispx/golang-rh-todo/wiki/Installation:-Openshift)
* [EC2 Instance](https://github.com/karolispx/golang-rh-todo/wiki/Installation:-EC2-Instance)

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
