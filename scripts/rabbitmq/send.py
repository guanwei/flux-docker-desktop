import logging
import pika
import ssl

logging.basicConfig(level=logging.INFO)
context = ssl.create_default_context(
    cafile="ssl/cacert.pem")
context.load_cert_chain("ssl/client.cert.pem",
                        "ssl/client.key.pem")
ssl_options = pika.SSLOptions(context, "localhost")
conn_params = pika.ConnectionParameters(port=5671,
                                        ssl_options=ssl_options)

with pika.BlockingConnection(conn_params) as conn:
    ch = conn.channel()
    ch.queue_declare("foobar")
    ch.basic_publish("", "foobar", "Hello, world!")
