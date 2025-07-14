# Backend Dockerfile for Go API with Gophercloud and Terraform
FROM golang:1.22-alpine

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

# Install Air for hot reload (use specific version compatible with Go 1.22)
RUN go install github.com/cosmtrek/air@v1.49.0

# Set working directory
WORKDIR /app

# Create necessary directories
RUN mkdir -p /app/terraform/state /app/tmp

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

# Default command - will be overridden in docker-compose for development
CMD ["tail", "-f", "/dev/null"]
