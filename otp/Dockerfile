FROM opentripplanner/opentripplanner:latest

# OTP uses /var/opentripplanner as its working directory
WORKDIR /var/opentripplanner

# Set Java memory limit (can still be overridden at runtime)
ENV JAVA_TOOL_OPTIONS="-Xmx500m"

# Download the GTFS zip
RUN --mount=type=bind,source=./data,target=/var/opentripplanner \
  curl -L -o /var/opentripplanner/gtfs.zip "$(cat /var/opentripplanner/gtfs-url.txt)"

# Build the OTP graph
RUN /docker-entrypoint.sh --build --save

# Expose the OTP server port
EXPOSE 8080

# Start OTP
CMD [ "--load", "--serve"]
