FROM alpine
ENV HOME=/
COPY /dist/scan_compare-linux-amd64 /bin/scan_compare
ENTRYPOINT ["/bin/scan_compare"]