#!/bin/bash

# Build the Docker image
docker build -t dream-forum .

# Run the Docker container
docker run --rm -p 4000:4000 dream-forum
