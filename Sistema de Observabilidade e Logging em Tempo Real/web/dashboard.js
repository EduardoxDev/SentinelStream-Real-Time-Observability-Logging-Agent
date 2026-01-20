let cpuChart, memoryChart, networkChart;
async function loadContainers() {
    try {
        const response = await fetch('/api/containers');
        const containers = await response.json();
        const select = document.getElementById('containerSelect');
        containers.forEach(container => {
            const option = document.createElement('option');
            option.value = container.ID;
            option.textContent = container.Names[0];
            select.appendChild(option);
        });
    } catch (error) {
        console.error('Failed to load containers:', error);
    }
}
async function loadData() {
    const containerID = document.getElementById('containerSelect').value;
    const timeRange = document.getElementById('timeRange').value;
    if (!containerID) {
        alert('Please select a container');
        return;
    }
    try {
        const response = await fetch(`/api/metrics?container_id=${containerID}&duration=${timeRange}`);
        const data = await response.json();
        updateCharts(data);
    } catch (error) {
        console.error('Failed to load metrics:', error);
    }
}
function updateCharts(data) {
    const timestamps = data.map(d => new Date(d.Timestamp).toLocaleTimeString());
    const cpuData = data.map(d => d.CPUPercent);
    const memoryData = data.map(d => d.MemoryPercent);
    const networkRxData = data.map(d => d.NetworkRx / 1024 / 1024); // MB
    const networkTxData = data.map(d => d.NetworkTx / 1024 / 1024); // MB
    if (cpuChart) cpuChart.destroy();
    cpuChart = new Chart(document.getElementById('cpuChart'), {
        type: 'line',
        data: {
            labels: timestamps,
            datasets: [{
                label: 'CPU %',
                data: cpuData,
                borderColor: '#3b82f6',
                backgroundColor: 'rgba(59, 130, 246, 0.1)',
                tension: 0.4
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            plugins: {
                legend: { labels: { color: '#e2e8f0' } }
            },
            scales: {
                y: {
                    beginAtZero: true,
                    max: 100,
                    ticks: { color: '#94a3b8' },
                    grid: { color: '#334155' }
                },
                x: {
                    ticks: { color: '#94a3b8' },
                    grid: { color: '#334155' }
                }
            }
        }
    });
    if (memoryChart) memoryChart.destroy();
    memoryChart = new Chart(document.getElementById('memoryChart'), {
        type: 'line',
        data: {
            labels: timestamps,
            datasets: [{
                label: 'Memory %',
                data: memoryData,
                borderColor: '#10b981',
                backgroundColor: 'rgba(16, 185, 129, 0.1)',
                tension: 0.4
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            plugins: {
                legend: { labels: { color: '#e2e8f0' } }
            },
            scales: {
                y: {
                    beginAtZero: true,
                    max: 100,
                    ticks: { color: '#94a3b8' },
                    grid: { color: '#334155' }
                },
                x: {
                    ticks: { color: '#94a3b8' },
                    grid: { color: '#334155' }
                }
            }
        }
    });
    if (networkChart) networkChart.destroy();
    networkChart = new Chart(document.getElementById('networkChart'), {
        type: 'line',
        data: {
            labels: timestamps,
            datasets: [
                {
                    label: 'RX (MB)',
                    data: networkRxData,
                    borderColor: '#8b5cf6',
                    backgroundColor: 'rgba(139, 92, 246, 0.1)',
                    tension: 0.4
                },
                {
                    label: 'TX (MB)',
                    data: networkTxData,
                    borderColor: '#ec4899',
                    backgroundColor: 'rgba(236, 72, 153, 0.1)',
                    tension: 0.4
                }
            ]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            plugins: {
                legend: { labels: { color: '#e2e8f0' } }
            },
            scales: {
                y: {
                    beginAtZero: true,
                    ticks: { color: '#94a3b8' },
                    grid: { color: '#334155' }
                },
                x: {
                    ticks: { color: '#94a3b8' },
                    grid: { color: '#334155' }
                }
            }
        }
    });
}
loadContainers();