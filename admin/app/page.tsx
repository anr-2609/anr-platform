"use client";

import { useEffect, useState } from "react";

interface Device {
  id: number;
  device_id: string;
  app_id: string;
  platform: string;
  os_version: string;
  app_version: string;
  last_active_at: string;
}

interface OverviewData {
  apps_count: number;
  devices_count: number;
  users_count: number;
  status: string;
}

export default function AdminOverviewPage() {
  const [overview, setOverview] = useState<OverviewData>({
    apps_count: 1,
    devices_count: 0,
    users_count: 1,
    status: "connecting",
  });
  const [devices, setDevices] = useState<Device[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const token = localStorage.getItem("anr_admin_token");
    if (!token) return;

    async function loadData() {
      try {
        const [overviewRes, devicesRes] = await Promise.all([
          fetch("https://api.anr-studio.com/api/v1/admin/overview", {
            headers: { Authorization: `Bearer ${token}` },
          }),
          fetch("https://api.anr-studio.com/api/v1/admin/devices", {
            headers: { Authorization: `Bearer ${token}` },
          }),
        ]);

        if (overviewRes.ok) {
          const json = await overviewRes.json();
          if (json.success) {
            setOverview(json.data);
          }
        }

        if (devicesRes.ok) {
          const json = await devicesRes.json();
          if (json.success && json.data?.devices) {
            setDevices(json.data.devices);
          }
        }
      } catch (err) {
        console.error(err);
      } finally {
        setLoading(false);
      }
    }

    loadData();
  }, []);

  return (
    <div>
      <div className="top-bar">
        <h1>Platform Overview</h1>
        <span className="env-badge">Production</span>
      </div>

      <div className="stats-grid">
        <div className="stat-card">
          <div className="stat-label">Applications</div>
          <div className="stat-value">{overview.apps_count}</div>
        </div>
        <div className="stat-card">
          <div className="stat-label">Registered Devices</div>
          <div className="stat-value">{overview.devices_count}</div>
        </div>
        <div className="stat-card">
          <div className="stat-label">Admin Accounts</div>
          <div className="stat-value">{overview.users_count}</div>
        </div>
        <div className="stat-card">
          <div className="stat-label">Core Engine</div>
          <div className="stat-value" style={{ color: "var(--success)" }}>
            {overview.status}
          </div>
        </div>
      </div>

      <div id="apps" className="table-card">
        <div className="table-header">
          <h2>Ecosystem Applications</h2>
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

      <div id="devices" className="table-card">
        <div className="table-header">
          <h2>Registered Devices</h2>
        </div>
        {loading ? (
          <p style={{ color: "var(--text-secondary)" }}>Loading devices...</p>
        ) : devices.length === 0 ? (
          <p style={{ color: "var(--text-secondary)" }}>No devices registered yet.</p>
        ) : (
          <table className="data-table">
            <thead>
              <tr>
                <th>Device ID</th>
                <th>App</th>
                <th>Platform</th>
                <th>OS</th>
                <th>App Version</th>
                <th>Last Active</th>
              </tr>
            </thead>
            <tbody>
              {devices.map((d) => (
                <tr key={d.id}>
                  <td><code>{d.device_id}</code></td>
                  <td>{d.app_id}</td>
                  <td>{d.platform}</td>
                  <td>{d.os_version || "N/A"}</td>
                  <td>{d.app_version || "N/A"}</td>
                  <td>{new Date(d.last_active_at).toLocaleString()}</td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
    </div>
  );
}
