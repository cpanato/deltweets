FROM iron/base
LABEL maintainer "carlos.panato <ctadeu@gmail.com>"
LABEL version="0.2"

WORKDIR /app

# copy binary into image
COPY deltweets /app

ENTRYPOINT ["./deltweets"]
