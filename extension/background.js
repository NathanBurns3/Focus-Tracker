let currentTabId = null;
let currentDomain = null;
let startTime = null;

const usageLogs = {}; // { domain: minutes }

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

// Function to calculate time spent and save it
async function trackTime() {
  if (!currentDomain || !startTime) return;

  const elapsedMinutes = (Date.now() - startTime) / 60000; // convert ms to minutes
  usageLogs[currentDomain] = (usageLogs[currentDomain] || 0) + elapsedMinutes;

  // Save logs to chrome.storage.local for periodic POSTing
  chrome.storage.local.set({ usageLogs });
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
    if (!usageLogs) return;

    // Convert to array of {domain, minutes}
    const entries = Object.entries(usageLogs).map(([domain, minutes]) => ({
      domain,
      minutes,
    }));

    fetch("http://localhost:8080/usage", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(entries),
    })
      .then(() => {
        console.log("Posting usage logs:", entries);
        // Clear logs after successful POST
        chrome.storage.local.set({ usageLogs: {} });
      })
      .catch((err) => {
        // Fail silently if Go server is not running
      });
  });
}

// Send logs every 30 seconds
setInterval(postUsageLogs, 30 * 1000);
