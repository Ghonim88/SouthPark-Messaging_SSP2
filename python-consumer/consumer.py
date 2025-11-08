#! /usr/bin/env python3
"""
South Park Message Consumer
Listens to RabbitMQ queue and prints messages from South Park characters.
"""
import pika
from config import Config
import json
import sys
import time
import signal
from datetime import datetime
from pika.exceptions import AMQPConnectionError, ChannelClosedByBroker

class SouthParkConsumer:
    """Consumer for South Park messages from RabbitMQ."""

    def __init__(self):
        self.connection = None
        self.channel = None
        self.should_stop = False

        #setup signal handlers for graceful shutdown
        signal.signal(signal.SIGINT, self.handle_signal)
        signal.signal(signal.SIGTERM, self.handle_signal)
        
    
    def handle_signal(self, signum, frame):
        """Handle termination signals for graceful shutdown."""
        print("Received termination signal. Shutting down...")
        self.should_stop = True
        self.close_connection()
        sys.exit(0)

    def connect(self):
        """Establish connection to RabbitMQ."""
        max_retries = 10
        retry_delay = 5  # seconds
        for attempt in range(1,max_retries +1):
            try:
                print(f"Attempting to connect to RabbitMQ (Attempt {attempt}/{max_retries})...")
                creadentials = pika.PlainCredentials(Config.RABBITMQ_USER, Config.RABBITMQ_PASSWORD)
                parameters = pika.ConnectionParameters(
                    host=Config.RABBITMQ_HOST,
                    port=Config.RABBITMQ_PORT,
                    virtual_host=Config.RABBITMQ_VHOST,
                    credentials=creadentials,
                    heartbeat=600, # heartbeat interval for connection health
                    blocked_connection_timeout=300, # timeout for blocked connections
                )
                #Establish connection
                self.connection = pika.BlockingConnection(parameters)
                self.channel = self.connection.channel()

                #Declare the queue (idempotent operation)
                self.channel.queue_declare(queue=Config.QUEUE_NAME, durable=True)
                print("Connected to RabbitMQ successfully.")
                print(f"Listening to queue: {Config.QUEUE_NAME}")
                print("-"*50)
                return True
            except AMQPConnectionError as e:
                print(f"Connection attempt {attempt} failed: {e}")
                if attempt < max_retries:
                    print(f"Retrying in {retry_delay} seconds...")
                    time.sleep(retry_delay)
                else:
                    print("Max retries reached. Exiting.")
                    return False
            except Exception as e:
                print(f"Unexpected error during connection: {e}")
                return False
            
        return False
    
    def callback(self, ch, method, properties, body):
        """Callback function to process received messages."""
        try:
            # Decode the message
            message_str = body.decode('utf-8')
            message = json.loads(message_str)

            # Extract message details (producer uses 'body' field)
            author = message.get('author', 'Unknown')
            body = message.get('body', '')
            sent_at = message.get('sent_at', '')

            # Format timestamp
            if sent_at:
                try:
                    dt = datetime.fromisoformat(sent_at.replace('Z', '+00:00')) 
                    formatted_time = dt.strftime('%Y-%m-%d %H:%M:%S')
                except:
                    formatted_time = sent_at
            else:
                formatted_time = 'Unknown Time'

            # Print the message
            print("\n"+"="*50)
            print(f"Author: {author}")
            print(f"Body: {body}")
            print(f"Sent At: {formatted_time}")
            print("="*50)

            # Acknowledge message
            ch.basic_ack(delivery_tag=method.delivery_tag)
        except json.JSONDecodeError as e:
            print(f"Failed to parse Json: {e}")
            print(f"Raw message: {body}")
            #Reject the message without requeuing
            ch.basic_nack(delivery_tag=method.delivery_tag, requeue=False)
        except Exception as e:
            print(f"Error processing message: {e}")
            #Reject the message without requeuing
            ch.basic_ack(delivery_tag=method.delivery_tag, requeue=False)

    def start_consuming(self):
        """Start consuming messages from the queue."""
        try:
            #Set up consumer
            self.channel.basic_qos(prefetch_count=1)  # Fair dispatch
            self.channel.basic_consume(queue=Config.QUEUE_NAME, on_message_callback=self.callback, auto_ack=False)

            print("Consumer is ready! Waiting for messages...")
            print(" (Press Ctrl+C to exit)")
            print()

            #Start consuming(blocking call)
            # Note: start_consuming is a channel method for BlockingConnection
            self.channel.start_consuming()
        except KeyboardInterrupt:
            print("\n Consumer interrupted by user.")
            self.close_connection()
        except ChannelClosedByBroker as e:
            print(f"Channel closed by broker: {e}")
            self.close_connection()
        except Exception as e:
            print(f"Unexpected error during consumption: {e}")
            self.close_connection()
    
    def close_connection(self):
        """Close the RabbitMQ connection gracefully."""
        try:
            if self.channel and self.channel.is_open:
                self.channel.stop_consuming()
                self.channel.close()
                print("Channel closed.")

            if self.connection and self.connection.is_open:
                self.connection.close()
            print("Connection to RabbitMQ closed.")
        except Exception as e:
            print(f"Error closing connection: {e}")

    def run(self):
        """Run the consumer."""
        if self.connect():
            self.start_consuming()
        else:
            print("Failed to connect to RabbitMQ. Exiting.")
            sys.exit(1)


def main():
    """Main entry point for the consumer."""
    Config.print_config()
    consumer = SouthParkConsumer()
    consumer.run()

if __name__ == "__main__":
    main()
