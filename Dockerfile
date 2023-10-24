# Use an official Go runtime as a parent image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /home/go/aplikacija

# Copy the Go application source code into the container
COPY . /home/go/aplikacija

# Build the Go application
RUN go build -o myapp

# Expose the port that the Go application will listen on
EXPOSE 9090

# Define the command to run your application
CMD ["./myapp"]
