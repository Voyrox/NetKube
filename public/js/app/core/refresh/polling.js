const DATA_REFRESH_INTERVAL_MS = 5000;

function startAutoRefresh(load, interval = DATA_REFRESH_INTERVAL_MS) {
  let inFlight = false;

  const run = async () => {
    if (inFlight || document.hidden) {
      return;
    }

    inFlight = true;
    try {
      await load();
    } finally {
      inFlight = false;
    }
  };

  run();
  const timer = window.setInterval(run, interval);
  window.addEventListener("beforeunload", () => window.clearInterval(timer), {
    once: true,
  });

  return {
    run,
    stop() {
      window.clearInterval(timer);
    },
  };
}
