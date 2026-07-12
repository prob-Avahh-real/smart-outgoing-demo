(function () {
  function getInputValue(id) {
    const element = document.getElementById(id);
    return element ? element.value.trim() : "";
  }

  function setAVPStatus(content, isError) {
    const statusContainer = document.getElementById("avp-status");
    if (!statusContainer) return;
    statusContainer.className = "text-sm " + (isError ? "text-red-700" : "text-gray-700");
    statusContainer.innerHTML = content;
  }

  function renderTask(task, session) {
    if (!task) return;
    const sessionId = (session && session.id) || task.session_id || "-";
    setAVPStatus(
      "<div><strong>Task ID:</strong> " + (task.id || "-") + "</div>" +
        "<div><strong>Session ID:</strong> " + sessionId + "</div>" +
        "<div><strong>Type:</strong> " + (task.task_type || "-") + "</div>" +
        "<div><strong>Status:</strong> " + (task.status || "-") + "</div>" +
        "<div><strong>Progress:</strong> " + (task.progress ?? "-") + "%</div>" +
        "<div><strong>Checkpoint:</strong> " + (task.last_checkpoint || "-") + "</div>",
      false
    );
  }

  function stopPolling() {
    if (window.ParkingApp.avpStatusPoller) {
      clearInterval(window.ParkingApp.avpStatusPoller);
      window.ParkingApp.avpStatusPoller = null;
    }
  }

  function startPolling() {
    stopPolling();
    window.ParkingApp.avpStatusPoller = setInterval(function () {
      window.queryAVPTask(true);
    }, 5000);
  }

  window.startAVPTask = async function () {
    try {
      const payload = {
        vehicle_id: getInputValue("avp-vehicle-id"),
        parking_lot_id: getInputValue("avp-lot-id"),
        dropoff_zone: getInputValue("avp-dropoff-zone"),
        target_space_id: getInputValue("avp-space-id"),
      };
      const result = await window.ParkingApi.startAVP(payload);
      document.getElementById("avp-task-id").value = result.task.id;
      renderTask(result.task, result.parking_session);
      startPolling();
      window.showSuccess("AVP auto-park task started");
    } catch (error) {
      setAVPStatus("Failed to start AVP task: " + error.message, true);
    }
  };

  window.startSummonTask = async function () {
    try {
      const payload = {
        vehicle_id: getInputValue("avp-vehicle-id"),
        parking_lot_id: getInputValue("avp-lot-id"),
        pickup_zone: getInputValue("avp-pickup-zone"),
      };
      const result = await window.ParkingApi.summonAVP(payload);
      document.getElementById("avp-task-id").value = result.task.id;
      renderTask(result.task);
      startPolling();
      window.showSuccess("AVP summon task started");
    } catch (error) {
      setAVPStatus("Failed to start summon task: " + error.message, true);
    }
  };

  window.queryAVPTask = async function (silent) {
    const taskId = getInputValue("avp-task-id");
    if (!taskId) {
      if (!silent) window.showError("Please provide a task ID");
      return;
    }
    try {
      const result = await window.ParkingApi.getAVPTask(taskId);
      renderTask(result.task);
      if (result.task.status === "completed" || result.task.status === "cancelled") {
        stopPolling();
      }
    } catch (error) {
      setAVPStatus("Failed to query AVP task: " + error.message, true);
      if (!silent) stopPolling();
    }
  };

  window.cancelAVPTask = async function () {
    const taskId = getInputValue("avp-task-id");
    if (!taskId) {
      window.showError("Please provide a task ID");
      return;
    }
    try {
      const result = await window.ParkingApi.cancelAVP(taskId);
      renderTask(result.task);
      stopPolling();
      window.showSuccess("AVP task cancelled");
    } catch (error) {
      setAVPStatus("Failed to cancel AVP task: " + error.message, true);
    }
  };
})();
