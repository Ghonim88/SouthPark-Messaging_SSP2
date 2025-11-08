FROM python:3.11-slim

# Set working directory
WORKDIR /app
# Install dependencies
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy the application files
COPY config.py .
COPY consumer.py .

#Make consumer.py executable
RUN chmod +x consumer.py

# Run the consumer
CMD ["python","-u","consumer.py"]