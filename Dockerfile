# Build stage for Go backend
FROM golang:1.24-alpine AS go-builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY cmd/ ./cmd/
COPY internal/ ./internal/

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/main.go

# Build stage for Frontend
FROM node:18-alpine AS frontend-builder

WORKDIR /app

# Install pnpm
RUN npm install -g pnpm

# Copy frontend files
COPY frontend/package.json frontend/pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile

# Copy frontend source and build
COPY frontend/ ./
# Set environment variable for build
# Use relative URL for Railway (proxied through Next.js)
ARG NEXT_PUBLIC_DEV_SERVER_URL=/api
ENV NEXT_PUBLIC_DEV_SERVER_URL=$NEXT_PUBLIC_DEV_SERVER_URL
RUN pnpm build

# Final stage - single container with both services
FROM node:18-alpine

WORKDIR /app

# Install pnpm for running Next.js
RUN npm install -g pnpm

# Copy Go binary
COPY --from=go-builder /app/main ./

# Copy built frontend
COPY --from=frontend-builder /app/.next ./frontend/.next
COPY --from=frontend-builder /app/public ./frontend/public
COPY --from=frontend-builder /app/package.json ./frontend/
COPY --from=frontend-builder /app/node_modules ./frontend/node_modules

# Create startup script
RUN echo '#!/bin/sh' > /app/start.sh && \
    echo 'echo "Starting Go backend..."' >> /app/start.sh && \
    echo './main &' >> /app/start.sh && \
    echo 'echo "Starting Next.js frontend..."' >> /app/start.sh && \
    echo 'cd frontend && PORT=3000 pnpm start &' >> /app/start.sh && \
    echo 'wait' >> /app/start.sh

RUN chmod +x /app/start.sh

EXPOSE 3000

# Start both services
CMD ["/app/start.sh"]