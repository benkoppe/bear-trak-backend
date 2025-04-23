FROM opentripplanner/opentripplanner:latest

# OTP uses /var/opentripplanner as its working directory
WORKDIR /var/opentripplanner

# make school configurable at build time
ARG SCHOOL
# (required for cache-busting)
ENV SCHOOL=$SCHOOL

# Set Java memory limit (can still be overridden at runtime)
ARG JAVA_TOOL_OPTIONS
ENV JAVA_TOOL_OPTIONS="${JAVA_TOOL_OPTIONS:--Xmx512m}"

# Copy static files into the container
COPY ./data/${SCHOOL}/graph.pbf /var/opentripplanner/graph.pbf
COPY ./data/${SCHOOL}/router-config.json /var/opentripplanner/router-config.json
COPY ./data/${SCHOOL}/build-config.json /var/opentripplanner/build-config.json

# Build the OTP graph
RUN /docker-entrypoint.sh --build --save

# Expose the OTP server port
EXPOSE 8080

# Start OTP
CMD [ "--load", "--serve"]
