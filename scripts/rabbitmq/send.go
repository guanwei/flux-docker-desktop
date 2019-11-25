package main

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

// 只能在安装 rabbitmq 的服务器上操作
func main() {
	cfg := new(tls.Config)
	cfg.RootCAs = x509.NewCertPool()

	if ca, err := ioutil.ReadFile("ssl/cacert.pem"); err == nil {
		cfg.RootCAs.AppendCertsFromPEM(ca)
	} else {
		failOnError(err, "Failed to append ca")
	}

	if cert, err := tls.LoadX509KeyPair("ssl/client_cert.pem", "ssl/client_key.pem"); err == nil {
		cfg.Certificates = append(cfg.Certificates, cert)
	} else {
		failOnError(err, "Failed to append client cert")
	}

	conn, err := amqp.DialTLS("amqps://admin:admin@localhost:5671/", cfg)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	body := "Hello World!"
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	log.Printf(" [x] Sent %s", body)
	failOnError(err, "Failed to publish a message")
}
