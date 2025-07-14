# Backend Dockerfile for Go API with Gophercloud and Terraform
FROM golang:1.23-alpine

# Install system dependencies
RUN apk add --no-cache \
    curl \
    wget \
    unzip \
    git \
    bash \
    ca-certificates

# Install Terraform CLI
ARG TERRAFORM_VERSION=1.6.6
RUN wget https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip && \
    unzip terraform_${TERRAFORM_VERSION}_linux_amd64.zip && \
    mv terraform /usr/local/bin/ && \
    rm terraform_${TERRAFORM_VERSION}_linux_amd64.zip

# Verify Terraform installation
RUN terraform version

# Install Air for hot reload
RUN go install github.com/cosmtrek/air@v1.49.0

# Set working directory
WORKDIR /app

# Copy source code
COPY src/backend/ ./
COPY src/terraform/ ./

# Initialize module if needed
RUN if [ ! -f go.mod ]; then go mod init cloud-provider; fi

# Add all dependencies explicitly
# RUN go get github.com/gin-gonic/gin@v1.9.1
# RUN go get github.com/gophercloud/gophercloud@v1.14.1
# RUN go get github.com/spf13/viper@v1.17.0
# RUN go get github.com/gophercloud/gophercloud/openstack/identity/v3/projects
# RUN go get github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/keypairs
# RUN go get github.com/gophercloud/gophercloud/openstack/compute/v2/servers
# RUN go get github.com/gophercloud/gophercloud/openstack/compute/v2/images
# RUN go get github.com/gophercloud/gophercloud/openstack/compute/v2/flavors
# RUN go get github.com/gophercloud/gophercloud/openstack/networking/v2/networks
# RUN go get github.com/gophercloud/gophercloud/openstack/networking/v2/subnets
# RUN go get github.com/gophercloud/gophercloud/openstack/networking/v2/ports
# RUN go get github.com/gophercloud/gophercloud/openstack/blockstorage/v3/volumes

# Download all dependencies and rebuild go.sum
RUN go mod download
RUN go mod tidy

# Verify all dependencies are properly cached
RUN go mod verify

# Create necessary directories
RUN mkdir -p /app/terraform/state /app/tmp

# Test build to ensure everything works
RUN go build -o /tmp/test ./cmd/api

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

# Run with Air for hot reload
CMD ["air", "-c", ".air.toml"]