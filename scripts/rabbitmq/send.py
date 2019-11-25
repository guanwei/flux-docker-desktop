import logging
import pika
import ssl

#logging.basicConfig(level=logging.INFO)

context = ssl.create_default_context(
    cafile="ssl/ca_certificate.pem")
context.load_cert_chain("ssl/client_certificate.pem",
                        "ssl/client_key.pem")
ssl_options = pika.SSLOptions(context, "localhost")
credentials = pika.PlainCredentials("admin", "admin")
parameters = pika.ConnectionParameters(host="localhost",
                                       port=5671,
                                       virtual_host="/",
                                       ssl_options=ssl_options,
                                       credentials=credentials)

body = "Hello World!"

with pika.BlockingConnection(parameters) as conn:
    ch = conn.channel()
    ch.queue_declare("hello")
    ch.basic_publish("", "hello", body)
    print(" [x] Sent {}".format(body))
