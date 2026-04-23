// Application State Management
class AppState {
  constructor() {
    this.config = null;
    this.map = null;
    this.vehicles = new Map();
    this.markers = new Map();
    this.routeLines = new Map();
    this.trailLines = new Map();
    this.destinationMarkers = new Map();
    this.sim = new Map();
    this.socket = null;
    this.pollingTimer = null;
    this.reconnectTimer = null;
    this.adminToken = "";
    this.playing = false;
    this.clockSec = 0;
    this.timer = null;
    this.viewMode = '2D'; // '2D' or '3D'
  }
}

// API Service
class ApiService {
  constructor(baseUrl = '') {
    this.baseUrl = baseUrl;
  }

  async request(endpoint, options = {}) {
    const url = `${this.baseUrl}${endpoint}`;
    const config = {
      headers: {
        'Content-Type': 'application/json',
        ...options.headers,
      },
      ...options,
    };

    const response = await fetch(url, config);
    
    if (!response.ok) {
      throw new Error(`HTTP ${response.status}: ${response.statusText}`);
    }

    return response.json();
  }

  async getConfig() {
    return this.request('/api/config');
  }

  async getVehicles() {
    return this.request('/api/vehicles');
  }

  async createVehicle(vehicleData) {
    return this.request('/api/vehicles', {
      method: 'POST',
      body: JSON.stringify(vehicleData),
    });
  }


  async getCacheStats() {
    return this.request('/api/cache/stats');
  }

  async cleanupCache() {
    return this.request('/api/cache/cleanup', { method: 'POST' });
  }

  async getExperimentStats(experimentId) {
    return this.request(`/api/ab/experiments/${experimentId}/stats`);
  }
}

// WebSocket Service
class WebSocketService {
  constructor(url) {
    this.url = url;
    this.socket = null;
    this.reconnectAttempts = 0;
    this.maxReconnectAttempts = 5;
    this.reconnectDelay = 2500;
  }

  connect(onMessage, onError, onClose) {
    try {
      this.socket = new WebSocket(this.url);

      this.socket.addEventListener('open', () => {
        console.log('WebSocket connected');
        this.reconnectAttempts = 0;
        onMessage({ type: 'connection', status: 'connected' });
      });

      this.socket.addEventListener('message', (event) => {
        try {
          const payload = JSON.parse(event.data);
          onMessage(payload);
        } catch (error) {
          console.error('WebSocket message parse error:', error);
        }
      });

      this.socket.addEventListener('close', () => {
        console.log('WebSocket disconnected');
        this.scheduleReconnect(onMessage, onError, onClose);
        onClose();
      });

      this.socket.addEventListener('error', (error) => {
        console.error('WebSocket error:', error);
        onError(error);
      });

    } catch (error) {
      console.error('WebSocket connection error:', error);
      onError(error);
    }
  }

  scheduleReconnect(onMessage, onError, onClose) {
    if (this.reconnectAttempts >= this.maxReconnectAttempts) {
      console.log('Max reconnect attempts reached');
      return;
    }

    this.reconnectAttempts++;
    console.log(`Scheduling reconnect attempt ${this.reconnectAttempts}`);

    setTimeout(() => {
      this.connect(onMessage, onError, onClose);
    }, this.reconnectDelay * this.reconnectAttempts);
  }

  disconnect() {
    if (this.socket) {
      this.socket.close();
      this.socket = null;
    }
  }
}

// Map Service
class MapService {
  constructor(containerId, config) {
    this.containerId = containerId;
    this.config = config;
    this.map = null;
    this.markers = new Map();
    this.routeLines = new Map();
    this.trailLines = new Map();
    this.destinationMarkers = new Map();
  }

  async initialize() {
    if (!this.config.amap_js_key) {
      throw new Error('AMap JS key is required');
    }

    return new Promise((resolve, reject) => {
      // Load AMap script
      const script = document.createElement('script');
      script.src = `https://webapi.amap.com/maps?v=2.0&key=${this.config.amap_js_key}`;
      script.onload = () => {
        this.createMap();
        resolve();
      };
      script.onerror = () => {
        reject(new Error('Failed to load AMap script'));
      };
      document.head.appendChild(script);
    });
  }

  createMap() {
    this.map = new AMap.Map(this.containerId, {
      zoom: 15,
      center: this.config.default_center,
      viewMode: '2D',
    });
  }

  addVehicleMarker(vehicle) {
    const marker = new AMap.Marker({
      position: vehicle.position || [vehicle.start_lng, vehicle.start_lat],
      content: '<div style="font-size:20px;">\ud83d\ude97</div>',
    });

    this.markers.set(vehicle.id, marker);
    marker.setMap(this.map);
  }

  updateVehicleMarker(vehicle, isFinished = false) {
    const marker = this.markers.get(vehicle.id);
    if (!marker) return;

    const icon = isFinished ? '\u26aa' : '\ud83d\ude97';
    marker.setContent(`<div style="font-size:20px;">${icon}</div>`);

    if (vehicle.position) {
      marker.setPosition(vehicle.position);
    }
  }

  addRouteLine(vehicleId, route) {
    if (!route || route.length < 2) return;

    const line = new AMap.Polyline({
      path: route,
      strokeColor: '#e94560',
      strokeWeight: 3,
    });

    this.routeLines.set(vehicleId, line);
    line.setMap(this.map);
  }

  updateTrailLine(vehicleId, trail) {
    if (!trail || trail.length < 2) return;

    let line = this.trailLines.get(vehicleId);
    if (!line) {
      line = new AMap.Polyline({
        path: trail,
        strokeColor: '#666666',
        strokeWeight: 2,
      });
      this.trailLines.set(vehicleId, line);
    }

    line.setPath(trail);
    line.setMap(this.map);
  }

  removeVehicle(vehicleId) {
    const marker = this.markers.get(vehicleId);
    if (marker) {
      marker.setMap(null);
      this.markers.delete(vehicleId);
    }

    const routeLine = this.routeLines.get(vehicleId);
    if (routeLine) {
      routeLine.setMap(null);
      this.routeLines.delete(vehicleId);
    }

    const trailLine = this.trailLines.get(vehicleId);
    if (trailLine) {
      trailLine.setMap(null);
      this.trailLines.delete(vehicleId);
    }

    const destMarker = this.destinationMarkers.get(vehicleId);
    if (destMarker) {
      destMarker.setMap(null);
      this.destinationMarkers.delete(vehicleId);
    }
  }

  cleanup() {
    this.markers.forEach(marker => marker.setMap(null));
    this.routeLines.forEach(line => line.setMap(null));
    this.trailLines.forEach(line => line.setMap(null));
    this.destinationMarkers.forEach(marker => marker.setMap(null));
    
    this.markers.clear();
    this.routeLines.clear();
    this.trailLines.clear();
    this.destinationMarkers.clear();
  }
}

// Simulation Service
class SimulationService {
  constructor(mapService, state) {
    this.mapService = mapService;
    this.state = state;
    this.frameInterval = 100; // 100ms per frame
  }

  start() {
    if (this.state.playing) return;
    
    this.state.playing = true;
    this.runFrame();
  }

  pause() {
    this.state.playing = false;
    if (this.state.timer) {
      clearTimeout(this.state.timer);
      this.state.timer = null;
    }
  }

  reset() {
    this.pause();
    this.state.clockSec = 0;
    this.updateClockDisplay();

    // Reset all vehicles to start positions
    this.state.vehicles.forEach((vehicle, id) => {
      const sim = this.state.sim.get(id);
      if (sim) {
        sim.progress = 0;
        sim.finished = false;
        sim.trail = [];
        
        // Reset marker position
        this.mapService.updateVehicleMarker(vehicle, false);
        
        // Clear trail
        const trailLine = this.mapService.trailLines.get(id);
        if (trailLine) {
          trailLine.setPath([]);
        }
      }
    });

    this.updateStatus();
  }

  runFrame() {
    if (!this.state.playing) return;

    this.state.clockSec += 0.1;
    this.updateClockDisplay();

    let doneCount = 0;
    const totalCount = this.state.vehicles.size;

    this.state.vehicles.forEach((vehicle, id) => {
      const sim = this.state.sim.get(id);
      if (!sim) return;

      const route = Array.isArray(vehicle.route) ? vehicle.route : [];
      if (route.length < 2) {
        sim.finished = true;
        doneCount++;
        return;
      }

      if (!sim.finished) {
        sim.progress += 0.08;
        if (sim.progress >= route.length - 1) {
          sim.progress = route.length - 1;
          sim.finished = true;
        }
      }

      const position = this.getPositionAtProgress(route, sim.progress);
      if (position) {
        // Update marker
        vehicle.position = position;
        this.mapService.updateVehicleMarker(vehicle, sim.finished);

        // Update trail
        sim.trail.push(position);
        if (sim.trail.length > 80) {
          sim.trail.shift();
        }
        
        this.mapService.updateTrailLine(id, sim.trail);
      }

      if (sim.finished) {
        doneCount++;
      }
    });

    this.updateStatus(doneCount, totalCount);
    this.state.timer = setTimeout(() => this.runFrame(), this.frameInterval);
  }

  getPositionAtProgress(route, progress) {
    if (!route || route.length < 2) return null;
    
    const index = Math.floor(progress);
    const fraction = progress - index;
    
    if (index >= route.length - 1) {
      return route[route.length - 1];
    }
    
    const start = route[index];
    const end = route[index + 1];
    
    return [
      start[0] + (end[0] - start[0]) * fraction,
      start[1] + (end[1] - start[1]) * fraction,
    ];
  }

  updateClockDisplay() {
    const element = document.getElementById('time');
    if (element) {
      element.textContent = this.state.clockSec.toFixed(1);
    }
  }

  updateStatus(done = null, total = null) {
    const doneElement = document.getElementById('done');
    const totalElement = document.getElementById('total');
    const slotElement = document.getElementById('slot');

    if (done !== null && doneElement) {
      doneElement.textContent = done;
    }
    if (total !== null && totalElement) {
      totalElement.textContent = total;
    }

    if (slotElement) {
      if (this.state.playing) {
        slotElement.textContent = 'RUNNING';
        slotElement.style.background = '#10b981';
      } else if (done === total && total > 0) {
        slotElement.textContent = 'COMPLETED';
        slotElement.style.background = '#3b82f6';
      } else {
        slotElement.textContent = 'WAITING';
        slotElement.style.background = '#e94560';
      }
    }
  }
}

// UI Controller
class UIController {
  constructor(state, apiService, mapService, simulationService) {
    this.state = state;
    this.apiService = apiService;
    this.mapService = mapService;
    this.simulationService = simulationService;
  }

  initialize() {
    this.bindEvents();
    this.loadAdminToken();
  }

  bindEvents() {
    // Control buttons
    document.getElementById('save-token-btn')?.addEventListener('click', () => {
      this.saveAdminToken();
    });

    document.getElementById('spawn-btn')?.addEventListener('click', () => {
      this.createDemoVehicle();
    });

    // Simulation controls
    window.start = () => this.simulationService.start();
    window.pause = () => this.simulationService.pause();
    window.reset = () => this.simulationService.reset();
    window.toggleViewMode = () => this.toggleViewMode();
  }

  async saveAdminToken() {
    const tokenInput = document.getElementById('admin-token');
    if (tokenInput) {
      this.state.adminToken = tokenInput.value;
      localStorage.setItem('adminToken', this.state.adminToken);
      alert('Token saved');
    }
  }

  loadAdminToken() {
    const savedToken = localStorage.getItem('adminToken');
    if (savedToken) {
      this.state.adminToken = savedToken;
      const tokenInput = document.getElementById('admin-token');
      if (tokenInput) {
        tokenInput.value = savedToken;
      }
    }
  }

  async createDemoVehicle() {
    try {
      const vehicleData = {
        name: `Demo Vehicle-${Math.floor(Math.random() * 90 + 10)}`,
        start_lng: 114.0448 + Math.random() * 0.01,
        start_lat: 22.6913 + Math.random() * 0.01,
        start_alt: Math.random() * 100,
      };

      await this.apiService.createVehicle(vehicleData);
      await this.refreshVehicles();
    } catch (error) {
      alert(`Failed to create vehicle: ${error.message}`);
    }
  }


  async refreshVehicles() {
    try {
      const data = await this.apiService.getVehicles();
      this.applySnapshot(data.items || []);
    } catch (error) {
      console.error('Failed to refresh vehicles:', error);
    }
  }

  applySnapshot(vehicles) {
    const validIds = new Set();

    vehicles.forEach((vehicle) => {
      validIds.add(vehicle.id);
      this.state.vehicles.set(vehicle.id, vehicle);
      
      // Ensure vehicle overlay exists
      this.ensureVehicleOverlay(vehicle);
      
      // Initialize simulation state for new vehicles
      if (!this.state.sim.has(vehicle.id)) {
        this.state.sim.set(vehicle.id, {
          progress: 0,
          finished: false,
          trail: [],
        });
      }
    });

    this.removeMissingOverlays(validIds);
    this.simulationService.updateStatus();
  }

  ensureVehicleOverlay(vehicle) {
    if (!this.mapService.markers.has(vehicle.id)) {
      this.mapService.addVehicleMarker(vehicle);
    }

    if (vehicle.route && vehicle.route.length > 0) {
      this.mapService.addRouteLine(vehicle.id, vehicle.route);
    }

    if (!this.state.destinationMarkers.has(vehicle.id)) {
      const destination = vehicle.destination || vehicle.route?.[vehicle.route.length - 1];
      if (destination) {
        const marker = new AMap.Marker({
          position: destination,
          content: '<div style="font-size:16px;">\ud83d\udccd</div>',
        });
        this.state.destinationMarkers.set(vehicle.id, marker);
        marker.setMap(this.mapService.map);
      }
    }
  }

  removeMissingOverlays(validIds) {
    this.state.vehicles.forEach((vehicle, id) => {
      if (!validIds.has(id)) {
        this.mapService.removeVehicle(id);
        this.state.sim.delete(id);
      }
    });
  }

  toggleViewMode() {
    const viewToggle = document.getElementById('view-toggle');
    const viewIcon = document.getElementById('view-icon');
    const viewText = document.getElementById('view-text');

    if (this.state.viewMode === '2D') {
      this.state.viewMode = '3D';
      if (viewIcon) viewIcon.textContent = '\ud83c\udfd4\ufe0f';
      if (viewText) viewText.textContent = '3D\u89c6\u56fe';
      if (viewToggle) viewToggle.classList.add('active');
    } else {
      this.state.viewMode = '2D';
      if (viewIcon) viewIcon.textContent = '\ud83d\uddfa\ufe0f';
      if (viewText) viewText.textContent = '2D\u89c6\u56fe';
      if (viewToggle) viewToggle.classList.remove('active');
    }

    // Update map view mode if supported
    if (this.mapService.map) {
      this.mapService.map.setViewMode(this.state.viewMode);
    }
  }

  updateSyncStatus(status) {
    const syncElement = document.getElementById('sync');
    if (syncElement) {
      syncElement.textContent = status;
    }
  }
}

// Main Application
class App {
  constructor() {
    this.state = new AppState();
    this.apiService = new ApiService();
    this.mapService = null;
    this.simulationService = null;
    this.uiController = null;
    this.wsService = null;
  }

  async initialize() {
    try {
      // Load configuration
      this.state.config = await this.apiService.getConfig();
      
      // Initialize map service
      this.mapService = new MapService('map', this.state.config);
      await this.mapService.initialize();
      
      // Initialize simulation service
      this.simulationService = new SimulationService(this.mapService, this.state);
      
      // Initialize UI controller
      this.uiController = new UIController(this.state, this.apiService, this.mapService, this.simulationService);
      this.uiController.initialize();
      
      // Initialize WebSocket service
      this.wsService = new WebSocketService(`ws://${window.location.host}/ws`);
      this.wsService.connect(
        (message) => this.handleWebSocketMessage(message),
        (error) => this.handleWebSocketError(error),
        () => this.handleWebSocketClose()
      );
      
      // Load initial data
      await this.refreshVehicles();
      
      console.log('Application initialized successfully');
      
    } catch (error) {
      console.error('Failed to initialize application:', error);
      alert(`Initialization failed: ${error.message}`);
    }
  }

  async refreshVehicles() {
    return this.uiController?.refreshVehicles();
  }

  handleWebSocketMessage(message) {
    if (Array.isArray(message.items)) {
      this.uiController?.applySnapshot(message.items);
    }
  }

  handleWebSocketError(error) {
    console.error('WebSocket error:', error);
    this.uiController?.updateSyncStatus('Connection Error');
  }

  handleWebSocketClose() {
    this.uiController?.updateSyncStatus('Reconnecting...');
    
    // Start polling as fallback
    if (!this.state.pollingTimer) {
      this.state.pollingTimer = setInterval(() => {
        this.refreshVehicles().catch(console.error);
      }, 2000);
    }
  }
}

// Initialize application when DOM is ready
document.addEventListener('DOMContentLoaded', () => {
  const app = new App();
  app.initialize().catch(console.error);
});
