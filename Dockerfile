FROM golang:1.21.4-alpine3.17

WORKDIR /code

# Copy the local package files to the container's workspace
COPY . .

# Download and install any required dependencies
RUN ["go", "install", "./cmd/gochatter"]

# Command to run the executable
CMD ["gochatter"]
