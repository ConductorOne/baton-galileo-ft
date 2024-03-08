FROM gcr.io/distroless/static-debian11:nonroot
ENTRYPOINT ["/baton-galileo-ft"]
COPY baton-galileo-ft /