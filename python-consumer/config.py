import os
from dotenv import load_dotenv

# Load environment variables from .env file if it exists
load_dotenv()

class Config:
    """Configuration for the python consumer."""

    # ReabbitMQ settings
    RABBITMQ_HOST = os.getenv('RABBITMQ_HOST', 'localhost')
    RABBITMQ_PORT = int(os.getenv('RABBITMQ_PORT', 5672))
    RABBITMQ_USER = os.getenv('RABBITMQ_USER', 'guest')
    RABBITMQ_PASSWORD = os.getenv('RABBITMQ_PASSWORD', 'guest')
    RABBITMQ_VHOST = os.getenv('RABBITMQ_VHOST', '/')

    # Queue name (must match the Go API queue name)
    QUEUE_NAME = 'southpark_messages'

    @classmethod
    def get_rabbitmq_url(cls):
        """Construct the RabbitMQ connection URL."""
        return f'amqp://{cls.RABBITMQ_USER}:{cls.RABBITMQ_PASSWORD}@{cls.RABBITMQ_HOST}:{cls.RABBITMQ_PORT}{cls.RABBITMQ_VHOST}'
    
    @classmethod
    def print_config(cls):
        """Print the current configuration (for debugging purposes)."""
        print("="*50)
        print("Python Consumer Configuration:")
        print(f"RabbitMQ Host: {cls.RABBITMQ_HOST}")
        print(f"RabbitMQ Port: {cls.RABBITMQ_PORT}")
        print(f"RabbitMQ User: {cls.RABBITMQ_USER}")
        print(f"Queue Name: {cls.QUEUE_NAME}")
        print("="*50)
