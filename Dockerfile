# Use an official Go runtime as a parent image
FROM golang:latest

# Set the working directory to /app
WORKDIR /app

# Copy the current directory contents into the container at /app
COPY . /app

# Download and install sqlite3
RUN go get -u github.com/mattn/go-sqlite3

# Build the Go application
RUN go build -o main .

# Expose port to the outside world
EXPOSE 4000

# Command to run the executable
CMD ["./main"]
