# Dream Forum authentication

Dream Forum is a web-based platform designed to facilitate user communication, allowing them to authenticate using Google or GitHub, create posts, comments, associate categories with posts, like/dislike content, and apply filters. Built with Go and SQLite, this project introduces users to web basics, session management, database manipulation, and Docker.

## Start the program:

### The faster way:

1. Install Docker Desktop app: https://www.docker.com/products/docker-desktop/

* If using Linux on WSL, intsall the Docker Desktop for Windows.

2. Verify that you have docker installed by running: ```docker --version```.

3. Open your Docker Desktop app.

4. ```bash run_forum.sh ``` or ```go run . ```

after that go to: http://localhost:4000

#### Dummy Users
You can use the following dummy users to log in for easier testing:

> **Email** siki@mail.ee | **Username** Siki | **Password** siki

> **Email** limpa@mail.com | **Username** Limpa | **Password** limpa

> **Email** papakoi@mail.com | **Username** Papakoi | **Password** papakoi


### The other way:
1. Install Docker Desktop app: https://www.docker.com/products/docker-desktop/

* If using Linux on WSL, intsall the Docker Desktop for Windows.

2. Verify that you have docker installed by running: ```docker --version```.

3. Open your Docker Desktop app.

4. Run: ```docker build -t dream-forum .```.

* Wait for Docker to build the image. It can take a few minutes.

5. Verify that the image has been built by running: ```docker images```.

* You should see the image named ```dream-forum```.

6. Run the image: ```docker run --rm -p 4000:4000 dream-forum```.

* This will delete the container after it has been stopped.

7. Go to: http://localhost:8080 in your browser and start posting, commenting, and liking. :)

## Remove image
When you are finished, run: ```sudo docker system prune -a```

* This will remove all the Docker images and containers that may have been left.

## Authors
> **purr-purr** 

> **mvolkons**

