# Simple TODO app built in GoLang using PostgreSQL for data storage.

---

## About This Application:

1. GoLang application is hosted on **Openshift** and on **AWS EC2 Instance** (as a backup).
2. PostgreSQL database is hosted on **AWS RDS**.
3. API documentation is currently hosted on **Swagger Hub** (https://app.swaggerhub.com/apis/karolispx/golang-rh-todo/1.0.0-oas3/) and can be used for making all API calls that are available in the application. More information contained in links below.
4. VueJS application has been built to be used with this application. VueJS application is currently hosted on AWS S3 bucket and can be viewed here: http://vue-rh-todo.s3-website-eu-west-1.amazonaws.com
    - Codebase can be found here: https://github.com/karolispx/vue-rh-todo
    - **This VueJS application has been put together very quickly, it's nowhere near being good enough to be used in production. **

---

## Application's Structure:

- **api** - files inside this folder contain logic for RestAPI calls.
- **bin** - folder for binary/openshift files.
- **helpers** - files inside this folder contain helper functions that are reused across the application.
- **models** - files inside this folder contain functions that are used to talk to database.
- **config.sample.json** - used as a template for storing application's config variables when working locally or on AWS EC2 instance.
- **Dockerfile** - used to build docker image for Openshift.
- **main.go** - main package of the application, setups routes.

---

## Features & Endpoints:

- [User Authentication](https://github.com/karolispx/golang-rh-todo/wiki/Feature:-User-Authentication)
- [Tasks](https://github.com/karolispx/golang-rh-todo/wiki/Feature:-Tasks)

---

## Installation:

- [Local Environment](https://github.com/karolispx/golang-rh-todo/wiki/Installation:-Local-Environment)
- [Openshift](https://github.com/karolispx/golang-rh-todo/wiki/Installation:-Openshift) - Including **Jenkins** pipeline. 
- [EC2 Instance](https://github.com/karolispx/golang-rh-todo/wiki/Installation:-EC2-Instance)

---

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

---

## Future Plans

1. Automate build and deployment to both **Openshift** and **AWS EC2 Instance** by possibly using a combination of **bash scripts** and **Jenkins** or **GOCD** - Jenkins has been used to automate the deployment to Openshift.
2. Automate database table creation.
3. Error logging or some sort of monitoring to get notifications if something goes wrong.
4. JWT validation - what to do with previously user generated tokens that haven't expired yet but are valid?
5. User email validation/password reset - possibly use **Mailgun** or **Mandrill** or **Sparkpost** for sending emails.
