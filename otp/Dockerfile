FROM opentripplanner/opentripplanner:latest

# OTP uses /var/opentripplanner as its working directory
WORKDIR /var/opentripplanner

# Copy static files into the container
COPY ithaca.pbf /var/opentripplanner/ithaca.pbf
COPY router-config.json /var/opentripplanner/router-config.json

# Download the GTFS zip
RUN curl -L -o /var/opentripplanner/tcat.gtfs.zip "https://realtimetcatbus.availtec.com/InfoPoint/GTFS-zip.ashx"

# Expose the OTP server port
EXPOSE 8080

# Set java options
ENV JAVA_TOOL_OPTIONS="-Xmx2g"

# Build the OTP graph
RUN /docker-entrypoint.sh --build --save

# Start OTP
CMD [ "--load", "--serve"]
