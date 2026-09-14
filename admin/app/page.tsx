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
    status: "operational",
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
    <>
      <div className="page-title-row">
        <div className="page-title-wrap">
          <h1 className="page-title">Platform Overview</h1>
          <span className="page-subtitle">Real-time ecosystem metrics & device activity</span>
        </div>
      </div>

      <div className="stats-grid">
        <div className="stat-card">
          <span className="stat-icon">+</span>
          <div className="stat-value-wrap">
            <span className="stat-value">{overview.apps_count}</span>
          </div>
          <span className="stat-label">Applications</span>
        </div>

        <div className="stat-card">
          <span className="stat-icon">*</span>
          <div className="stat-value-wrap">
            <span className="stat-value">{overview.devices_count}</span>
          </div>
          <span className="stat-label">Registered Devices</span>
        </div>

        <div className="stat-card">
          <span className="stat-icon">#</span>
          <div className="stat-value-wrap">
            <span className="stat-value">{overview.users_count}</span>
          </div>
          <span className="stat-label">Admin Accounts</span>
        </div>

        <div className="stat-card">
          <span className="stat-icon">%</span>
          <div className="stat-value-wrap">
            <span className="stat-value status-text">{overview.status}</span>
          </div>
          <span className="stat-label">Core Engine</span>
        </div>
      </div>

      <div id="apps" className="table-card">
        <div className="table-header">
          <h2 className="table-title">Ecosystem Applications</h2>
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
              <td><span className="code-pill">anr-001-wallpaper</span></td>
              <td>ANR Wallpaper</td>
              <td>
                <span className="badge-status">
                  <span className="badge-status-dot" />
                  Active
                </span>
              </td>
              <td>Android</td>
            </tr>
          </tbody>
        </table>
      </div>

      <div id="devices" className="table-card">
        <div className="table-header">
          <h2 className="table-title">Registered Devices</h2>
        </div>
        {loading ? (
          <p className="table-empty">Loading devices...</p>
        ) : devices.length === 0 ? (
          <p className="table-empty">No devices registered yet.</p>
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
                  <td><span className="code-pill">{d.device_id}</span></td>
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
    </>
  );
}
