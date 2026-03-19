FROM envoyproxy/envoy:v1.30-latest

COPY files/envoy.compose.yaml /etc/envoy/envoy.yaml

CMD ["-c", "/etc/envoy/envoy.yaml", "--log-level", "info"]
