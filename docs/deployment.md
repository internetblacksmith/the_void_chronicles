# Deployment Guide

Complete guide for deploying the Void Reavers SSH Reader in production environments, from single server setups to scalable cloud deployments.

## 🎯 Deployment Overview

This guide covers:
- **Local Production**: Single server deployment
- **Cloud Deployment**: AWS, GCP, Azure setups
- **Container Deployment**: Docker and Kubernetes
- **High Availability**: Load balancing and clustering
- **Monitoring**: Logging, metrics, and alerts

## 🏠 Local Production Deployment

### Prerequisites

- Linux server (Ubuntu 20.04+ recommended)
- Sudo access for system configuration
- Basic networking knowledge
- Domain name (optional, for external access)

### Step 1: System Preparation

#### Update System
```bash
sudo apt update && sudo apt upgrade -y
```

#### Install Dependencies
```bash
# Ubuntu/Debian
sudo apt install -y golang-go openssh-client git curl wget

# CentOS/RHEL/Fedora
sudo dnf install -y golang openssh-clients git curl wget
```

#### Create Service User
```bash
sudo useradd --system --home-dir /opt/void-reader --create-home --shell /bin/false voidreader
```

### Step 2: Application Deployment

#### Using Deployment Script (Recommended)
```bash
# Clone/copy application code
git clone <repository> void-reavers-reader
cd void-reavers-reader

# Deploy with automated script
sudo ./deploy.sh
```

#### Manual Deployment
```bash
# Build application
./build.sh

# Create installation directory
sudo mkdir -p /opt/void-reader
sudo mkdir -p /opt/void-reader/.ssh
sudo mkdir -p /opt/void-reader/.void_reader_data

# Copy files
sudo cp void-reader /opt/void-reader/
sudo cp -r book1_void_reavers /opt/void-reader/
sudo cp .ssh/id_ed25519* /opt/void-reader/.ssh/

# Set permissions
sudo chown -R voidreader:voidreader /opt/void-reader
sudo chmod 755 /opt/void-reader/void-reader
sudo chmod 600 /opt/void-reader/.ssh/id_ed25519
```

### Step 3: System Service Configuration

#### Install Systemd Service
```bash
# Copy service file
sudo cp systemd/void-reader.service /etc/systemd/system/

# Reload systemd
sudo systemctl daemon-reload

# Enable and start service
sudo systemctl enable void-reader
sudo systemctl start void-reader
```

#### Service Management Commands
```bash
# Check status
sudo systemctl status void-reader

# View logs
sudo journalctl -u void-reader -f

# Restart service
sudo systemctl restart void-reader

# Stop service
sudo systemctl stop void-reader
```

### Step 4: Network Configuration

#### Firewall Setup
```bash
# Ubuntu (UFW)
sudo ufw allow 23234/tcp
sudo ufw enable

# CentOS/RHEL (Firewalld)
sudo firewall-cmd --permanent --add-port=23234/tcp
sudo firewall-cmd --reload

# Direct iptables
sudo iptables -A INPUT -p tcp --dport 23234 -j ACCEPT
sudo iptables-save > /etc/iptables/rules.v4
```

#### Reverse Proxy (Optional)

**Nginx TCP Proxy:**
```nginx
# /etc/nginx/nginx.conf
stream {
    upstream void_reader {
        server localhost:23234;
    }
    
    server {
        listen 22222;
        proxy_pass void_reader;
        proxy_timeout 1s;
        proxy_responses 1;
        proxy_bind $remote_addr transparent;
    }
}
```

**HAProxy Configuration:**
```haproxy
# /etc/haproxy/haproxy.cfg
listen void_reader
    bind *:22222
    mode tcp
    server void1 localhost:23234 check
```

### Step 5: SSL/TLS Setup (Advanced)

#### Stunnel Configuration
```bash
# Install stunnel
sudo apt install stunnel4

# Configure stunnel
cat > /etc/stunnel/void-reader.conf << EOF
[void-reader]
accept = 443
connect = localhost:23234
cert = /etc/ssl/certs/void-reader.pem
key = /etc/ssl/private/void-reader.key
EOF

# Start stunnel
sudo systemctl enable stunnel4
sudo systemctl start stunnel4
```

## ☁️ Cloud Deployment

### AWS Deployment

#### EC2 Instance Setup

**Launch Instance:**
```bash
# Use Ubuntu 20.04 LTS AMI
# Instance type: t3.micro (for testing) or t3.small (production)
# Security group: Allow SSH (22) and custom port (23234)
```

**User Data Script:**
```bash
#!/bin/bash
apt update && apt upgrade -y
apt install -y golang-go git

# Clone and deploy application
cd /opt
git clone <repository> void-reader
cd void-reader
./deploy.sh

# Configure security group
aws ec2 authorize-security-group-ingress \
    --group-id sg-xxxxxxxxx \
    --protocol tcp \
    --port 23234 \
    --cidr 0.0.0.0/0
```

#### ECS Deployment

**Task Definition:**
```json
{
  "family": "void-reader",
  "networkMode": "awsvpc",
  "requiresCompatibilities": ["FARGATE"],
  "cpu": "256",
  "memory": "512",
  "taskRoleArn": "arn:aws:iam::account:role/void-reader-task-role",
  "containerDefinitions": [
    {
      "name": "void-reader",
      "image": "void-reader:latest",
      "portMappings": [
        {
          "containerPort": 23234,
          "protocol": "tcp"
        }
      ],
      "logConfiguration": {
        "logDriver": "awslogs",
        "options": {
          "awslogs-group": "/ecs/void-reader",
          "awslogs-region": "us-east-1",
          "awslogs-stream-prefix": "ecs"
        }
      }
    }
  ]
}
```

#### Lambda Deployment (Event-driven)

For serverless deployment with Lambda + API Gateway:

```go
// lambda/main.go
package main

import (
    "context"
    "github.com/aws/aws-lambda-go/lambda"
    "github.com/aws/aws-lambda-go/events"
)

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    // Initialize SSH reader session
    // Handle WebSocket connection
    return events.APIGatewayProxyResponse{
        StatusCode: 200,
        Body: "SSH Reader session started",
    }, nil
}

func main() {
    lambda.Start(handleRequest)
}
```

### Google Cloud Platform

#### Compute Engine
```bash
# Create instance
gcloud compute instances create void-reader \
    --image-family=ubuntu-2004-lts \
    --image-project=ubuntu-os-cloud \
    --machine-type=e2-micro \
    --tags=void-reader

# Configure firewall
gcloud compute firewall-rules create allow-void-reader \
    --allow tcp:23234 \
    --source-ranges 0.0.0.0/0 \
    --target-tags void-reader
```

#### Cloud Run
```yaml
# cloudrun.yaml
apiVersion: serving.knative.dev/v1
kind: Service
metadata:
  name: void-reader
spec:
  template:
    metadata:
      annotations:
        autoscaling.knative.dev/maxScale: "10"
    spec:
      containers:
      - image: gcr.io/PROJECT_ID/void-reader
        ports:
        - containerPort: 23234
        resources:
          limits:
            cpu: 1000m
            memory: 512Mi
```

### Azure Deployment

#### Container Instances
```bash
# Create resource group
az group create --name void-reader-rg --location eastus

# Deploy container
az container create \
    --resource-group void-reader-rg \
    --name void-reader \
    --image void-reader:latest \
    --ports 23234 \
    --ip-address public \
    --cpu 1 \
    --memory 1
```

## 🐳 Container Deployment

### Docker Deployment

#### Build and Run
```bash
# Build image
docker build -t void-reader .

# Run container
docker run -d \
    --name void-reader \
    -p 23234:23234 \
    -v $(pwd)/book1_void_reavers:/app/book1_void_reavers:ro \
    -v void_reader_data:/app/.void_reader_data \
    void-reader
```

#### Docker Compose
```yaml
# docker-compose.yml
version: '3.8'

services:
  void-reader:
    build: .
    ports:
      - "23234:23234"
    volumes:
      - ./book1_void_reavers:/app/book1_void_reavers:ro
      - void_reader_data:/app/.void_reader_data
      - ssh_keys:/app/.ssh
    restart: unless-stopped
    environment:
      - VOID_READER_HOST=0.0.0.0
      - VOID_READER_MAX_USERS=100
    healthcheck:
      test: ["CMD", "nc", "-z", "localhost", "23234"]
      interval: 30s
      timeout: 10s
      retries: 3

volumes:
  void_reader_data:
  ssh_keys:
```

#### Multi-stage Production Build
```dockerfile
# Dockerfile.prod
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o void-reader

FROM alpine:latest
RUN apk --no-cache add ca-certificates openssh-keygen
WORKDIR /app
COPY --from=builder /app/void-reader .
RUN addgroup -g 1001 -S voidreader && \
    adduser -u 1001 -S voidreader -G voidreader && \
    mkdir -p .ssh .void_reader_data && \
    ssh-keygen -t ed25519 -f .ssh/id_ed25519 -N "" && \
    chown -R voidreader:voidreader /app
USER voidreader
EXPOSE 23234
CMD ["./void-reader"]
```

### Kubernetes Deployment

#### Deployment Manifest
```yaml
# k8s/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: void-reader
  labels:
    app: void-reader
spec:
  replicas: 3
  selector:
    matchLabels:
      app: void-reader
  template:
    metadata:
      labels:
        app: void-reader
    spec:
      containers:
      - name: void-reader
        image: void-reader:latest
        ports:
        - containerPort: 23234
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "256Mi"
            cpu: "200m"
        volumeMounts:
        - name: book-content
          mountPath: /app/book1_void_reavers
          readOnly: true
        - name: user-data
          mountPath: /app/.void_reader_data
        env:
        - name: VOID_READER_HOST
          value: "0.0.0.0"
        - name: VOID_READER_MAX_USERS
          value: "50"
      volumes:
      - name: book-content
        configMap:
          name: book-content
      - name: user-data
        persistentVolumeClaim:
          claimName: void-reader-data
```

#### Service Configuration
```yaml
# k8s/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: void-reader-service
spec:
  selector:
    app: void-reader
  ports:
  - protocol: TCP
    port: 23234
    targetPort: 23234
  type: LoadBalancer
```

#### Ingress (TCP)
```yaml
# k8s/ingress.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: void-reader-ingress
  annotations:
    nginx.ingress.kubernetes.io/tcp-services-configmap: default/tcp-services
spec:
  rules:
  - host: void-reader.example.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: void-reader-service
            port:
              number: 23234
```

## 🔧 High Availability Deployment

### Load Balancing

#### HAProxy Configuration
```haproxy
# /etc/haproxy/haproxy.cfg
global
    daemon
    user haproxy
    group haproxy

defaults
    mode tcp
    timeout connect 5000ms
    timeout client 50000ms
    timeout server 50000ms

listen void_reader_cluster
    bind *:23234
    balance roundrobin
    server reader1 10.0.1.10:23234 check
    server reader2 10.0.1.11:23234 check
    server reader3 10.0.1.12:23234 check
```

#### Nginx TCP Load Balancer
```nginx
# /etc/nginx/nginx.conf
stream {
    upstream void_reader {
        least_conn;
        server 10.0.1.10:23234;
        server 10.0.1.11:23234;
        server 10.0.1.12:23234;
    }
    
    server {
        listen 23234;
        proxy_pass void_reader;
        proxy_timeout 1s;
        proxy_responses 1;
    }
}
```

### Session Persistence

#### Shared Storage for User Data
```bash
# NFS setup for shared user data
sudo apt install nfs-common

# Mount shared storage
sudo mount -t nfs nfs-server:/exports/void-reader-data /opt/void-reader/.void_reader_data

# Add to /etc/fstab
echo "nfs-server:/exports/void-reader-data /opt/void-reader/.void_reader_data nfs defaults 0 0" | sudo tee -a /etc/fstab
```

#### Redis for Session State
```go
// Add Redis support for distributed sessions
import "github.com/go-redis/redis/v8"

type RedisProgressManager struct {
    client *redis.Client
}

func NewRedisProgressManager(addr string) *RedisProgressManager {
    rdb := redis.NewClient(&redis.Options{
        Addr: addr,
    })
    
    return &RedisProgressManager{client: rdb}
}
```

### Database Backend

#### PostgreSQL Setup
```sql
-- Database schema for user progress
CREATE TABLE user_progress (
    username VARCHAR(255) PRIMARY KEY,
    current_chapter INTEGER NOT NULL DEFAULT 0,
    scroll_offset INTEGER NOT NULL DEFAULT 0,
    last_read TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    reading_time INTERVAL NOT NULL DEFAULT '0',
    progress_data JSONB
);

CREATE TABLE bookmarks (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    chapter INTEGER NOT NULL,
    scroll_offset INTEGER NOT NULL,
    note TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_bookmarks_username ON bookmarks(username);
```

## 📊 Monitoring and Observability

### Logging

#### Structured Logging
```go
import "github.com/sirupsen/logrus"

func setupLogging() {
    logrus.SetFormatter(&logrus.JSONFormatter{})
    logrus.SetLevel(logrus.InfoLevel)
    
    // Log to file
    file, err := os.OpenFile("void-reader.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err == nil {
        logrus.SetOutput(file)
    }
}
```

#### ELK Stack Integration
```yaml
# docker-compose.monitoring.yml
version: '3.8'

services:
  elasticsearch:
    image: docker.elastic.co/elasticsearch/elasticsearch:7.14.0
    environment:
      - discovery.type=single-node
    ports:
      - "9200:9200"

  logstash:
    image: docker.elastic.co/logstash/logstash:7.14.0
    volumes:
      - ./logstash.conf:/usr/share/logstash/pipeline/logstash.conf
    ports:
      - "5044:5044"

  kibana:
    image: docker.elastic.co/kibana/kibana:7.14.0
    ports:
      - "5601:5601"
    environment:
      - ELASTICSEARCH_HOSTS=http://elasticsearch:9200
```

### Metrics

#### Prometheus Integration
```go
import "github.com/prometheus/client_golang/prometheus"
import "github.com/prometheus/client_golang/prometheus/promhttp"

var (
    connectionsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "void_reader_connections_total",
            Help: "Total number of SSH connections",
        },
        []string{"status"},
    )
    
    activeUsers = prometheus.NewGauge(
        prometheus.GaugeOpts{
            Name: "void_reader_active_users",
            Help: "Number of currently active users",
        },
    )
)

func setupMetrics() {
    prometheus.MustRegister(connectionsTotal)
    prometheus.MustRegister(activeUsers)
    
    http.Handle("/metrics", promhttp.Handler())
    go http.ListenAndServe(":8080", nil)
}
```

#### Grafana Dashboard
```json
{
  "dashboard": {
    "title": "Void Reader Metrics",
    "panels": [
      {
        "title": "Active Users",
        "type": "graph",
        "targets": [
          {
            "expr": "void_reader_active_users",
            "legendFormat": "Active Users"
          }
        ]
      },
      {
        "title": "Connection Rate",
        "type": "graph",
        "targets": [
          {
            "expr": "rate(void_reader_connections_total[5m])",
            "legendFormat": "Connections/sec"
          }
        ]
      }
    ]
  }
}
```

### Health Checks

#### Application Health Endpoint
```go
func healthHandler(w http.ResponseWriter, r *http.Request) {
    health := map[string]interface{}{
        "status": "healthy",
        "timestamp": time.Now(),
        "version": "1.0.0",
        "active_connections": getActiveConnections(),
        "uptime": time.Since(startTime),
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(health)
}
```

#### Kubernetes Probes
```yaml
# Add to deployment.yaml
livenessProbe:
  httpGet:
    path: /health
    port: 8080
  initialDelaySeconds: 30
  periodSeconds: 10

readinessProbe:
  httpGet:
    path: /ready
    port: 8080
  initialDelaySeconds: 5
  periodSeconds: 5
```

## 🔐 Security Hardening

### Network Security

#### VPN Access Only
```bash
# Configure OpenVPN or WireGuard
# Allow SSH only from VPN subnet
sudo ufw allow from 10.8.0.0/24 to any port 23234
sudo ufw deny 23234
```

#### SSH Key Management
```bash
# Generate strong SSH keys
ssh-keygen -t ed25519 -b 521 -f /opt/void-reader/.ssh/id_ed25519

# Regular key rotation
0 0 1 * * root /opt/void-reader/scripts/rotate-ssh-keys.sh
```

### Application Security

#### Rate Limiting
```go
import "golang.org/x/time/rate"

type rateLimiter struct {
    visitors map[string]*visitor
    mu       sync.RWMutex
}

type visitor struct {
    limiter  *rate.Limiter
    lastSeen time.Time
}

func (rl *rateLimiter) getVisitor(ip string) *rate.Limiter {
    rl.mu.Lock()
    defer rl.mu.Unlock()
    
    v, exists := rl.visitors[ip]
    if !exists {
        limiter := rate.NewLimiter(1, 3) // 1 per second, burst of 3
        rl.visitors[ip] = &visitor{limiter, time.Now()}
        return limiter
    }
    
    v.lastSeen = time.Now()
    return v.limiter
}
```

#### Input Validation
```go
func validateUsername(username string) error {
    if len(username) > 50 {
        return errors.New("username too long")
    }
    
    matched, _ := regexp.MatchString("^[a-zA-Z0-9_-]+$", username)
    if !matched {
        return errors.New("invalid username format")
    }
    
    return nil
}
```

## 🚀 Performance Optimization

### Caching

#### In-Memory Cache
```go
import "github.com/patrickmn/go-cache"

var bookCache = cache.New(5*time.Minute, 10*time.Minute)

func getCachedBook(bookPath string) (*Book, error) {
    if cached, found := bookCache.Get(bookPath); found {
        return cached.(*Book), nil
    }
    
    book, err := LoadBook(bookPath)
    if err != nil {
        return nil, err
    }
    
    bookCache.Set(bookPath, book, cache.DefaultExpiration)
    return book, nil
}
```

#### Redis Cache
```go
func (rpm *RedisProgressManager) GetCachedProgress(username string) (*UserProgress, error) {
    data, err := rpm.client.Get(ctx, "progress:"+username).Result()
    if err == redis.Nil {
        return nil, errors.New("progress not found")
    } else if err != nil {
        return nil, err
    }
    
    var progress UserProgress
    err = json.Unmarshal([]byte(data), &progress)
    return &progress, err
}
```

### Connection Optimization

#### Connection Pooling
```go
type ConnectionPool struct {
    pool chan net.Conn
    maxConnections int
}

func NewConnectionPool(maxConnections int) *ConnectionPool {
    return &ConnectionPool{
        pool: make(chan net.Conn, maxConnections),
        maxConnections: maxConnections,
    }
}
```

### Resource Management

#### Memory Limits
```bash
# Systemd service limits
echo "MemoryMax=512M" | sudo tee -a /etc/systemd/system/void-reader.service
echo "CPUQuota=50%" | sudo tee -a /etc/systemd/system/void-reader.service
sudo systemctl daemon-reload
sudo systemctl restart void-reader
```

## 📋 Deployment Checklist

### Pre-Deployment
- [ ] System requirements verified
- [ ] Dependencies installed
- [ ] Security patches applied
- [ ] Firewall configured
- [ ] SSL certificates obtained (if needed)
- [ ] DNS records configured
- [ ] Monitoring setup ready

### Deployment
- [ ] Application built successfully
- [ ] Configuration files reviewed
- [ ] Service user created
- [ ] File permissions set correctly
- [ ] Systemd service installed
- [ ] Service starts without errors
- [ ] Health checks passing
- [ ] Log files created and writable

### Post-Deployment
- [ ] Connectivity tested from multiple clients
- [ ] User data persistence verified
- [ ] Backup system tested
- [ ] Monitoring alerts configured
- [ ] Performance baseline established
- [ ] Security scan completed
- [ ] Documentation updated

### Production Readiness
- [ ] Load testing completed
- [ ] Disaster recovery plan tested
- [ ] Security hardening applied
- [ ] Compliance requirements met
- [ ] Team training completed
- [ ] Runbooks updated

## 🆘 Deployment Troubleshooting

### Common Issues

#### Service Won't Start
```bash
# Check service status
sudo systemctl status void-reader

# Check logs
sudo journalctl -u void-reader -f --no-pager

# Check file permissions
ls -la /opt/void-reader/
```

#### Connection Refused
```bash
# Check if service is listening
sudo netstat -tlnp | grep 23234

# Check firewall
sudo ufw status
sudo iptables -L

# Test local connection
telnet localhost 23234
```

#### Performance Issues
```bash
# Check resource usage
top -p $(pgrep void-reader)
iostat -x 1

# Check connection limits
ss -tuln | grep 23234
```

---

**Deployment Complete!** 🚀✨

Your Void Reavers SSH Reader is now ready for production use. Monitor the system closely during the first few days and adjust configurations as needed.

**Next Steps:**
- Set up monitoring and alerting
- Create backup and disaster recovery procedures  
- Plan for scaling as user base grows
- Consider adding new features based on user feedback