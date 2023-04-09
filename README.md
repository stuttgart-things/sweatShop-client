# stuttgart-things/yacht-application-client

YAC is a gRPC Client for importing & sending revisionRuns to YAS

## GET REVISIONRUN FROM GIT REPOSITROY

```
yacht-application-client get \
--repo https://github.com/stuttgart-things/yacht-application-server.git
```

## SEND REVISIONRUN(S) TO YAS

```
yacht-application-client send \
--endpoint yas.dev.sthings.tiab.ssc.sva.de \
--file yacht.json 
```

License
-------


Author Information
------------------

Patrick Hermann, stuttgart-things 04/2023
