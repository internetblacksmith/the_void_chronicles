# Troubleshooting Guide

Comprehensive troubleshooting guide for the Void Reavers SSH Reader covering common issues, diagnostics, and solutions.

## 🔍 Quick Diagnostics

### Health Check Commands

Run these commands to quickly assess system health:

```bash
# Check if service is running
sudo systemctl status void-reader

# Check if port is listening  
sudo netstat -tlnp | grep 23234
# Alternative: sudo ss -tlnp | grep 23234

# Test local connection
ssh localhost -p 23234

# Check logs
sudo journalctl -u void-reader -f --lines=50

# Check file permissions
ls -la /opt/void-reader/
ls -la /opt/void-reader/.void_reader_data/

# Check disk space
df -h /opt/void-reader/

# Check memory usage
free -h
ps aux | grep void-reader
```

### Log Analysis

#### Enable Debug Logging

Add to your application:
```go
import "github.com/charmbracelet/log"

func init() {
    log.SetLevel(log.DebugLevel)
    log.SetReportCaller(true)
}
```

#### Log Locations

| Deployment Type | Log Location |
|----------------|--------------|
| Systemd Service | `journalctl -u void-reader` |
| Docker | `docker logs void-reader` |
| Docker Compose | `docker-compose logs void-reader` |
| Manual | Console output or specified log file |
| Kubernetes | `kubectl logs deployment/void-reader` |

## 🚨 Common Issues

### Connection Issues

#### Issue: "Connection refused" when connecting via SSH

**Symptoms:**
```bash
$ ssh localhost -p 23234
ssh: connect to host localhost port 23234: Connection refused
```

**Diagnosis:**
```bash
# Check if service is running
sudo systemctl status void-reader

# Check if port is listening
sudo netstat -tlnp | grep 23234

# Check if binary exists and is executable
ls -la /opt/void-reader/void-reader
```

**Solutions:**

1. **Service not running:**
   ```bash
   sudo systemctl start void-reader
   sudo systemctl enable void-reader
   ```

2. **Port already in use:**
   ```bash
   # Find what's using the port
   sudo lsof -i :23234
   
   # Kill the process or change port in main.go
   sudo kill -9 <PID>
   ```

3. **Binary not executable:**
   ```bash
   sudo chmod +x /opt/void-reader/void-reader
   sudo chown voidreader:voidreader /opt/void-reader/void-reader
   ```

4. **Firewall blocking:**
   ```bash
   # Ubuntu/Debian
   sudo ufw allow 23234/tcp
   
   # CentOS/RHEL
   sudo firewall-cmd --permanent --add-port=23234/tcp
   sudo firewall-cmd --reload
   ```

#### Issue: SSH connection drops immediately

**Symptoms:**
```bash
$ ssh localhost -p 23234
Connection to localhost closed.
```

**Diagnosis:**
```bash
# Check logs for errors
sudo journalctl -u void-reader -f

# Test with verbose SSH
ssh -vvv localhost -p 23234
```

**Solutions:**

1. **SSH key issues:**
   ```bash
   # Regenerate SSH host key
   sudo rm /opt/void-reader/.ssh/id_ed25519*
   sudo -u voidreader ssh-keygen -t ed25519 -f /opt/void-reader/.ssh/id_ed25519 -N ""
   sudo systemctl restart void-reader
   ```

2. **Permission issues:**
   ```bash
   sudo chmod 600 /opt/void-reader/.ssh/id_ed25519
   sudo chmod 644 /opt/void-reader/.ssh/id_ed25519.pub
   sudo chown voidreader:voidreader /opt/void-reader/.ssh/id_ed25519*
   ```

3. **Application crash:**
   ```bash
   # Check for panic in logs
   sudo journalctl -u void-reader --since "1 hour ago" | grep -i panic
   
   # Restart service
   sudo systemctl restart void-reader
   ```

#### Issue: Can't connect from remote machines

**Symptoms:**
- Local connections work: `ssh localhost -p 23234` ✅
- Remote connections fail: `ssh server-ip -p 23234` ❌

**Diagnosis:**
```bash
# Check what interfaces the service binds to
sudo netstat -tlnp | grep 23234
# Should show 0.0.0.0:23234, not 127.0.0.1:23234

# Test from remote machine
telnet server-ip 23234
```

**Solutions:**

1. **Change bind address:**
   ```go
   // In main.go
   const (
       host = "0.0.0.0"    // Change from "localhost"
       port = "23234"
   )
   ```

2. **Firewall configuration:**
   ```bash
   # Allow from specific network
   sudo ufw allow from 192.168.1.0/24 to any port 23234
   
   # Or allow from anywhere (less secure)
   sudo ufw allow 23234/tcp
   ```

3. **Cloud security groups:**
   ```bash
   # AWS
   aws ec2 authorize-security-group-ingress \
       --group-id sg-xxxxxxxxx \
       --protocol tcp \
       --port 23234 \
       --cidr 0.0.0.0/0
   ```

### Application Issues

#### Issue: Book content not loading

**Symptoms:**
```
Error Loading Book
Could not load book content
```

**Diagnosis:**
```bash
# Check book directory structure
ls -la /opt/void-reader/book1_void_reavers/
ls -la /opt/void-reader/book1_void_reavers/markdown/

# Check file permissions
sudo -u voidreader ls -la /opt/void-reader/book1_void_reavers/markdown/

# Check file content
head /opt/void-reader/book1_void_reavers/markdown/chapter01.md
```

**Solutions:**

1. **Missing book directory:**
   ```bash
   # Copy book content
   sudo cp -r /path/to/source/book1_void_reavers /opt/void-reader/
   sudo chown -R voidreader:voidreader /opt/void-reader/book1_void_reavers
   ```

2. **Permission issues:**
   ```bash
   sudo chmod -R 644 /opt/void-reader/book1_void_reavers/
   sudo chmod 755 /opt/void-reader/book1_void_reavers/
   sudo chmod 755 /opt/void-reader/book1_void_reavers/markdown/
   ```

3. **File encoding issues:**
   ```bash
   # Check file encoding
   file -i /opt/void-reader/book1_void_reavers/markdown/chapter01.md
   
   # Convert if needed
   iconv -f ISO-8859-1 -t UTF-8 chapter01.md > chapter01.md.utf8
   ```

#### Issue: User progress not saving

**Symptoms:**
- Reading position resets on reconnection
- Bookmarks disappear
- Progress screen shows no data

**Diagnosis:**
```bash
# Check data directory
ls -la /opt/void-reader/.void_reader_data/

# Check permissions
sudo -u voidreader touch /opt/void-reader/.void_reader_data/test.json

# Check disk space
df -h /opt/void-reader/

# Check for errors in logs
sudo journalctl -u void-reader | grep -i "progress\|save\|error"
```

**Solutions:**

1. **Permission issues:**
   ```bash
   sudo mkdir -p /opt/void-reader/.void_reader_data
   sudo chown -R voidreader:voidreader /opt/void-reader/.void_reader_data
   sudo chmod 755 /opt/void-reader/.void_reader_data
   ```

2. **Disk space full:**
   ```bash
   # Clean up old logs
   sudo journalctl --vacuum-time=7d
   
   # Clean up old progress files
   sudo find /opt/void-reader/.void_reader_data -name "*.json" -mtime +30 -delete
   ```

3. **Corrupted progress files:**
   ```bash
   # Check for malformed JSON
   cd /opt/void-reader/.void_reader_data
   for file in *.json; do
       echo "Checking $file:"
       python3 -m json.tool "$file" > /dev/null || echo "Invalid JSON: $file"
   done
   
   # Remove corrupted files
   sudo rm /opt/void-reader/.void_reader_data/corrupted_user.json
   ```

#### Issue: Application crashes or panics

**Symptoms:**
```
panic: runtime error: index out of range
Service exits unexpectedly
```

**Diagnosis:**
```bash
# Check for panic messages
sudo journalctl -u void-reader | grep -A 10 -B 10 "panic"

# Check system resources
free -h
df -h
```

**Solutions:**

1. **Out of memory:**
   ```bash
   # Increase memory limits
   sudo systemctl edit void-reader
   
   # Add:
   [Service]
   MemoryMax=1G
   ```

2. **Out of disk space:**
   ```bash
   # Clean up space
   sudo apt autoremove
   sudo apt autoclean
   sudo journalctl --vacuum-size=100M
   ```

3. **File corruption:**
   ```bash
   # Restore from backup
   sudo systemctl stop void-reader
   sudo cp -r /backup/void-reader/* /opt/void-reader/
   sudo systemctl start void-reader
   ```

### Performance Issues

#### Issue: Slow response times

**Symptoms:**
- Slow to connect via SSH
- Laggy interface response
- High CPU usage

**Diagnosis:**
```bash
# Check CPU and memory usage
top -p $(pgrep void-reader)
htop

# Check I/O usage
iostat -x 1

# Check network
netstat -i
ss -tuln | grep 23234

# Profile the application
go tool pprof http://localhost:6060/debug/pprof/profile
```

**Solutions:**

1. **High CPU usage:**
   ```bash
   # Limit CPU usage
   sudo systemctl edit void-reader
   
   # Add:
   [Service]
   CPUQuota=50%
   ```

2. **Memory issues:**
   ```bash
   # Add swap if needed
   sudo fallocate -l 2G /swapfile
   sudo chmod 600 /swapfile
   sudo mkswap /swapfile
   sudo swapon /swapfile
   ```

3. **I/O bottleneck:**
   ```bash
   # Move to faster storage
   sudo rsync -av /opt/void-reader/ /fast/storage/void-reader/
   sudo systemctl stop void-reader
   sudo mv /opt/void-reader /opt/void-reader.old
   sudo ln -s /fast/storage/void-reader /opt/void-reader
   sudo systemctl start void-reader
   ```

#### Issue: Too many connections

**Symptoms:**
```
accept tcp :23234: too many open files
Connection limit reached
```

**Diagnosis:**
```bash
# Check current connections
ss -tuln | grep 23234 | wc -l

# Check file descriptor limits
ulimit -n
cat /proc/$(pgrep void-reader)/limits | grep "Max open files"

# Check system limits
cat /proc/sys/fs/file-max
```

**Solutions:**

1. **Increase file descriptor limits:**
   ```bash
   # For systemd service
   sudo systemctl edit void-reader
   
   # Add:
   [Service]
   LimitNOFILE=65536
   ```

2. **System-wide limits:**
   ```bash
   # Edit /etc/security/limits.conf
   echo "voidreader soft nofile 65536" | sudo tee -a /etc/security/limits.conf
   echo "voidreader hard nofile 65536" | sudo tee -a /etc/security/limits.conf
   ```

3. **Application-level connection limiting:**
   ```go
   // Add to main.go
   const maxConnections = 100
   
   var connectionCount int32
   
   func checkConnectionLimit() bool {
       return atomic.LoadInt32(&connectionCount) < maxConnections
   }
   ```

## 🐳 Container-Specific Issues

### Docker Issues

#### Issue: Container won't start

**Diagnosis:**
```bash
# Check container status
docker ps -a

# Check container logs
docker logs void-reader

# Check image
docker images | grep void-reader

# Check Docker daemon
sudo systemctl status docker
```

**Solutions:**

1. **Build issues:**
   ```bash
   # Rebuild image
   docker build --no-cache -t void-reader .
   
   # Check Dockerfile syntax
   docker build --dry-run -t void-reader .
   ```

2. **Port conflicts:**
   ```bash
   # Check what's using port
   sudo lsof -i :23234
   
   # Use different port
   docker run -p 23235:23234 void-reader
   ```

3. **Volume mount issues:**
   ```bash
   # Check volume permissions
   ls -la ./book1_void_reavers/
   
   # Fix permissions
   sudo chown -R 1001:1001 ./book1_void_reavers/
   ```

#### Issue: Docker Compose networking problems

**Diagnosis:**
```bash
# Check services
docker-compose ps

# Check networks
docker network ls
docker-compose logs

# Test connectivity
docker-compose exec void-reader nc -z localhost 23234
```

**Solutions:**

1. **Service discovery:**
   ```yaml
   # In docker-compose.yml
   services:
     void-reader:
       hostname: void-reader
       networks:
         - void-network
   
   networks:
     void-network:
       driver: bridge
   ```

2. **Port mapping:**
   ```yaml
   services:
     void-reader:
       ports:
         - "23234:23234"
       expose:
         - "23234"
   ```

### Kubernetes Issues

#### Issue: Pod crashes or restarts

**Diagnosis:**
```bash
# Check pod status
kubectl get pods -l app=void-reader

# Check pod logs
kubectl logs -l app=void-reader --previous

# Describe pod
kubectl describe pod <pod-name>

# Check events
kubectl get events --sort-by=.metadata.creationTimestamp
```

**Solutions:**

1. **Resource limits:**
   ```yaml
   # Increase resources in deployment.yaml
   resources:
     limits:
       memory: "512Mi"
       cpu: "500m"
     requests:
       memory: "256Mi"
       cpu: "250m"
   ```

2. **Health check failures:**
   ```yaml
   # Adjust health checks
   livenessProbe:
     httpGet:
       path: /health
       port: 8080
     initialDelaySeconds: 60  # Increase delay
     timeoutSeconds: 10
   ```

3. **Storage issues:**
   ```yaml
   # Add persistent volume
   volumeMounts:
   - name: user-data
     mountPath: /app/.void_reader_data
   volumes:
   - name: user-data
     persistentVolumeClaim:
       claimName: void-reader-pvc
   ```

## 🔐 Security Issues

### Authentication Problems

#### Issue: SSH key warnings or rejections

**Symptoms:**
```
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@    WARNING: REMOTE HOST IDENTIFICATION HAS CHANGED!     @
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
```

**Solutions:**

1. **Remove old host key:**
   ```bash
   ssh-keygen -R "[localhost]:23234"
   ```

2. **Accept new key:**
   ```bash
   ssh -o StrictHostKeyChecking=no localhost -p 23234
   ```

3. **Use consistent host keys:**
   ```bash
   # Backup your host key
   sudo cp /opt/void-reader/.ssh/id_ed25519 /secure/backup/
   
   # Restore after reinstall
   sudo cp /secure/backup/id_ed25519 /opt/void-reader/.ssh/
   ```

### Permission Issues

#### Issue: File permission errors

**Symptoms:**
```
permission denied: .void_reader_data/user.json
cannot create directory: permission denied
```

**Solutions:**

1. **Fix ownership:**
   ```bash
   sudo chown -R voidreader:voidreader /opt/void-reader/
   ```

2. **Fix permissions:**
   ```bash
   # Application binary
   sudo chmod 755 /opt/void-reader/void-reader
   
   # Data directory
   sudo chmod 755 /opt/void-reader/.void_reader_data
   sudo chmod 644 /opt/void-reader/.void_reader_data/*.json
   
   # SSH keys
   sudo chmod 600 /opt/void-reader/.ssh/id_ed25519
   sudo chmod 644 /opt/void-reader/.ssh/id_ed25519.pub
   ```

3. **SELinux issues (CentOS/RHEL):**
   ```bash
   # Check SELinux status
   sestatus
   
   # Set context
   sudo setsebool -P ssh_sysadm_login on
   sudo semanage fcontext -a -t ssh_exec_t "/opt/void-reader/void-reader"
   sudo restorecon -v /opt/void-reader/void-reader
   ```

## 🔧 Advanced Diagnostics

### Performance Profiling

#### CPU Profiling
```go
// Add to main.go
import _ "net/http/pprof"

func init() {
    go func() {
        log.Println(http.ListenAndServe("localhost:6060", nil))
    }()
}
```

Access profiling:
```bash
# CPU profile
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30

# Memory profile
go tool pprof http://localhost:6060/debug/pprof/heap

# Goroutine profile
go tool pprof http://localhost:6060/debug/pprof/goroutine
```

#### Memory Analysis
```bash
# Check memory usage
ps -eo pid,ppid,cmd,%mem,%cpu --sort=-%mem | head

# Monitor memory over time
while true; do
    echo "$(date): $(ps -o pid,vsz,rss,comm -p $(pgrep void-reader))"
    sleep 10
done
```

#### Network Analysis
```bash
# Monitor connections
watch 'ss -tuln | grep 23234'

# Check connection states
ss -tan | grep 23234 | awk '{print $1}' | sort | uniq -c

# Network throughput
iftop -i any -P
```

### Log Analysis

#### Parse Structured Logs
```bash
# Extract error messages
sudo journalctl -u void-reader --since "1 hour ago" | jq -r 'select(.level == "error") | .message'

# Count connection attempts
sudo journalctl -u void-reader --since today | grep "connection" | wc -l

# Find slow requests
sudo journalctl -u void-reader | grep "slow" | tail -20
```

#### Log Rotation Issues
```bash
# Check log rotation
sudo logrotate -d /etc/logrotate.d/void-reader

# Manual rotation
sudo logrotate -f /etc/logrotate.d/void-reader
```

### Database Diagnostics (if using database backend)

```sql
-- Check user progress records
SELECT username, current_chapter, last_read 
FROM user_progress 
ORDER BY last_read DESC 
LIMIT 10;

-- Find corrupted data
SELECT username, progress_data 
FROM user_progress 
WHERE progress_data IS NULL 
OR NOT jsonb_valid(progress_data);

-- Performance analysis
EXPLAIN ANALYZE SELECT * FROM bookmarks WHERE username = 'testuser';
```

## 🛠️ Recovery Procedures

### Data Recovery

#### Backup and Restore
```bash
# Create backup
sudo tar -czf void-reader-backup-$(date +%Y%m%d).tar.gz \
    /opt/void-reader/.void_reader_data \
    /opt/void-reader/.ssh

# Restore from backup
sudo systemctl stop void-reader
sudo tar -xzf void-reader-backup-20240115.tar.gz -C /
sudo chown -R voidreader:voidreader /opt/void-reader/
sudo systemctl start void-reader
```

#### Corrupt Data Cleanup
```bash
# Find and remove corrupted progress files
cd /opt/void-reader/.void_reader_data
for file in *.json; do
    if ! python3 -m json.tool "$file" > /dev/null 2>&1; then
        echo "Corrupted: $file"
        sudo mv "$file" "$file.corrupted"
    fi
done
```

### Service Recovery

#### Complete Service Reset
```bash
# Stop service
sudo systemctl stop void-reader

# Backup data
sudo cp -r /opt/void-reader/.void_reader_data /tmp/backup-user-data

# Remove and reinstall
sudo rm -rf /opt/void-reader
sudo ./deploy.sh

# Restore user data
sudo cp -r /tmp/backup-user-data/* /opt/void-reader/.void_reader_data/
sudo chown -R voidreader:voidreader /opt/void-reader/.void_reader_data

# Start service
sudo systemctl start void-reader
```

#### Emergency Restart Script
```bash
#!/bin/bash
# emergency-restart.sh

echo "Emergency restart procedure starting..."

# Stop service
sudo systemctl stop void-reader

# Kill any remaining processes
sudo pkill -f void-reader

# Check for port conflicts
if lsof -i :23234; then
    echo "Port 23234 still in use, killing processes..."
    sudo fuser -k 23234/tcp
fi

# Start service
sudo systemctl start void-reader

# Wait and check
sleep 5
if sudo systemctl is-active --quiet void-reader; then
    echo "Service restarted successfully"
else
    echo "Service failed to start, check logs:"
    sudo journalctl -u void-reader --lines=20
fi
```

## 📞 Getting Help

### Self-Help Resources

1. **Check application logs first**
2. **Review configuration files**
3. **Test with minimal setup**
4. **Search existing documentation**

### Gathering Information for Support

When reporting issues, include:

```bash
# System information
uname -a
cat /etc/os-release

# Service status
sudo systemctl status void-reader --no-pager -l

# Recent logs (last 50 lines)
sudo journalctl -u void-reader --lines=50 --no-pager

# Configuration
ls -la /opt/void-reader/
ps aux | grep void-reader

# Network status
sudo netstat -tlnp | grep 23234
sudo ss -tlnp | grep 23234

# Resource usage
df -h
free -h
```

### Community Resources

- **Documentation**: Complete guides in `/docs` directory
- **GitHub Issues**: Report bugs and feature requests
- **Discussions**: Community help and questions
- **Wiki**: Community-contributed solutions

---

**Troubleshooting Complete!** 🔧✅

Most issues can be resolved by following this guide systematically. Remember to check logs first, verify basic connectivity, and ensure proper permissions before diving into complex solutions.

**Need More Help?**
- Review the [Configuration Guide](configuration.md) for advanced options
- Check the [Deployment Guide](deployment.md) for production setups
- See the [Development Guide](development.md) for code-level debugging