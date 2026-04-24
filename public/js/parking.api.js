(function () {
  async function requestJSON(url, options) {
    const response = await fetch(url, options || {});
    let payload = {};
    try {
      payload = await response.json();
    } catch (_) {}
    if (!response.ok) {
      throw new Error(payload.error || ("HTTP " + response.status));
    }
    return payload;
  }

  window.ParkingApi = {
    findParking(body) {
      return requestJSON("/api/parking/find", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(body),
      });
    },
    getLot(lotId) {
      return requestJSON("/api/parking/lots/" + encodeURIComponent(lotId));
    },
    getLotSpaces(lotId) {
      return requestJSON("/api/parking/lots/" + encodeURIComponent(lotId) + "/spaces");
    },
    reserveSpace(body) {
      return requestJSON("/api/parking/reserve", {
        method: "POST",
        headers: { "Content-Type": "application/json", "x-user-id": "demo_user" },
        body: JSON.stringify(body),
      });
    },
    startAVP(body) {
      return requestJSON("/api/parking/avp/start", {
        method: "POST",
        headers: { "Content-Type": "application/json", "x-user-id": "demo_user" },
        body: JSON.stringify(body),
      });
    },
    summonAVP(body) {
      return requestJSON("/api/parking/avp/summon", {
        method: "POST",
        headers: { "Content-Type": "application/json", "x-user-id": "demo_user" },
        body: JSON.stringify(body),
      });
    },
    getAVPTask(taskId) {
      return requestJSON("/api/parking/avp/tasks/" + encodeURIComponent(taskId));
    },
    cancelAVP(taskId) {
      return requestJSON("/api/parking/avp/tasks/" + encodeURIComponent(taskId) + "/cancel", {
        method: "POST",
      });
    },
  };
})();
