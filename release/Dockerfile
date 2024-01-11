FROM --platform=$TARGETPLATFORM registry.access.redhat.com/ubi9-minimal

# Set working directory
WORKDIR /app

ARG project_name
ARG TARGETPLATFORM
COPY ./dist/$TARGETPLATFORM .

ARG version
ARG tag
ENV APP_BINARY=${project_name}
ENV VERSION=${version:-develop}
ENV TAG=${tag:-master}

# Set entrypoint
ENTRYPOINT ["./$APP_BINARY"]
