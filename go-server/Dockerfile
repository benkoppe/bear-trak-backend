# syntax=docker/dockerfile:1

################################################################################
# Create a stage for building the application.
ARG GO_VERSION=1.24.2
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION} AS build
WORKDIR /src

# Download dependencies as a separate step to take advantage of Docker's caching.
# Leverage a cache mount to /go/pkg/mod/ to speed up subsequent builds.
# Leverage bind mounts to go.sum and go.mod to avoid having to copy them into
# the container.
RUN --mount=type=cache,target=/go/pkg/mod/ \
  --mount=type=bind,source=go.sum,target=go.sum \
  --mount=type=bind,source=go.mod,target=go.mod \
  go mod download -x

# This is the architecture you're building for, which is passed in by the builder.
# Placing it here allows the previous steps to be cached across architectures.
ARG TARGETARCH

# Build the application.
# Leverage a cache mount to /go/pkg/mod/ to speed up subsequent builds.
# Leverage a bind mount to the current directory to avoid having to copy the
# source code into the container.
RUN --mount=type=cache,target=/go/pkg/mod/ \
  --mount=type=bind,target=. \
  CGO_ENABLED=0 GOARCH=$TARGETARCH go build -o /bin/server .

################################################################################
# Create a stage for running the application that contains the minimal
# runtime dependencies for the application. 
FROM chromedp/headless-shell:latest AS final

# Install any runtime dependencies.
# Leverage a cache mount to /var/cache/apk/ to speed up subsequent builds.
RUN --mount=type=cache,target=/var/cache/apt \
  apt-get update && \
  apt-get install -y --no-install-recommends \
  ca-certificates \
  tzdata \
  && rm -rf /var/lib/apt/lists/* \
  && update-ca-certificates


# Create a non-privileged user that the server will run under.
# See https://docs.docker.com/go/dockerfile-user-best-practices/
# ARG UID=10001
# RUN useradd -u ${UID} -M -s /usr/sbin/nologin appuser
# USER appuser

# Copy the executable from the "build" stage.
COPY --from=build /bin/server /bin/

# Expose the port that the application listens on.
EXPOSE 3000

# What the container should run when it is started.
ENTRYPOINT [ "/bin/server" ]
