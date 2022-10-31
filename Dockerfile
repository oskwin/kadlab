FROM golang:1.17-alpine

# Add the commands needed to put your compiled go binary in the container and
# run it when the container starts.
#
# See https://docs.docker.com/engine/reference/builder/ for a reference of all
# the commands you can use in this file.
#
# change to correct directory
 WORKDIR /Users/oskwin/Documents/courses/D7024E/lab/mykadlab/
 # copy all files to image
 COPY   /d7024e /usr/local/go/src/kadlab/d7024e
 COPY main.go /usr/local/go/src/kadlab
 
 COPY go.mod ./
 RUN go mod download

 COPY *.go ./

 RUN go build -o /kadlab

 EXPOSE 4000

 CMD [ "/kadlab" ]

# In order to use this file together with the docker-compose.yml file in the
# same directory, you need to ensure the image you build gets the name
# "kadlab", which you do by using the following command:
#
# $ docker build . -t kadlab
