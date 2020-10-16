FROM ghcr.io/bryk-io/shell:0.1.0

# Expose required ports
EXPOSE 8080

# Expose required volumes
VOLUME /etc/dlt4eu

# Add application binary and use it as default entrypoint
COPY dlt4eu /bin/dlt4eu
ENTRYPOINT ["/bin/dlt4eu"]
