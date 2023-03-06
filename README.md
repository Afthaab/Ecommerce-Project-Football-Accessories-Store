# Ecommerce-Project-Football-Accessories-Store
This project is written purely in Go Language. Gin (Http Web Frame Work) is used in this project. PostgreSQL Database is used to manage the data.
## Framework Used
Gin-Gonic: This whole project is built on Gin frame work. Its is a popular http web frame work. 
```
go get -u github.com/gin-gonic/gin
```
## Database used:
PostgreSQL: PostgreSQL is a powerful, open source object-relational database. The data managment in this project is done using PostgreSQL. ORM tool named GORM is also been used to simplify the forms of queries for better understanding.

```
go get -u gorm.io/gorm
go get -u gorm.io/driver/postgres
```
## External Packages Used
#### Razorpay
For Payment I have used the test case of Razorpay.
```
github.com/razorpay/razorpay-go
```
#### Validator
Package validator implements value validations for structs and individual fields based on tags.
```
github.com/go-playground/validator/v10
```
#### Twilio
The twilio-go helper library lets you write Go code to make HTTP requests to the Twilio API and get the OTP. This is open source library.
```
github.com/twilio/twilio-go/rest/api/v2010
```
#### Gomail
Gomail is a simple and efficient package to send emails. It is well tested and documented.
```
gopkg.in/mail.v2
```
#### JWT 
JSON Web Tokens are an open, industry standard RFC 7519 method for representing claims securely between two parties.
```
github.com/golang-jwt/jwt/v4
```
#### Commands to run project:
```
go run main.go
```

## API Platform Used
API platforms Postman is used to run all the API's Provided by this project

#### API Documentation
```
https://documenter.getpostman.com/view/25380689/2s93CUJVxP
```
#### Database Design
Databse is designed using draw SQL website
```
https://footballestoreimagebucket.s3.ap-south-1.amazonaws.com/drawSQL-e-commerce-db-export-2023-03-06.png
```
