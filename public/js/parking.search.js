(function () {
  function formatFeature(feature) {
    const featureMap = {
      covered: "Covered",
      "24/7": "24/7",
      ev_charging: "EV Charging",
      security: "Security",
      disabled: "Disabled Access",
    };
    return featureMap[feature] || feature;
  }

  function createParkingCard(recommendation, index) {
    const card = document.createElement("div");
    card.className = "parking-card bg-white rounded-lg shadow-md p-4";
    const lot = recommendation.parking_lot || {};

    card.innerHTML = `
      <div class="flex justify-between items-start mb-2">
        <h4 class="font-semibold text-gray-900">${lot.name || "Unknown"}</h4>
        ${index === 0 ? '<span class="recommendation-badge text-xs text-white px-2 py-1 rounded-full">Best Match</span>' : ""}
      </div>
      <div class="text-sm text-gray-600 mb-3">${lot.address || "N/A"}</div>
      <div class="grid grid-cols-2 gap-2 mb-3">
        <div class="text-sm"><span class="text-gray-500">Distance:</span> <span class="font-medium">${(lot.distance || 0).toFixed(2)} km</span></div>
        <div class="text-sm"><span class="text-gray-500">Price:</span> <span class="font-medium">¥${lot.price_per_hour || 0}/hr</span></div>
        <div class="text-sm"><span class="text-gray-500">Available:</span> <span class="font-medium">${lot.available_spaces || 0}/${lot.total_spaces || 0}</span></div>
        <div class="text-sm"><span class="text-gray-500">Score:</span> <span class="font-medium">${Math.round(recommendation.score || 0)}/100</span></div>
      </div>
      <div class="mb-3">
        <div class="text-sm text-gray-500 mb-1">Features:</div>
        <div class="flex flex-wrap gap-1">${(lot.features || []).map((f) => `<span class="text-xs bg-blue-100 text-blue-800 px-2 py-1 rounded">${formatFeature(f)}</span>`).join("")}</div>
      </div>
      <div class="flex space-x-2">
        <button onclick="showParkingDetails('${lot.id}')" class="flex-1 px-3 py-2 bg-blue-500 text-white text-sm rounded hover:bg-blue-600">View Details</button>
        <button onclick="reserveSpace('${lot.id}')" class="flex-1 px-3 py-2 bg-green-500 text-white text-sm rounded hover:bg-green-600">Reserve</button>
      </div>`;
    return card;
  }

  function addMapMarkers(recommendations) {
    const app = window.ParkingApp;
    app.markers.forEach((marker) => app.map.remove(marker));
    app.markers = [];

    recommendations.forEach((rec) => {
      const lot = rec.parking_lot;
      if (!lot) return;
      const marker = new AMap.Marker({
        position: [lot.longitude, lot.latitude],
        title: lot.name || "Unknown",
      });
      marker.on("click", function () {
        window.showParkingDetails(lot.id);
      });
      app.map.add(marker);
      app.markers.push(marker);
    });
  }

  window.getCurrentLocation = function () {
    if (!navigator.geolocation) return;
    navigator.geolocation.getCurrentPosition(
      function (position) {
        const lat = Math.round(position.coords.latitude * 1000) / 1000;
        const lng = Math.round(position.coords.longitude * 1000) / 1000;
        window.ParkingApp.userLocation = { lat: lat, lng: lng };
        document.getElementById("location-input").value = lat.toFixed(6) + ", " + lng.toFixed(6);
        if (window.ParkingApp.map) window.ParkingApp.map.setCenter([lng, lat]);
      },
      function () {}
    );
  };

  window.findParking = async function () {
    const app = window.ParkingApp;
    const maxPrice = parseFloat(document.getElementById("max-price").value) || 0;
    const maxDistance = parseFloat(document.getElementById("max-distance").value) || 5;

    document.getElementById("loading").classList.remove("hidden");
    document.getElementById("recommendations-section").classList.add("hidden");
    try {
      const data = await window.ParkingApi.findParking({
        latitude: app.userLocation.lat,
        longitude: app.userLocation.lng,
        max_price: maxPrice,
        max_distance: maxDistance,
        limit: 10,
      });
      const recommendations = data.recommendations || [];
      const container = document.getElementById("parking-recommendations");
      container.innerHTML = "";
      recommendations.forEach((rec, i) => container.appendChild(createParkingCard(rec, i)));
      document.getElementById("recommendations-section").classList.remove("hidden");
      addMapMarkers(recommendations);
    } catch (error) {
      window.showError("Failed to find parking spots: " + error.message);
    } finally {
      document.getElementById("loading").classList.add("hidden");
    }
  };

  window.showParkingDetails = async function (lotId) {
    try {
      const data = await window.ParkingApi.getLot(lotId);
      const lot = data.parking_lot || {};
      document.getElementById("modal-title").textContent = lot.name || "Parking Details";
      document.getElementById("modal-content").innerHTML = `
      <div class="space-y-4">
        <div><h4 class="font-medium text-gray-900">Address</h4><p class="text-gray-600">${lot.address || "N/A"}</p></div>
        <div class="flex space-x-3 pt-4 border-t">
          <button onclick="navigateToParking('${lot.id}')" class="flex-1 px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600">Navigate</button>
          <button onclick="reserveSpace('${lot.id}')" class="flex-1 px-4 py-2 bg-green-500 text-white rounded hover:bg-green-600">Reserve Space</button>
        </div>
      </div>`;
      document.getElementById("parking-modal").classList.remove("hidden");
    } catch (error) {
      window.showError("Failed to load parking details: " + error.message);
    }
  };

  window.reserveSpace = function (lotId) {
    const modal = document.getElementById("parking-modal");
    document.getElementById("modal-title").textContent = "Reserve Parking Space";
    document.getElementById("modal-content").innerHTML = `
      <form onsubmit="submitReservation(event, '${lotId}')" class="space-y-4">
        <div><label class="block text-sm font-medium text-gray-700">Start Time</label><input type="datetime-local" id="start-time" required class="w-full px-3 py-2 border border-gray-300 rounded-md"></div>
        <div><label class="block text-sm font-medium text-gray-700">End Time</label><input type="datetime-local" id="end-time" required class="w-full px-3 py-2 border border-gray-300 rounded-md"></div>
        <div class="flex space-x-3 pt-4 border-t">
          <button type="button" onclick="closeModal()" class="flex-1 px-4 py-2 border border-gray-300 text-gray-700 rounded hover:bg-gray-50">Cancel</button>
          <button type="submit" class="flex-1 px-4 py-2 bg-green-500 text-white rounded hover:bg-green-600">Confirm Reservation</button>
        </div>
      </form>`;
    modal.classList.remove("hidden");
  };

  window.submitReservation = async function (event, lotId) {
    event.preventDefault();
    try {
      const spacesData = await window.ParkingApi.getLotSpaces(lotId);
      const suitableSpace = (spacesData.parking_spaces || []).find((space) => space.is_available && !space.is_reserved);
      if (!suitableSpace) {
        window.showError("No suitable spaces available");
        return;
      }
      await window.ParkingApi.reserveSpace({
        parking_lot_id: lotId,
        space_id: suitableSpace.id,
        start_time: document.getElementById("start-time").value,
        end_time: document.getElementById("end-time").value,
      });
      window.showSuccess("Parking space reserved successfully");
      window.closeModal();
    } catch (error) {
      window.showError("Failed to reserve parking: " + error.message);
    }
  };

  window.navigateToParking = function () {
    alert("Navigation feature coming soon!");
  };
})();
