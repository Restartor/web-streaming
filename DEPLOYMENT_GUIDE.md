# Deployment Guide — Web Streaming Backend

Step-by-step panduan deployment dari development ke production dengan best practices.

---

## 📋 Pre-Deployment Checklist

### Code Quality
- [ ] All linting passed (`golangci-lint`)
- [ ] All tests passing (`go test ./...`)
- [ ] Code reviewed & approved
- [ ] No sensitive data in code (API keys, passwords)
- [ ] Comments removed / cleaned up

### Database
- [ ] Migration scripts created & tested
- [ ] Backup strategy documented
- [ ] Connection pooling configured
- [ ] Indexes created on frequently queried columns

### Infrastructure
- [ ] SSL/TLS certificate ready
- [ ] Load balancer configured
- [ ] Health check endpoint implemented
- [ ] Secrets manager setup (Vault, AWS Secrets, K8s Secrets)
- [ ] Monitoring & alerting configured

### Configuration
- [ ] Environment variables documented
- [ ] Rate limits tuned for expected load
- [ ] CORS origins configured
- [ ] JWT secret generated (min 32 chars, random)
- [ ] Refresh token secret generated

### Security
- [ ] Input validation enabled
- [ ] CSRF protection configured
- [ ] Rate limiting active
- [ ] SQL injection prevention verified (using GORM ORM)
- [ ] Password hashing using bcrypt

---

## 🐳 Docker Deployment

### Step 1: Build Docker Image

**Dockerfile (production-optimized):**
```dockerfile
# Build stage
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Build with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -a -installsuffix cgo \
    -o /bin/web-streaming ./backend

# Runtime stage
FROM alpine:3.18
RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /bin/web-streaming /bin/web-streaming

ENV GIN_MODE=release
EXPOSE 1010

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
  CMD /bin/web-streaming health || exit 1

ENTRYPOINT ["/bin/web-streaming"]
```

**Build the image:**
```bash
# From root directory
docker build -f Dockerfile -t web-streaming:latest .
docker tag web-streaming:latest your-registry.azurecr.io/web-streaming:latest
docker push your-registry.azurecr.io/web-streaming:latest

# Tag with version
docker tag web-streaming:latest your-registry.azurecr.io/web-streaming:v1.0.0
docker push your-registry.azurecr.io/web-streaming:v1.0.0
```

### Step 2: Run with docker-compose (Staging)

**docker-compose.prod.yml:**
```yaml
version: '3.8'

services:
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: ${DB_NAME}
      POSTGRES_USER: ${DB_USER}
      POSTGRES_PASSWORD: ${DB_PASSWORD}
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${DB_USER}"]
      interval: 10s
      timeout: 5s
      retries: 5
    restart: unless-stopped

  app:
    image: your-registry.azurecr.io/web-streaming:latest
    environment:
      DB_HOST: postgres
      DB_PORT: 5432
      DB_USER: ${DB_USER}
      DB_PASSWORD: ${DB_PASSWORD}
      DB_NAME: ${DB_NAME}
      DB_SSLMODE: require
      JWT_SECRET: ${JWT_SECRET}
      REFRESH_TOKEN_SECRET: ${REFRESH_TOKEN_SECRET}
      ACCESS_TOKEN_DURATION: ${ACCESS_TOKEN_DURATION}
      REFRESH_TOKEN_DURATION: ${REFRESH_TOKEN_DURATION}
      ALLOWED_ORIGINS: ${ALLOWED_ORIGINS}
      GIN_MODE: release
    ports:
      - "1010:1010"
    depends_on:
      postgres:
        condition: service_healthy
    healthcheck:
      test: ["CMD", "wget", "--quiet", "--tries=1", "--spider", "http://localhost:1010/health"]
      interval: 30s
      timeout: 3s
      retries: 3
    restart: unless-stopped
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"

volumes:
  postgres_data:
    driver: local
```

**Environment file (.env.prod):**
```bash
DB_USER=prod_postgres
DB_PASSWORD=your-strong-password-here
DB_NAME=web_streaming_prod
JWT_SECRET=your-very-long-random-jwt-secret-min-32-chars
REFRESH_TOKEN_SECRET=another-long-random-secret-min-32-chars
ACCESS_TOKEN_DURATION=15m
REFRESH_TOKEN_DURATION=168h
ALLOWED_ORIGINS=https://yourdomain.com,https://www.yourdomain.com
```

**Deploy:**
```bash
docker-compose -f docker-compose.prod.yml \
  --env-file .env.prod up -d

# View logs
docker-compose -f docker-compose.prod.yml logs -f app

# Stop
docker-compose -f docker-compose.prod.yml down
```

---

## ☸️ Kubernetes Deployment

### Step 1: Create Kubernetes Manifests

**k8s/namespace.yaml:**
```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: web-streaming
```

**k8s/secrets.yaml:**
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: db-credentials
  namespace: web-streaming
type: Opaque
data:
  db-host: cG9zdGdyZXM=  # base64 encoded "postgres"
  db-port: NTQzMg==  # "5432"
  db-user: cHJvZF91c2Vy  # Your prod user
  db-password: eW91ci1zdHJvbmctcGFzc3dvcmQ=  # Your prod password
  db-name: d2ViX3N0cmVhbWluZ19wcm9k  # "web_streaming_prod"

---
apiVersion: v1
kind: Secret
metadata:
  name: jwt-credentials
  namespace: web-streaming
type: Opaque
data:
  jwt-secret: eW91ci12ZXJ5LWxvbmctcmFuZG9tLWp3dC1zZWNyZXQ=  # Your JWT secret
  refresh-secret: YW5vdGhlci1sb25nLXJhbmRvbS1zZWNyZXQ=  # Your refresh secret
```

**k8s/configmap.yaml:**
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: app-config
  namespace: web-streaming
data:
  ACCESS_TOKEN_DURATION: "15m"
  REFRESH_TOKEN_DURATION: "168h"
  ALLOWED_ORIGINS: "https://yourdomain.com,https://www.yourdomain.com"
  GIN_MODE: "release"
```

**k8s/deployment.yaml:**
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: web-streaming-api
  namespace: web-streaming
  labels:
    app: web-streaming-api
spec:
  replicas: 3
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 0
  selector:
    matchLabels:
      app: web-streaming-api
  template:
    metadata:
      labels:
        app: web-streaming-api
      annotations:
        prometheus.io/scrape: "true"
        prometheus.io/port: "1010"
        prometheus.io/path: "/metrics"
    spec:
      serviceAccountName: web-streaming
      securityContext:
        runAsNonRoot: true
        runAsUser: 1000
        fsGroup: 1000
      containers:
      - name: api
        image: your-registry.azurecr.io/web-streaming:v1.0.0
        imagePullPolicy: IfNotPresent
        ports:
        - name: http
          containerPort: 1010
          protocol: TCP
        env:
        - name: DB_HOST
          valueFrom:
            secretKeyRef:
              name: db-credentials
              key: db-host
        - name: DB_PORT
          valueFrom:
            secretKeyRef:
              name: db-credentials
              key: db-port
        - name: DB_USER
          valueFrom:
            secretKeyRef:
              name: db-credentials
              key: db-user
        - name: DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: db-credentials
              key: db-password
        - name: DB_NAME
          valueFrom:
            secretKeyRef:
              name: db-credentials
              key: db-name
        - name: DB_SSLMODE
          value: "require"
        - name: JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: jwt-credentials
              key: jwt-secret
        - name: REFRESH_TOKEN_SECRET
          valueFrom:
            secretKeyRef:
              name: jwt-credentials
              key: refresh-secret
        - name: ACCESS_TOKEN_DURATION
          valueFrom:
            configMapKeyRef:
              name: app-config
              key: ACCESS_TOKEN_DURATION
        - name: REFRESH_TOKEN_DURATION
          valueFrom:
            configMapKeyRef:
              name: app-config
              key: REFRESH_TOKEN_DURATION
        - name: ALLOWED_ORIGINS
          valueFrom:
            configMapKeyRef:
              name: app-config
              key: ALLOWED_ORIGINS
        - name: GIN_MODE
          valueFrom:
            configMapKeyRef:
              name: app-config
              key: GIN_MODE
        
        livenessProbe:
          httpGet:
            path: /health
            port: http
          initialDelaySeconds: 30
          periodSeconds: 30
          timeoutSeconds: 3
          failureThreshold: 3

        readinessProbe:
          httpGet:
            path: /health
            port: http
          initialDelaySeconds: 10
          periodSeconds: 10
          timeoutSeconds: 3
          failureThreshold: 3

        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"

        securityContext:
          allowPrivilegeEscalation: false
          readOnlyRootFilesystem: true
          capabilities:
            drop:
              - ALL
        
        volumeMounts:
        - name: tmp
          mountPath: /tmp

      volumes:
      - name: tmp
        emptyDir: {}

      affinity:
        podAntiAffinity:
          preferredDuringSchedulingIgnoredDuringExecution:
          - weight: 100
            podAffinityTerm:
              labelSelector:
                matchExpressions:
                - key: app
                  operator: In
                  values:
                  - web-streaming-api
              topologyKey: kubernetes.io/hostname
```

**k8s/service.yaml:**
```yaml
apiVersion: v1
kind: Service
metadata:
  name: web-streaming-api-svc
  namespace: web-streaming
  labels:
    app: web-streaming-api
spec:
  type: LoadBalancer
  ports:
  - port: 80
    targetPort: 1010
    protocol: TCP
    name: http
  selector:
    app: web-streaming-api
  sessionAffinity: ClientIP
  sessionAffinityConfig:
    clientIP:
      timeoutSeconds: 10800
```

**k8s/hpa.yaml (Auto-scaling):**
```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: web-streaming-hpa
  namespace: web-streaming
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: web-streaming-api
  minReplicas: 3
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: 80
  behavior:
    scaleDown:
      stabilizationWindowSeconds: 300
      policies:
      - type: Percent
        value: 50
        periodSeconds: 60
    scaleUp:
      stabilizationWindowSeconds: 0
      policies:
      - type: Percent
        value: 100
        periodSeconds: 30
      - type: Pods
        value: 2
        periodSeconds: 60
      selectPolicy: Max
```

### Step 2: Deploy to Kubernetes

```bash
# Create namespace
kubectl apply -f k8s/namespace.yaml

# Create secrets (update values first!)
kubectl apply -f k8s/secrets.yaml

# Create configmap
kubectl apply -f k8s/configmap.yaml

# Apply all manifests
kubectl apply -f k8s/

# Verify deployment
kubectl get pods -n web-streaming
kubectl get svc -n web-streaming
kubectl logs -f deployment/web-streaming-api -n web-streaming

# Port forward to test locally
kubectl port-forward svc/web-streaming-api-svc 8080:80 -n web-streaming

# Update image (rolling update)
kubectl set image deployment/web-streaming-api \
  api=your-registry.azurecr.io/web-streaming:v1.0.1 \
  -n web-streaming

# Check rollout status
kubectl rollout status deployment/web-streaming-api -n web-streaming
```

---

## 🔄 CI/CD Pipeline (GitHub Actions)

**.github/workflows/deploy.yml:**
```yaml
name: Deploy to Production

on:
  push:
    tags:
      - 'v*'
  workflow_dispatch:

jobs:
  test:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:15-alpine
        env:
          POSTGRES_DB: web_streaming_test
          POSTGRES_USER: postgres
          POSTGRES_PASSWORD: postgres
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
        ports:
          - 5432:5432

    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v4
        with:
          go-version: '1.25'
      - run: go mod download
      - run: golangci-lint run ./...
      - run: go test -v -race -coverprofile=coverage.out ./...
      - uses: codecov/codecov-action@v3

  build-and-push:
    needs: test
    runs-on: ubuntu-latest
    permissions:
      contents: read
      packages: write

    steps:
      - uses: actions/checkout@v4

      - name: Set version
        id: version
        run: echo "VERSION=${GITHUB_REF#refs/tags/}" >> $GITHUB_OUTPUT

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v2

      - name: Login to Azure Container Registry
        uses: docker/login-action@v2
        with:
          registry: your-registry.azurecr.io
          username: ${{ secrets.ACR_USERNAME }}
          password: ${{ secrets.ACR_PASSWORD }}

      - name: Build and push
        uses: docker/build-push-action@v4
        with:
          context: ./backend
          push: true
          tags: |
            your-registry.azurecr.io/web-streaming:latest
            your-registry.azurecr.io/web-streaming:${{ steps.version.outputs.VERSION }}
          cache-from: type=registry,ref=your-registry.azurecr.io/web-streaming:buildcache
          cache-to: type=registry,ref=your-registry.azurecr.io/web-streaming:buildcache,mode=max

  deploy-staging:
    needs: build-and-push
    runs-on: ubuntu-latest
    environment:
      name: staging

    steps:
      - uses: actions/checkout@v4

      - name: Deploy to Azure Container Instances
        env:
          AZURE_RESOURCE_GROUP: web-streaming-staging
          AZURE_CONTAINER_INSTANCE: web-streaming-api-staging
        run: |
          az login --service-principal \
            -u ${{ secrets.AZURE_CLIENT_ID }} \
            -p ${{ secrets.AZURE_CLIENT_SECRET }} \
            --tenant ${{ secrets.AZURE_TENANT_ID }}
          
          az container create \
            --resource-group $AZURE_RESOURCE_GROUP \
            --name $AZURE_CONTAINER_INSTANCE \
            --image your-registry.azurecr.io/web-streaming:latest \
            --environment-variables \
              DB_HOST=${{ secrets.STAGING_DB_HOST }} \
              DB_USER=${{ secrets.STAGING_DB_USER }} \
              JWT_SECRET=${{ secrets.STAGING_JWT_SECRET }} \
            --registry-login-server your-registry.azurecr.io \
            --registry-username ${{ secrets.ACR_USERNAME }} \
            --registry-password ${{ secrets.ACR_PASSWORD }} \
            --ports 1010 \
            --cpu 1 --memory 1 \
            --restart-policy OnFailure

      - name: Smoke tests
        run: |
          sleep 10
          curl -f http://${{ secrets.STAGING_API_HOST }}/api/v1/health || exit 1

  deploy-production:
    needs: deploy-staging
    runs-on: ubuntu-latest
    environment:
      name: production
    if: github.event_name == 'push'

    steps:
      - uses: actions/checkout@v4

      - name: Deploy to Kubernetes
        env:
          KUBECONFIG_CONTENT: ${{ secrets.KUBE_CONFIG }}
        run: |
          mkdir -p $HOME/.kube
          echo "$KUBECONFIG_CONTENT" > $HOME/.kube/config
          
          VERSION=${GITHUB_REF#refs/tags/}
          kubectl set image deployment/web-streaming-api \
            api=your-registry.azurecr.io/web-streaming:$VERSION \
            -n web-streaming
          
          kubectl rollout status deployment/web-streaming-api \
            -n web-streaming \
            --timeout=5m

      - name: Verify deployment
        run: |
          kubectl get pods -n web-streaming
          kubectl get svc -n web-streaming
```

---

## 🚀 Manual Deployment Steps

### On Production Server

```bash
# 1. SSH ke server
ssh deploy@production.example.com

# 2. Pull latest code
cd /var/www/web-streaming
git fetch origin
git checkout tags/v1.0.0

# 3. Backup database
pg_dump -U postgres web_streaming > backup_$(date +%Y%m%d_%H%M%S).sql

# 4. Run migrations
export DATABASE_URL="postgres://postgres:password@localhost/web_streaming?sslmode=require"
migrate -path ./backend/migrations -database $DATABASE_URL up

# 5. Build binary
cd backend
go build -o web-streaming .

# 6. Stop old service
sudo systemctl stop web-streaming

# 7. Replace binary
sudo cp web-streaming /usr/local/bin/web-streaming

# 8. Start service
sudo systemctl start web-streaming

# 9. Verify
curl http://localhost:1010/health

# 10. Check logs
sudo journalctl -u web-streaming -f
```

### Systemd Service File (/etc/systemd/system/web-streaming.service)
```ini
[Unit]
Description=Web Streaming API
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/var/www/web-streaming
ExecStart=/usr/local/bin/web-streaming
Restart=on-failure
RestartSec=5s
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
```

---

## 📊 Monitoring & Logging

### Log Aggregation (ELK Stack)
```yaml
# filebeat configuration
filebeat.inputs:
- type: log
  enabled: true
  paths:
    - /var/log/web-streaming.log

output.elasticsearch:
  hosts: ["https://elasticsearch.example.com:9200"]
  username: "elastic"
  password: "${ELASTIC_PASSWORD}"
  ssl.verification_mode: certificate

processors:
  - add_docker_metadata: ~
  - add_kubernetes_metadata: ~
```

### Prometheus Scrape Config
```yaml
scrape_configs:
  - job_name: 'web-streaming'
    static_configs:
      - targets: ['localhost:1010']
    metrics_path: '/metrics'
    scrape_interval: 30s
```

### Alert Rules (Alertmanager)
```yaml
groups:
- name: web_streaming
  rules:
  - alert: APIDown
    expr: up{job="web-streaming"} == 0
    for: 2m
    annotations:
      summary: "API is down"

  - alert: HighErrorRate
    expr: rate(http_requests_total{status=~"5.."}[5m]) > 0.05
    for: 5m
    annotations:
      summary: "High error rate detected"

  - alert: HighLatency
    expr: histogram_quantile(0.95, http_request_duration_seconds) > 1
    for: 10m
    annotations:
      summary: "High API latency"

  - alert: DatabaseConnectionPoolFull
    expr: db_connections_active >= db_connections_max
    for: 5m
    annotations:
      summary: "Database connection pool exhausted"
```

---

## ✅ Post-Deployment Verification

```bash
# 1. Check API health
curl https://api.yourdomain.com/api/v1/health

# 2. Test registration
curl -X POST https://api.yourdomain.com/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{"username":"test","email":"test@test.com","password":"Pass123"}'

# 3. Test login
curl -X POST https://api.yourdomain.com/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@test.com","password":"Pass123"}'

# 4. Verify logs
kubectl logs deployment/web-streaming-api -n web-streaming

# 5. Check metrics
curl http://localhost:9090/api/v1/targets

# 6. Database connectivity
psql -h prod-db.example.com -U postgres -d web_streaming -c "SELECT count(*) FROM users;"

# 7. Test external API from frontend
# Monitor network tab in browser for CORS errors
# Check response formats match expectations
```

---

## 🔙 Rollback Procedure

```bash
# If deployment fails:

# Kubernetes rollback
kubectl rollout undo deployment/web-streaming-api -n web-streaming

# Check rollout history
kubectl rollout history deployment/web-streaming-api -n web-streaming

# Revert to specific revision
kubectl rollout undo deployment/web-streaming-api --to-revision=3 -n web-streaming

# Database rollback (if migration failed)
export DATABASE_URL="postgres://postgres:password@localhost/web_streaming"
migrate -path ./backend/migrations -database $DATABASE_URL down 1
```

---

## 📝 Deployment Summary

| Environment | Database | Secrets | Scaling | Monitoring |
|---|---|---|---|---|
| **Development** | Local PostgreSQL | .env file | 1 instance | Basic logging |
| **Staging** | RDS/Managed | AWS Secrets Manager | 2-3 instances | CloudWatch |
| **Production** | RDS HA | Azure Key Vault / Vault | 3-10+ auto-scaled | DataDog / ELK |

