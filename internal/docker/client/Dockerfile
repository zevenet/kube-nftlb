# Download the latest Debian image
FROM debian:stable

# Default shell when executing RUN
SHELL ["/bin/bash", "-c"]

# With --build-arg=KEY, you can change the default key "12345"
ARG KEY=12345
ENV KEY_NFTLB=${KEY}

# After compiling the Golang client, copy the binary to / and make it the entrypoint
COPY ./app /app
CMD [ "/bin/bash", "-c", "./app $KEY_NFTLB" ]
