# Download the latest Debian image
FROM debian:stable

# Default shell when executing RUN
SHELL ["/bin/bash", "-c"]

RUN apt-get update
RUN apt-get install -y gnupg ca-certificates wget

# Install nftlb and dependencies
RUN wget -O - http://repo.zevenet.com/zevenet.com.gpg.key | apt-key add -
RUN echo "deb [arch=amd64] http://repo.zevenet.com/ce/v5 buster main" | tee -a /etc/apt/sources.list

RUN apt-get update
RUN apt-get install -y nftlb

# With --build-arg=KEY, you can change the default key "12345"
ARG KEY=12345
ENV KEY_NFTLB=${KEY}

# After compiling everything, make nftlb as the entrypoint
CMD [ "/bin/bash", "-c", "/usr/local/zevenet/app/nftlb/sbin/nftlb -k $KEY_NFTLB" ]
