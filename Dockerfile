# Dockerfile definition for Backend application service.

# From which image we want to build. This is basically our environment.
FROM golang:1.21-alpine AS Build

WORKDIR /app

# This will copy all the files in our repo to the inside the container at root location.
COPY . .

# Generate the OpenAPI server package expected by the application build.
RUN mkdir -p generated \
	&& go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest \
	&& /go/bin/oapi-codegen --package generated \
		-generate types,server,spec api.yml > generated/api.gen.go

# Build our binary at root location.
RUN GOPATH= go build -o /main cmd/main.go

####################################################################
# This is the actual image that we will be using in production.
FROM alpine:latest

# We need to copy the binary from the build image to the production image.
COPY --from=Build /main .

# This is the port that our application will be listening on.
EXPOSE 1323

# This is the command that will be executed when the container is started.
ENTRYPOINT ["./main"]