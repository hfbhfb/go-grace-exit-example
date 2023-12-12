FROM alpine
# RUN apk add fuse3 
ADD ./static/ /static/
COPY ./server /

ENTRYPOINT ["/server"]