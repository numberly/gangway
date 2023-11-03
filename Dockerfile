# Use a specific version of golang alpine image
FROM golang:1.21.3-alpine3.18 AS buildstage

# Install curl, which is more reliable than ADD for downloading files
RUN apk add --no-cache curl

WORKDIR /go/src/github.com/numberly/gangway

# Use curl to download assets
RUN mkdir assets \
    && curl -fSL -o assets/materialize.min.css https://raw.githubusercontent.com/Dogfalo/materialize/v1-dev/dist/css/materialize.min.css \
    && curl -fSL -o assets/materialize.min.js https://raw.githubusercontent.com/Dogfalo/materialize/v1-dev/dist/js/materialize.min.js \
    && curl -fSL -o assets/prism-core.min.js https://raw.githubusercontent.com/PrismJS/prism/v1.28.0/components/prism-core.min.js \
    && curl -fSL -o assets/prism-bash.min.js https://raw.githubusercontent.com/PrismJS/prism/v1.28.0/components/prism-bash.min.js \
    && curl -fSL -o assets/prism-yaml.min.js https://raw.githubusercontent.com/PrismJS/prism/v1.28.0/components/prism-yaml.min.js \
    && curl -fSL -o assets/prism-powershell.min.js https://raw.githubusercontent.com/PrismJS/prism/v1.28.0/components/prism-powershell.min.js \
    && curl -fSL -o assets/prism-tomorrow.min.css https://raw.githubusercontent.com/PrismJS/prism/v1.28.0/themes/prism-tomorrow.min.css

# Copy the local files to the container
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY assets/ assets/
COPY cmd/ cmd/
COPY internal/ internal/
COPY templates/ templates/

# Build the application
RUN cd cmd/gangway && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/bin/gangway

# Use distroless static image for minimal size and surface area
FROM gcr.io/distroless/static:nonroot

# Ensure non-root user has necessary permissions
COPY --from=buildstage --chown=nonroot:nonroot /go/bin/gangway /bin/gangway

# Use array syntax for ENTRYPOINT, which doesn't invoke a shell
# This requires the configuration file to be provided as an argument
ENTRYPOINT ["/bin/gangway"]
CMD ["-config", "/gangway/gangway.yaml"]