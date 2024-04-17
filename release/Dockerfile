FROM --platform=$TARGETPLATFORM registry.access.redhat.com/ubi9-minimal

WORKDIR /app

ARG project_name
ARG TARGETPLATFORM
ARG version
ARG tag
ENV APP_BINARY=${project_name}
ENV VERSION=${version:-develop}
ENV TAG=${tag:-main}
COPY ./release/entrypoint.sh .
RUN chmod +x entrypoint.sh

COPY ./dist/$TARGETPLATFORM/$APP_BINARY .

ENTRYPOINT ["./entrypoint.sh"]
