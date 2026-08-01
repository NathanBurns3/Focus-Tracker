# Focus Tracker

A comprehensive productivity tracking system that monitors application and website usage on your computer. Focus Tracker combines a macOS daemon for desktop app tracking, a Chrome extension for website monitoring, and a beautiful dashboard for visualizing your productivity insights.

## Features

- **Desktop App Tracking**: Automatically tracks which applications you're using and for how long
- **Website Monitoring**: Chrome extension logs time spent on websites
- **Real-time Dashboard**: Beautiful Next.js web interface with interactive charts
- **Multiple Chart Types**: Visualize data with bar, line, or pie charts
- **Daily Reports**: CLI command to generate terminal-based productivity reports
- **Application Aliasing**: Customize display names for apps and websites

## Architecture

### Components

1. **CLI (Go)** - Background daemon and command-line interface

   - `serve` runs as a launchd-managed daemon (installed via `focus-tracker install`)
   - Polls for the active macOS application every 10 seconds
   - Starts an HTTP server to receive data from the Chrome extension
   - Generates terminal reports of daily usage
   - Manages database connections and data insertion

2. **Dashboard (Next.js)** - Web-based visualization

   - Displays top 10 apps and websites for the current day
   - Interactive charts with multiple visualization options
   - Real-time data fetching from API routes
   - Dark mode design for comfortable viewing

3. **Chrome Extension** - Website usage tracking

   - Tracks active tabs and time spent on each domain
   - Sends aggregated usage data to the local server every 30 seconds
   - Handles tab switches, window focus changes, and extension suspend events

4. **Database (PostgreSQL)** - Data persistence
   - Stores all usage records with timestamps
   - Supports both desktop (source: "desktop") and website (source: "chrome") tracking

## Database Schema

```sql
CREATE TABLE daily_usage (
   id SERIAL PRIMARY KEY,
   name TEXT NOT NULL,
   usage_date DATE NOT NULL DEFAULT NOW(),
   minutes_used NUMERIC(5,2) NOT NULL DEFAULT 0,
   source TEXT NOT NULL CHECK (source IN ('desktop', 'chrome')) DEFAULT 'desktop'
);
```

- `name`: Application or website name
- `usage_date`: Date of the usage record
- `minutes_used`: Total minutes spent
- `source`: Either "desktop" (macOS app) or "chrome" (website)

Timezone note: the CLI sets the PostgreSQL session time zone to `America/New_York`. Set `TZ` before running or adjust in `cli/internal/db/db.go` if you need a different zone.

## Installation

### Prerequisites

- Go 1.16+ (for CLI)
- Node.js 16+ (for dashboard)
- PostgreSQL 13+ (or Docker)
- Chrome browser (for extension)
- macOS (daemon uses osascript/AppleScript)

### Setup

1. **Clone the repository**

   ```bash
   git clone <repository-url>
   cd Focus-Tracker
   ```

2. **Environment Configuration**
   Create a `.env` file in the root directory:

   ```env
   POSTGRES_USER=your_user
   POSTGRES_PASSWORD=your_password
   POSTGRES_DB=trackerdb
   NEXT_PUBLIC_DB_PATH=postgres://your_user:your_password@localhost:5432/trackerdb
   ```

3. **Start PostgreSQL**
   Using Docker Compose:

   ```bash
   docker-compose up -d
   ```

4. **Setup CLI (install binary on PATH)**

   ```bash
   cd cli
   go mod download
   go install .
   echo 'export PATH="$HOME/go/bin:$PATH"' >> ~/.zshrc
   source ~/.zshrc
   ```

5. **Setup Dashboard**

   ```bash
   cd dashboard
   npm install
   ```

6. **Load Chrome Extension**
   - Open Chrome and navigate to `chrome://extensions`
   - Enable "Developer mode" (top right)
   - Click "Load unpacked"
   - Select the `extension/` directory

## Usage

### Starting the Tracker

1. **One-time install of the launchd service** (from `cli/` directory):

   ```bash
   focus-tracker install
   ```

   This writes `~/Library/LaunchAgents/com.focustracker.daemon.plist` pointing to `focus-tracker serve` and loads it with launchd.

2. **Start the daemon** (from anywhere once `focus-tracker` is on PATH):

   ```bash
   focus-tracker start
   ```

   This will:

   - Automatically start Docker Desktop
   - Start polling for active applications every 10 seconds
   - Start an HTTP server on `http://localhost:8080` to receive website data
   - Insert/update usage records in the database

3. **Start the dashboard** (from `dashboard/` directory):

   ```bash
   npm run dev
   ```

   Access at `http://localhost:3000`

4. **The Chrome extension** runs automatically once loaded

   - Tracks active tabs in the background
   - Sends usage data to the local server every 30 seconds

5. **Check daemon status (optional)**

   ```bash
   focus-tracker status
   ```

### Stopping the Tracker

```bash
focus-tracker stop
```

This will:

- Stop the background daemon
- Automatically quit Docker Desktop

### Viewing Daily Reports

```bash
focus-tracker report
```

Displays a formatted terminal report with:

- Top 10 applications by usage time
- Top 10 websites by usage time
- ASCII bar charts for visual comparison
- Color-coded output

### Application Aliases

Customize display names by editing `cli/internal/config/aliases.json`:

```json
{
  "Electron": "VS Code",
  "chatgpt.com": "ChatGPT",
  "docs.google.com": "Google Docs"
}
```

## API Routes

### GET `/api/usage`

Fetches top 10 apps and websites for the current day.

**Response:**

```json
[
  {
    "name": "VS Code",
    "minutes_used": 120.5,
    "source": "desktop"
  },
  {
    "name": "ChatGPT",
    "minutes_used": 45.25,
    "source": "chrome"
  }
]
```

## Project Structure

```
Focus-Tracker/
├── cli/
│   ├── main.go
│   ├── cmd/
│   │   ├── install.go   # Writes/loads launchd plist (serve)
│   │   ├── start.go     # launchctl start
│   │   ├── stop.go      # launchctl stop/unload
│   │   ├── status.go    # launchctl list wrapper
│   │   ├── serve.go     # Daemon entrypoint (called by launchd)
│   │   └── report.go    # Terminal reports
│   └── internal/
│       ├── config/      # Configuration & aliases
│       ├── tracker/     # Application polling (formerly daemon)
│       ├── db/          # Database operations
│       └── server/      # HTTP server for extension
├── dashboard/
│   ├── src/
│   │   ├── app/
│   │   │   ├── page.tsx             # Main dashboard
│   │   │   └── api/usage/route.ts   # API endpoint
│   │   └── components/
│   │       └── UsageChart.tsx       # Chart visualization
│   └── package.json
├── extension/
│   ├── manifest.json
│   └── background.js
├── docker-compose.yml
├── .env
└── README.md
```

## Development

### CLI Development

The CLI uses the Cobra framework for command management. To add new commands, create a new file in `cli/cmd/` and register it in the `init()` function.

### Dashboard Development

The dashboard uses Next.js with TypeScript and Tailwind CSS. Charts are built with Recharts for interactive visualizations.

### Extension Development

The extension uses Chrome's Manifest V3 API. Modify `extension/background.js` to change tracking behavior or `extension/manifest.json` to adjust permissions.
