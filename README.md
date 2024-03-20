# HealthBite

Is the name of the site three friends started working on where we are determined to revolutionize the approach to public health. Recognizing the importance of preventive measures, we initiated the implementation of a groundbreaking web app, aimed at promoting healthier choices among the population.
This project is made for the FIICode Web & Dev competition. You can find more info [here](https://fiicode-api.asii.ro/static/webdev2024)

## Key aspects

#### Backend Development in Go
The site's backend will be written in Go, a programming language known for its speed and efficiency. Go compiles directly to machine code, resulting in a fast backend, and it also supports multithreading, allowing for concurrent execution of tasks.

#### Input Sanitization & Regex Checking
To ensure data security and integrity, the site will implement input sanitization and regular expression (regex) checking. This will help prevent malicious input and enforce specific patterns for user-provided data, enhancing the overall security and reliability of the system.

#### Minimal JavaScript Usage
The site will prioritize minimal use of JavaScript to promote fast loading times and reduce client-side processing. By keeping JavaScript usage to a minimum, the site will deliver a streamlined user experience while maintaining high performance and efficient resource utilization.

#### Secure data storage
Users and products are referenced by their UUIDs, and sensitive information is encrypted with bcrypt. This makes the project comply with modern security practices.
#### Continous Integration
After each commit, the source code goes through a CI/CD pipeline to ensure that every change in the code does not break the project. This ensures the project will run every time and without any hiccups in development!

#### Powered by Amazon Web Services
Our project utilizes AWS to power the project. With the power of the cloud, we are able to compute high amounts of users and products without any downtime!

## Prerequisites

To run the application smoothly, make sure to configure a PostgreSQL database with the constants set in the constants.go file.

```
package main

const (
        // Postgresql
        host     = "<PG_IP>"
        port     = <DB_PORT>
        user     = "<PG_USERNAME>"
        password = "<PASSWORD>"
        dbname   = "<DATABASE>"

        // AWS
        aws_region = "<AWS_REGION>"
        aws_access = "<AWS_ACCESS_KEY>"
        aws_secret = "<AWS_SECRET>"

        // S3
        bucketName = "<BUCKET_NAME>"
)
```

## Deployment

```
git clone https://github.com/wr47hcnf/HealthBite.git
cd HealthBite
go run ./src
```

## Built With

* [HTML](https://www.w3schools.com/html/) - basic framework for the content on the webpage
* [Bootstrap](http://www.dropwizard.io/1.0.2/docs/) - front-end framework
* [CSS](https://www.w3schools.com/css/) - control the look and feel of web pages
* [JavaScript](https://www.w3schools.com/js/default.asp) - enhancing interactivity on websites and web applications.
* [GO](https://www.w3schools.com/go/index.php) - provide a simple and efficient way to create reliable and high-performance software

## Authors

* **Patrick Covaci** - *Back end (computing)* - [Ty3r0X](https://github.com/Ty3r0X)
* **Ivența Rareș** - *Front-end development* - [iv3bta](https://github.com/iv3bta)
* **Chiriță Rareș** - *Dynamic front end functionality* - [Raresbomba](https://github.com/Raresbomba)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details
