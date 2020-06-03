# Download the latest Debian image
FROM debian:stable

# Default shell when executing RUN
SHELL ["/bin/bash", "-c"]


RUN cat /etc/apt/sources.list

RUN apt-get update
RUN apt-get install -y gnupg ca-certificates wget libterm-readline-perl-perl supervisor


# Install nftlb and dependencies
RUN wget -O - http://repo.zevenet.com/zevenet.com.gpg.key | apt-key add -
RUN echo "deb [arch=amd64] http://repo.zevenet.com/ce/v5 buster main" | tee -a /etc/apt/sources.list

RUN apt-get update
RUN apt-get install -y nftlb


# Replacing nftlb with a devel version
COPY ./nftlb.deb /nftlb.deb
RUN if [ -s "nftlb.deb" ] ; then dpkg -i ./nftlb.deb ; fi


# With --build-arg=KEY, you can change the default key "12345"
ARG KEY=12345
ENV KEY_NFTLB=${KEY}

# After compiling the Golang client, copy the binary to / and make it the entrypoint
COPY ./app /goclient

COPY ./start.sh /start.sh
RUN sed -Ei "s/#KEY/$KEY_NFTLB/g" /start.sh
CMD [ "/bin/bash", "-c", "/start.sh" ]
