This is a special AVL module that loads data from the Trak api.

To get this working, I had to follow these steps:

- Download `Core.jar` with:

```bash
curl -s https://api.github.com/repos/TheTransitClock/transitime/releases/latest | jq -r ".assets[].browser_download_url" | grep 'Core.jar' | xargs -L1 wget
```

- Install file to `mvn` with:

```bash
mvn install:install-file -Dfile=./Core.jar -DgroupId=org.transitclock -DartifactId=transitclock-core -Dversion=2.2 -Dpackaging=jar
```

- That should be everything! See `pom.xml` to see how this is referenced locally as a dependency. Build with `mvn package`.
