FROM golang:1.17-alpine


# Copy the local package files to the container’s workspace.
ADD . /go/src/github.com/William9923/httpmiddleware

# Build the golang-docker command inside the container.
RUN go install github.com/William9923/httpmiddleware

# Run the golang-docker command when the container starts.
ENTRYPOINT /go/bin/httpmiddleware

# get dependencies
RUN go mod vendor
RUN go tidy


#Expose port
EXPOSE 8080

# Run the application
CMD ["make", "run-demo"]


