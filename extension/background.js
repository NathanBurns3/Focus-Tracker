let currentTabId = null;
let currentDomain = null;
let startTime = null;

// Listen for tab activation
chrome.tabs.onActivated.addListener(async (activeInfo) => {
  await trackTime(); // save previous tab's usage

  const tab = await chrome.tabs.get(activeInfo.tabId);
  currentTabId = tab.id;
  currentDomain = extractDomain(tab.url);
  startTime = Date.now();
});

// Listen for tab updates (e.g., URL change)
chrome.tabs.onUpdated.addListener(async (tabId, changeInfo) => {
  if (tabId === currentTabId && changeInfo.url) {
    await trackTime();
    currentDomain = extractDomain(changeInfo.url);
    startTime = Date.now();
  }
});

// Listen for tab removal
chrome.tabs.onRemoved.addListener(async (tabId) => {
  if (tabId === currentTabId) {
    await trackTime();
    currentTabId = null;
    currentDomain = null;
    startTime = null;
  }
});

// Listen for window focus changes
chrome.windows.onFocusChanged.addListener(async (windowId) => {
  if (windowId === chrome.windows.WINDOW_ID_NONE) {
    // User switched away from Chrome
    await trackTime();
    startTime = null;
  } else if (currentTabId !== null) {
    // User returned to Chrome, reset timer
    startTime = Date.now();
  }
});

// Listen for extension suspend (browser shutdown)
chrome.runtime.onSuspend.addListener(async () => {
  await trackTime();
});

// Function to calculate time spent and save it
async function trackTime() {
  if (!currentDomain || !startTime) return;

  const elapsedMinutes = (Date.now() - startTime) / 60000; // convert ms to minutes

  // Get current logs from storage, update them, and save back
  chrome.storage.local.get("usageLogs", ({ usageLogs }) => {
    const logs = usageLogs || {};
    logs[currentDomain] = (logs[currentDomain] || 0) + elapsedMinutes;
    chrome.storage.local.set({ usageLogs: logs });
  });
}

// Utility to extract domain from URL
function extractDomain(url) {
  try {
    const u = new URL(url);
    return u.hostname;
  } catch (e) {
    return "unknown";
  }
}

// Function to POST usage logs to Go server
function postUsageLogs() {
  chrome.storage.local.get("usageLogs", ({ usageLogs }) => {
    if (!usageLogs || Object.keys(usageLogs).length === 0) return;

    // Convert to array of {domain, minutes}
    const entries = Object.entries(usageLogs).map(([domain, minutes]) => ({
      domain,
      minutes,
    }));

    // Clear logs BEFORE sending to prevent duplicates
    chrome.storage.local.set({ usageLogs: {} });

    fetch("http://localhost:8080/usage", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(entries),
    })
      .then(() => {
        console.log("Posted usage logs:", entries);
      })
      .catch((err) => {
        console.log("Failed to post usage logs (server down?)");
      });
  });
}

// Send logs every 30 seconds
setInterval(postUsageLogs, 30 * 1000);
