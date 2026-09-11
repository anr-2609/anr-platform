export default function AdminOverviewPage() {
  return (
    <div>
      <div className="top-bar">
        <h1>Platform Overview</h1>
        <span className="env-badge">Production</span>
      </div>

      <div className="stats-grid">
        <div className="stat-card">
          <div className="stat-label">Applications</div>
          <div className="stat-value">1</div>
        </div>
        <div className="stat-card">
          <div className="stat-label">Active Devices</div>
          <div className="stat-value">Live</div>
        </div>
        <div className="stat-card">
          <div className="stat-label">Database</div>
          <div className="stat-value" style={{ color: "var(--success)" }}>Healthy</div>
        </div>
        <div className="stat-card">
          <div className="stat-label">Redis Cache</div>
          <div className="stat-value" style={{ color: "var(--success)" }}>Connected</div>
        </div>
      </div>

      <div className="table-card">
        <div className="table-header">
          <h2>Registered Applications</h2>
        </div>
        <table className="data-table">
          <thead>
            <tr>
              <th>App ID</th>
              <th>Name</th>
              <th>Status</th>
              <th>Platform</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td><code>anr-001-wallpaper</code></td>
              <td>ANR Wallpaper</td>
              <td><span className="badge-active">active</span></td>
              <td>Android / iOS</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  );
}
