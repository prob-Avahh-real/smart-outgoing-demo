(function () {
  window.ParkingApp = {
    map: null,
    userLocation: { lat: 22.6913, lng: 114.0448 },
    markers: [],
    avpStatusPoller: null,
  };

  window.showError = function (message) {
    alert("Error: " + message);
  };

  window.showSuccess = function (message) {
    alert("Success: " + message);
  };

  window.closeModal = function () {
    const modal = document.getElementById("parking-modal");
    if (modal) modal.classList.add("hidden");
  };

  function initMap() {
    if (typeof AMap === "undefined") {
      setTimeout(initMap, 500);
      return;
    }
    try {
      const app = window.ParkingApp;
      app.map = new AMap.Map("map", {
        zoom: 13,
        center: [app.userLocation.lng, app.userLocation.lat],
        viewMode: "2D",
        resizeEnable: true,
      });
      app.map.addControl(new AMap.Scale());
      app.map.addControl(new AMap.ToolBar());
    } catch (error) {
      console.error("Error initializing map:", error);
    }
  }

  window.onload = function () {
    initMap();
    window.getCurrentLocation();
  };
})();
