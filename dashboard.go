package main

// dashboardHTML is the entire single-file frontend. It fetches /api/status,
// renders a live status page styled with the dui.css system via shell.js.
const dashboardHTML = `<!DOCTYPE html>
<html lang="en" data-theme="dark" class="dark">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <title>Up — Real-Time Health & Status Dashboard | dhanur.me</title>
  <meta name="description" content="Real-time uptime monitoring and system health dashboard across all dhanur.me microservices, APIs, and infrastructure endpoints." />
  <link rel="canonical" href="https://up.dhanur.me/" />
  <link rel="icon" type="image/x-icon" href="https://dhanur.me/favicon.ico" />

  <!-- Open Graph / Facebook -->
  <meta property="og:type" content="website" />
  <meta property="og:url" content="https://up.dhanur.me/" />
  <meta property="og:title" content="Up — Real-Time Health & Status Dashboard | dhanur.me" />
  <meta property="og:description" content="Real-time uptime monitoring and system health dashboard across all dhanur.me microservices, APIs, and infrastructure endpoints." />
  <meta property="og:image" content="https://raw.githubusercontent.com/kascit/kascit.github.io/raw-mirror/images/branding/logo-light.png" />

  <!-- Twitter -->
  <meta name="twitter:card" content="summary" />
  <meta name="twitter:title" content="Up — Real-Time Health & Status Dashboard | dhanur.me" />
  <meta name="twitter:description" content="Real-time uptime monitoring and system health dashboard across all dhanur.me microservices, APIs, and infrastructure endpoints." />

  <!-- Trusted Types default policy (must run before ANY script that uses innerHTML/script.src) -->
  <script>
  (function(){
    if (window.trustedTypes && window.trustedTypes.createPolicy) {
      try {
        window.__defaultPolicy = window.trustedTypes.createPolicy('default', {
          createScriptURL: function(s) { return s; },
          createHTML: function(s) { return s; }
        });
      } catch(e) {}
    }
    function installFramePolicy(m) {
      var orig = Element.prototype[m];
      Element.prototype[m] = function() {
        var r = orig.apply(this, arguments);
        var n = arguments[0];
        if (n && n.tagName === 'IFRAME') {
          try {
            var w = n.contentWindow;
            if (w && w.trustedTypes && !w.trustedTypes.defaultPolicy) {
              w.trustedTypes.createPolicy('default', {
                createHTML: function(s){return s},
                createScript: function(s){return s},
                createScriptURL: function(s){return s}
              });
            }
          } catch(e) {}
        }
        return r;
      };
    }
    installFramePolicy('appendChild');
    installFramePolicy('insertBefore');
  })();
  </script>

  <script>
  window.SiteNavConfig = {
    mode: "navbar",
    nav: [{ name: "Home", url: "https://dhanur.me/", icon: "fa-solid fa-house" }],
    sidebarNav: [],
    activePath: window.location.pathname,
    badge: { text: "STATUS", class: "badge-neutral" },
    showSearch: false,
    showAppsGrid: true,
    showThemeToggle: true,
    enablePwa: false,
    favicon: false
  };
  </script>
  <link rel="stylesheet" href="https://dhanur.me/css/font-awesome.min.css" />
  <link rel="stylesheet" href="https://dhanur.me/css/dui.css" />
  <script src="https://dhanur.me/js/shell.min.js" defer></script>

  <style>
    /* ── Pulse animation for live indicator ── */
    @keyframes pulse-ring {
      0%   { box-shadow: 0 0 0 0 color-mix(in oklab, var(--color-success) 55%, transparent); }
      70%  { box-shadow: 0 0 0 7px transparent; }
      100% { box-shadow: 0 0 0 0 transparent; }
    }
    .pulse { animation: pulse-ring 2s cubic-bezier(0.455,0.03,0.515,0.955) infinite; }

    /* ── Spark bar (90 ticks) ── */
    .spark-bar { display: flex; gap: 2px; align-items: flex-end; height: 28px; }
    .spark-tick {
      flex: 1 1 0;
      min-width: 2px;
      border-radius: 2px;
      transition: height 0.3s ease;
    }
    .tick-up      { background: var(--color-success); opacity: 0.8; }
    .tick-degraded{ background: var(--color-warning); opacity: 0.8; }
    .tick-down    { background: var(--color-error);   opacity: 0.8; }
    .tick-unknown { background: color-mix(in oklab, var(--color-base-content) 20%, transparent); }

    /* ── Status pill ── */
    .pill {
      display: inline-flex; align-items: center; gap: 6px;
      font-size: 11px; font-weight: 700; letter-spacing: 0.05em;
      text-transform: uppercase; padding: 3px 10px;
      border-radius: 999px;
    }
    .pill-up       { background: color-mix(in oklab, var(--color-success) 14%, transparent); color: var(--color-success); border: 1px solid color-mix(in oklab, var(--color-success) 30%, transparent); }
    .pill-degraded { background: color-mix(in oklab, var(--color-warning) 14%, transparent); color: var(--color-warning); border: 1px solid color-mix(in oklab, var(--color-warning) 30%, transparent); }
    .pill-down     { background: color-mix(in oklab, var(--color-error) 14%, transparent);   color: var(--color-error);   border: 1px solid color-mix(in oklab, var(--color-error) 30%, transparent); }
    .pill-unknown  { background: color-mix(in oklab, var(--color-base-content) 8%, transparent); color: var(--color-base-content); border: 1px solid color-mix(in oklab, var(--color-base-content) 15%, transparent); }

    /* ── Latency badge ── */
    .latency-chip {
      font-size: 10px; font-family: ui-monospace, monospace;
      color: color-mix(in oklab, var(--color-base-content) 45%, transparent);
      padding: 2px 6px;
      background: color-mix(in oklab, var(--color-base-content) 6%, transparent);
      border-radius: 4px;
    }
  </style>
</head>
<body class="bg-base-100 text-base-content antialiased min-h-screen">
<header id="dhanur-shell"></header>
<main class="eco-page eco-container">

  <!-- Page header -->
  <div class="mb-10">
    <div class="flex items-center gap-3 mb-3">
      <span id="global-dot" class="w-3 h-3 rounded-full bg-base-content/20"></span>
      <h1 class="text-4xl lg:text-5xl font-extrabold tracking-tight" id="global-title">Checking…</h1>
    </div>
    <p class="text-base-content/60 text-sm" id="global-sub">Fetching service health…</p>
  </div>

  <!-- Service cards -->
  <div id="cards" class="flex flex-col gap-4"></div>

  <!-- Footer -->
  <footer class="eco-footer">
    <span>&copy; <span id="year">2026</span> Dhanur Relhan · refreshing every 30s</span>
    <div class="flex items-center gap-4">
      <span id="last-checked" class="font-mono text-xs opacity-60">—</span>
      <a href="https://dhanur.me/" class="hover:text-base-content transition-colors">dhanur.me</a>
      <a href="https://github.com/kascit/up" class="hover:text-base-content transition-colors" target="_blank" rel="noopener noreferrer">Source</a>
    </div>
  </footer>
</main>

<script>
const STATUS_COLORS = {
  up:       'pill-up',
  degraded: 'pill-degraded',
  down:     'pill-down',
  unknown:  'pill-unknown',
};
const STATUS_LABELS = {
  up:       'Operational',
  degraded: 'Degraded',
  down:     'Down',
  unknown:  'Unknown',
};
const GLOBAL_TEXT = {
  up:       ['All Systems Operational', 'Everything looks good.'],
  degraded: ['Partial Outage', 'Some services are experiencing issues.'],
  down:     ['Major Outage', 'One or more services are unreachable.'],
  unknown:  ['Checking…', 'Fetching service health…'],
};
const GLOBAL_DOT = {
  up:       'bg-success pulse',
  degraded: 'bg-warning pulse',
  down:     'bg-error pulse',
  unknown:  'bg-base-content/20',
};

function latencyText(ms) {
  if (ms < 0) return 'n/a';
  if (ms < 1000) return ms + ' ms';
  return (ms / 1000).toFixed(2) + ' s';
}

function tickClass(status) {
  switch (status) {
    case 'up':       return 'tick-up';
    case 'degraded': return 'tick-degraded';
    case 'down':     return 'tick-down';
    default:         return 'tick-unknown';
  }
}

function tickHeight(status) {
  switch (status) {
    case 'up':       return '100%';
    case 'degraded': return '55%';
    case 'down':     return '20%';
    default:         return '35%';
  }
}

function renderSpark(history) {
  // history is newest-first; we want to display oldest-left
  const slots = history.slice().reverse();
  // Pad to historyLen if shorter
  while (slots.length < 90) slots.unshift(null);
  return slots.map(h => {
    const cls  = h ? tickClass(h.status) : 'tick-unknown';
    const hpct = h ? tickHeight(h.status) : '35%';
    return '<div class="spark-tick ' + cls + '" style="height:' + hpct + '"></div>';
  }).join('');
}

function renderCard(svc) {
  const cur  = svc.current;
  const pill = STATUS_COLORS[cur.status] || 'pill-unknown';
  const lbl  = STATUS_LABELS[cur.status] || 'Unknown';

  return '<div class="eco-card p-5">' +
    '<div class="flex flex-col gap-3">' +
      // Row 1: icon + name + status + latency
      '<div class="flex items-center justify-between gap-3">' +
        '<div class="flex items-center gap-3">' +
          '<div class="flex items-center justify-center w-9 h-9 rounded-xl bg-base-200 border border-base-content/10 flex-shrink-0">' +
            '<i class="' + svc.config.icon + ' text-base-content/70"></i>' +
          '</div>' +
          '<div>' +
            '<a href="' + svc.config.homepageUrl + '" target="_blank" rel="noopener"' +
              ' class="font-bold text-[15px] text-base-content hover:underline decoration-1 underline-offset-2">' +
              svc.config.displayName +
            '</a>' +
            '<div class="text-[10px] text-base-content/40 font-mono mt-0.5">' + svc.config.category + '</div>' +
          '</div>' +
        '</div>' +
        '<div class="flex items-center gap-2">' +
          '<span class="latency-chip">' + latencyText(cur.latencyMs) + '</span>' +
          '<span class="pill ' + pill + '">' + lbl + '</span>' +
        '</div>' +
      '</div>' +
      // Row 2: spark + uptime
      '<div>' +
        '<div class="spark-bar">' + renderSpark(svc.history) + '</div>' +
        '<div class="flex justify-between text-[10px] text-base-content/35 mt-1 font-mono">' +
          '<span>45 min ago</span>' +
          '<span>' + svc.uptimePercent.toFixed(1) + '% uptime</span>' +
        '</div>' +
      '</div>' +
    '</div>' +
  '</div>';
}

async function refresh() {
  try {
    const res  = await fetch('/api/status');
    if (!res.ok) throw new Error('HTTP ' + res.status);
    const data = await res.json();

    // Update global header
    const [title, sub] = GLOBAL_TEXT[data.overall] || GLOBAL_TEXT.unknown;
    document.getElementById('global-title').textContent = title;
    document.getElementById('global-sub').textContent   = sub;

    const dot = document.getElementById('global-dot');
    dot.className = 'w-3 h-3 rounded-full ' + (GLOBAL_DOT[data.overall] || GLOBAL_DOT.unknown);

    // Render cards
    document.getElementById('cards').innerHTML =
      data.services.map(renderCard).join('');

    // Footer timestamp
    document.getElementById('last-checked').textContent =
      'Last checked: ' + new Date(data.updatedAt).toLocaleTimeString();

  } catch (e) {
    console.error('Failed to fetch status:', e);
  }
}

const yearEl = document.getElementById('year');
if (yearEl) yearEl.textContent = new Date().getFullYear();

// Initial load + polling every 30 s (matching server-side interval)
refresh();
setInterval(refresh, 30_000);
</script>
</body>
</html>`
