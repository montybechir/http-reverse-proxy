# tests/Dockerfile  

FROM alpine:latest  

# Install necessary packages like bash, curl, and netcat  
RUN apk add --no-cache bash curl netcat-openbsd  

# Set the working directory  
WORKDIR /app  

# Copy the tests and scripts into the container  
COPY tests /app/tests  
COPY scripts/wait-for.sh /usr/local/bin/wait-for.sh  
COPY scripts/run-tests.sh /usr/local/bin/run-tests.sh  

# Make the test scripts executable  
RUN chmod +x /usr/local/bin/wait-for.sh /usr/local/bin/run-tests.sh  

# Define the entrypoint  
ENTRYPOINT ["/usr/local/bin/run-tests.sh"]  
