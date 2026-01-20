let ws;
let reconnectInterval;
const containers = new Map();
let metricsHistory = {
    cpu: [],
    memory: [],
    network: []
};
let activityLog = [];
document.querySelectorAll('.nav-link').forEach(link => {
    link.addEventListener('click', function() {
        const tabName = this.getAttribute('data-tab');
        switchTab(tabName);
    });
});
function switchTab(tabName) {
    document.querySelectorAll('.nav-link').forEach(link => {
        link.classList.remove('active');
    });
    document.querySelector(`[data-tab="${tabName}"]`).classList.add('active');
    document.querySelectorAll('.tab-content').forEach(content => {
        content.classList.remove('active');
    });
    document.getElementById(`${tabName}-tab`).classList.add('active');
    const titles = {
        'overview': { title: 'Overview', subtitle: 'Real-time system monitoring dashboard' },
        'processes': { title: 'Processes', subtitle: 'Monitor all running processes' },
        'metrics': { title: 'System Metrics', subtitle: 'Detailed system performance metrics' },
        'network': { title: 'Network Traffic', subtitle: 'Network usage and bandwidth monitoring' },
        'alerts': { title: 'Alerts & Notifications', subtitle: 'System alerts and warnings' },
        'activity': { title: 'Activity Log', subtitle: 'System activity and event history' },
        'settings': { title: 'Settings', subtitle: 'Configure dashboard preferences' }
    };
    document.getElementById('page-title').textContent = titles[tabName].title;
    document.getElementById('page-subtitle').textContent = titles[tabName].subtitle;
}
function connect() {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsUrl = `${protocol}//${window.location.host}/ws`;
    ws = new WebSocket(wsUrl);
    ws.onopen = () => {
        console.log('✅ WebSocket connected');
        updateStatus(true);
        clearInterval(reconnectInterval);
        addActivityLog('success', 'WebSocket Connected', 'Successfully connected to monitoring server');
    };
    ws.onmessage = (event) => {
        const data = JSON.parse(event.data);
        updateDashboard(data);
    };
    ws.onerror = (error) => {
        console.error('❌ WebSocket error:', error);
        addActivityLog('error', 'Connection Error', 'Failed to connect to monitoring server');
    };
    ws.onclose = () => {
        console.log('🔌 WebSocket disconnected');
        updateStatus(false);
        reconnectInterval = setInterval(connect, 5000);
        addActivityLog('warning', 'WebSocket Disconnected', 'Connection lost, attempting to reconnect...');
    };
}
function updateStatus(connected) {
    const statusEl = document.getElementById('status');
    if (connected) {
        statusEl.className = 'status-badge connected';
        statusEl.innerHTML = '<i class="fas fa-circle"></i><span>CONNECTED</span>';
    } else {
        statusEl.className = 'status-badge disconnected';
        statusEl.innerHTML = '<i class="fas fa-circle"></i><span>DISCONNECTED</span>';
    }
}
function updateDashboard(metrics) {
    if (!metrics || metrics.length === 0) {
        showEmptyState();
        return;
    }
    let totalCpu = 0;
    let totalMemory = 0;
    let totalAlerts = 0;
    let totalRx = 0;
    let totalTx = 0;
    let maxCpu = 0;
    let maxMemory = 0;
    metrics.forEach(metric => {
        totalCpu += metric.CPUPercent;
        totalMemory += metric.MemoryPercent;
        totalRx += metric.NetworkRx;
        totalTx += metric.NetworkTx;
        if (metric.CPUPercent > maxCpu) maxCpu = metric.CPUPercent;
        if (metric.MemoryPercent > maxMemory) maxMemory = metric.MemoryPercent;
        if (metric.CPUPercent > 90 || metric.MemoryPercent > 85) {
            totalAlerts++;
        }
        let card = containers.get(metric.ContainerID);
        if (!card) {
            card = createContainerCard(metric);
            containers.set(metric.ContainerID, card);
            document.getElementById('containers').appendChild(card.element);
        }
        updateContainerCard(card, metric);
    });
    const avgCpu = totalCpu / metrics.length;
    const avgMemory = totalMemory / metrics.length;
    const avgNetwork = ((totalRx + totalTx) / metrics.length / 1024 / 1024).toFixed(2);
    document.getElementById('total-processes').textContent = metrics.length;
    document.getElementById('avg-cpu').textContent = avgCpu.toFixed(1) + '%';
    document.getElementById('avg-memory').textContent = avgMemory.toFixed(1) + '%';
    document.getElementById('total-alerts').textContent = totalAlerts;
    updateCircularProgress('cpu-circle', avgCpu);
    updateCircularProgress('memory-circle', avgMemory);
    updateCircularProgress('network-circle', Math.min(avgNetwork * 10, 100));
    document.getElementById('metrics-cpu-total').textContent = avgCpu.toFixed(1) + '%';
    document.getElementById('metrics-cpu-bar').style.width = avgCpu + '%';
    document.getElementById('metrics-cpu-peak').textContent = maxCpu.toFixed(1) + '%';
    document.getElementById('metrics-cpu-avg').textContent = avgCpu.toFixed(1) + '%';
    document.getElementById('metrics-mem-total').textContent = avgMemory.toFixed(1) + '%';
    document.getElementById('metrics-mem-bar').style.width = avgMemory + '%';
    document.getElementById('metrics-mem-peak').textContent = maxMemory.toFixed(1) + '%';
    document.getElementById('metrics-mem-avg').textContent = avgMemory.toFixed(1) + '%';
    document.getElementById('metrics-proc-active').textContent = metrics.length;
    document.getElementById('metrics-proc-running').textContent = metrics.length;
    document.getElementById('metrics-proc-idle').textContent = '0';
    document.getElementById('net-total-rx').textContent = formatBytes(totalRx);
    document.getElementById('net-total-tx').textContent = formatBytes(totalTx);
    document.getElementById('net-total').textContent = formatBytes(totalRx + totalTx);
    document.getElementById('net-bandwidth').textContent = avgNetwork + ' MB/s';
    updateNetworkTable(metrics);
    updateAlertsList(metrics);
    updateResourceAllocation(metrics);
}
function updateCircularProgress(id, percentage) {
    const circle = document.getElementById(id);
    if (!circle) return;
    const progressCircle = circle.querySelector('.progress');
    const percentageText = circle.querySelector('.percentage');
    const circumference = 326.73;
    const offset = circumference - (percentage / 100) * circumference;
    progressCircle.style.strokeDashoffset = offset;
    percentageText.textContent = Math.round(percentage) + '%';
}
function createContainerCard(metric) {
    const card = document.createElement('div');
    card.className = 'container-card';
    card.innerHTML = `
        <div class="container-header">
            <div class="container-name">
                <i class="fas fa-cube"></i>
                <span>${metric.ContainerName}</span>
            </div>
            <div class="container-id">${metric.ContainerID.substring(8, 20)}</div>
        </div>
        <div class="metric">
            <div class="metric-header">
                <div class="metric-label"><i class="fas fa-microchip"></i>CPU</div>
                <div class="metric-value cpu-value">0%</div>
            </div>
            <div class="metric-bar">
                <div class="metric-fill cpu-fill" style="width: 0%"></div>
            </div>
        </div>
        <div class="metric">
            <div class="metric-header">
                <div class="metric-label"><i class="fas fa-memory"></i>Memory</div>
                <div class="metric-value memory-value">0%</div>
            </div>
            <div class="metric-bar">
                <div class="metric-fill memory-fill" style="width: 0%"></div>
            </div>
        </div>
        <div class="network-info">
            <div class="network-stat">
                <div class="label"><i class="fas fa-download"></i>Download</div>
                <div class="value rx-value">0 B</div>
            </div>
            <div class="network-stat">
                <div class="label"><i class="fas fa-upload"></i>Upload</div>
                <div class="value tx-value">0 B</div>
            </div>
        </div>
        <div class="alert-container"></div>
    `;
    return { element: card };
}
function updateContainerCard(card, metric) {
    const cpuValue = card.element.querySelector('.cpu-value');
    const cpuFill = card.element.querySelector('.cpu-fill');
    const memoryValue = card.element.querySelector('.memory-value');
    const memoryFill = card.element.querySelector('.memory-fill');
    const rxValue = card.element.querySelector('.rx-value');
    const txValue = card.element.querySelector('.tx-value');
    const alertContainer = card.element.querySelector('.alert-container');
    cpuValue.textContent = `${metric.CPUPercent.toFixed(1)}%`;
    cpuFill.style.width = `${Math.min(metric.CPUPercent, 100)}%`;
    memoryValue.textContent = `${metric.MemoryPercent.toFixed(1)}%`;
    memoryFill.style.width = `${Math.min(metric.MemoryPercent, 100)}%`;
    rxValue.textContent = formatBytes(metric.NetworkRx);
    txValue.textContent = formatBytes(metric.NetworkTx);
    alertContainer.innerHTML = '';
    if (metric.CPUPercent > 90) {
        alertContainer.innerHTML += `
            <div class="alert">
                <i class="fas fa-exclamation-triangle"></i>
                <span>HIGH CPU USAGE</span>
            </div>
        `;
        addActivityLog('warning', 'High CPU Alert', `${metric.ContainerName} CPU usage at ${metric.CPUPercent.toFixed(1)}%`);
    }
    if (metric.MemoryPercent > 85) {
        alertContainer.innerHTML += `
            <div class="alert">
                <i class="fas fa-exclamation-triangle"></i>
                <span>HIGH MEMORY USAGE</span>
            </div>
        `;
        addActivityLog('warning', 'High Memory Alert', `${metric.ContainerName} memory usage at ${metric.MemoryPercent.toFixed(1)}%`);
    }
}
function updateNetworkTable(metrics) {
    const tbody = document.getElementById('network-table');
    tbody.innerHTML = '';
    metrics.forEach(metric => {
        const row = document.createElement('tr');
        const total = metric.NetworkRx + metric.NetworkTx;
        const status = total > 1024 * 1024 * 100 ? 'warning' : 'success';
        row.innerHTML = `
            <td><i class="fas fa-cube" style="color: #22c55e; margin-right: 8px;"></i>${metric.ContainerName}</td>
            <td>${formatBytes(metric.NetworkRx)}</td>
            <td>${formatBytes(metric.NetworkTx)}</td>
            <td>${formatBytes(total)}</td>
            <td><span class="badge ${status}">${status === 'success' ? 'Normal' : 'High'}</span></td>
        `;
        tbody.appendChild(row);
    });
}
function updateAlertsList(metrics) {
    const alertsList = document.getElementById('alerts-list');
    alertsList.innerHTML = '';
    let hasAlerts = false;
    metrics.forEach(metric => {
        if (metric.CPUPercent > 90) {
            hasAlerts = true;
            alertsList.innerHTML += `
                <div class="activity-item">
                    <div class="activity-icon error">
                        <i class="fas fa-exclamation-triangle"></i>
                    </div>
                    <div class="activity-content">
                        <div class="activity-title">High CPU Usage Alert</div>
                        <div class="activity-desc">${metric.ContainerName} is using ${metric.CPUPercent.toFixed(1)}% CPU</div>
                        <div class="activity-time">${new Date().toLocaleTimeString()}</div>
                    </div>
                </div>
            `;
        }
        if (metric.MemoryPercent > 85) {
            hasAlerts = true;
            alertsList.innerHTML += `
                <div class="activity-item">
                    <div class="activity-icon warning">
                        <i class="fas fa-exclamation-triangle"></i>
                    </div>
                    <div class="activity-content">
                        <div class="activity-title">High Memory Usage Alert</div>
                        <div class="activity-desc">${metric.ContainerName} is using ${metric.MemoryPercent.toFixed(1)}% memory</div>
                        <div class="activity-time">${new Date().toLocaleTimeString()}</div>
                    </div>
                </div>
            `;
        }
    });
    if (!hasAlerts) {
        alertsList.innerHTML = `
            <div class="empty-state">
                <i class="fas fa-check-circle"></i>
                <h3>No active alerts</h3>
                <p>All systems are operating normally</p>
            </div>
        `;
    }
}
function updateResourceAllocation(metrics) {
    const container = document.getElementById('resource-allocation');
    container.innerHTML = '';
    metrics.slice(0, 5).forEach(metric => {
        container.innerHTML += `
            <div class="metric">
                <div class="metric-header">
                    <div class="metric-label"><i class="fas fa-cube"></i>${metric.ContainerName}</div>
                    <div class="metric-value">${metric.CPUPercent.toFixed(1)}%</div>
                </div>
                <div class="metric-bar">
                    <div class="metric-fill cpu-fill" style="width: ${Math.min(metric.CPUPercent, 100)}%"></div>
                </div>
            </div>
        `;
    });
}
function addActivityLog(type, title, description) {
    const log = {
        type,
        title,
        description,
        time: new Date().toLocaleTimeString()
    };
    activityLog.unshift(log);
    if (activityLog.length > 50) activityLog.pop();
    updateActivityLog();
    updateRecentActivity();
}
function updateActivityLog() {
    const container = document.getElementById('activity-log');
    container.innerHTML = '';
    activityLog.forEach(log => {
        container.innerHTML += `
            <div class="activity-item">
                <div class="activity-icon ${log.type}">
                    <i class="fas fa-${getIconForType(log.type)}"></i>
                </div>
                <div class="activity-content">
                    <div class="activity-title">${log.title}</div>
                    <div class="activity-desc">${log.description}</div>
                    <div class="activity-time">${log.time}</div>
                </div>
            </div>
        `;
    });
}
function updateRecentActivity() {
    const container = document.getElementById('recent-activity');
    container.innerHTML = '';
    activityLog.slice(0, 5).forEach(log => {
        container.innerHTML += `
            <div class="activity-item">
                <div class="activity-icon ${log.type}">
                    <i class="fas fa-${getIconForType(log.type)}"></i>
                </div>
                <div class="activity-content">
                    <div class="activity-title">${log.title}</div>
                    <div class="activity-desc">${log.description}</div>
                    <div class="activity-time">${log.time}</div>
                </div>
            </div>
        `;
    });
}
function getIconForType(type) {
    const icons = {
        'success': 'check-circle',
        'error': 'times-circle',
        'warning': 'exclamation-triangle',
        'info': 'info-circle'
    };
    return icons[type] || 'info-circle';
}
function showEmptyState() {
    const containerGrid = document.getElementById('containers');
    containerGrid.innerHTML = `
        <div class="empty-state">
            <i class="fas fa-server"></i>
            <h3>No processes detected</h3>
            <p>Waiting for metrics...</p>
        </div>
    `;
}
function formatBytes(bytes) {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return (bytes / Math.pow(k, i)).toFixed(1) + ' ' + sizes[i];
}
connect();
addActivityLog('info', 'Dashboard Initialized', 'Observability dashboard started successfully');