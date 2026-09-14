"use client";

import { useEffect, useState } from "react";
import { usePathname, useRouter } from "next/navigation";
import "./globals.css";

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const pathname = usePathname();
  const router = useRouter();
  const [authorized, setAuthorized] = useState(false);
  const [adminUser, setAdminUser] = useState<{ email?: string; role?: string } | null>(null);

  useEffect(() => {
    if (pathname === "/login") {
      setAuthorized(true);
      return;
    }

    const token = localStorage.getItem("anr_admin_token");
    if (!token) {
      router.replace("/login");
    } else {
      setAuthorized(true);
      try {
        const stored = localStorage.getItem("anr_admin_user");
        if (stored) {
          setAdminUser(JSON.parse(stored));
        }
      } catch (e) {
        setAdminUser({ email: "admin@anr-studio.com", role: "admin" });
      }
    }
  }, [pathname, router]);

  function handleLogout() {
    localStorage.removeItem("anr_admin_token");
    localStorage.removeItem("anr_admin_user");
    router.replace("/login");
  }

  return (
    <html lang="en">
      <head>
        <title>ANR Platform — Admin Dashboard</title>
        <meta name="description" content="Admin Dashboard for ANR Studio Ecosystem" />
        <link rel="preconnect" href="https://fonts.googleapis.com" />
        <link rel="preconnect" href="https://fonts.gstatic.com" crossOrigin="anonymous" />
        <link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&display=swap" rel="stylesheet" />
        <link href="https://db.onlinewebfonts.com/c/8cb707a9b8a73f8a7403336b861c3074?family=BubbledotICG-FinePos" rel="stylesheet" />
        <link rel="icon" href="/assets/logo.webp?v=2" type="image/webp" />
        <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.5.2/css/all.min.css" integrity="sha512-SnH5WK+bZxgPHs44uWIX+LLJAJ9/2PkPKZ5QiAj6Ta86w+fsb2TkcmfRyVX3pBnMFcV7oQPJkl9QevSCWr3W6A==" crossOrigin="anonymous" referrerPolicy="no-referrer" />
      </head>
      <body>
        {pathname === "/login" ? (
          children
        ) : authorized ? (
          <div className="dashboard-layout">
            <aside className="dashboard-sidebar">
              <div className="sidebar-brand">
                <a href="/" className="logo-btn" aria-label="Dashboard">
                  <img
                    src="/assets/logo.webp?v=2"
                    alt="ANR Studio"
                    width={38}
                    height={38}
                    className="logo-img"
                  />
                </a>
                <div className="brand-info">
                  <span className="brand-title">ANR PLATFORM</span>
                  <span className="brand-badge">MANAGEMENT CONSOLE</span>
                </div>
              </div>

              <div className="sidebar-section-title">NAVIGATION</div>
              <nav className="sidebar-nav">
                <a href="/" className={`sidebar-nav-item ${pathname === "/" ? "active" : ""}`}>
                  <i className="fa-solid fa-chart-pie nav-icon" />
                  <span>Platform Overview</span>
                </a>
                <a href="/#apps" className="sidebar-nav-item">
                  <i className="fa-solid fa-cubes nav-icon" />
                  <span>Applications</span>
                </a>
                <a href="/#devices" className="sidebar-nav-item">
                  <i className="fa-solid fa-microchip nav-icon" />
                  <span>Devices</span>
                </a>
                <a href="http://localhost:3000/status" target="_blank" rel="noopener noreferrer" className="sidebar-nav-item">
                  <i className="fa-solid fa-heart-pulse nav-icon" />
                  <span>System Status</span>
                  <i className="fa-solid fa-arrow-up-right-from-square nav-external-icon" />
                </a>
              </nav>

              <div className="sidebar-footer">
                <div className="admin-profile">
                  <div className="admin-avatar">
                    <i className="fa-solid fa-user-shield" />
                  </div>
                  <div className="admin-meta">
                    <span className="admin-name">{adminUser?.email || "admin@anr-studio.com"}</span>
                    <span className="admin-role">Administrator</span>
                  </div>
                </div>
                <button
                  type="button"
                  onClick={handleLogout}
                  className="btn-sidebar-logout"
                  title="Sign out"
                >
                  <i className="fa-solid fa-arrow-right-from-bracket" />
                  <span>Sign Out</span>
                </button>
              </div>
            </aside>

            <div className="dashboard-body">
              <header className="dashboard-topbar">
                <div className="topbar-breadcrumb">
                  <span className="breadcrumb-root">Admin</span>
                  <span className="breadcrumb-sep">/</span>
                  <span className="breadcrumb-active">Platform Overview</span>
                </div>
                <div className="topbar-actions">
                  <div className="badge-production">
                    <span className="pulse-dot" />
                    Production Environment
                  </div>
                  <a
                    href="https://anr-studio.com"
                    target="_blank"
                    rel="noopener noreferrer"
                    className="topbar-link-btn"
                  >
                    <i className="fa-solid fa-arrow-up-right-from-square" />
                    <span>View Landing</span>
                  </a>
                </div>
              </header>

              <main className="dashboard-content">
                {children}
              </main>
            </div>
          </div>
        ) : (
          <div style={{ minHeight: "100vh", display: "grid", placeItems: "center", background: "#06090e", color: "#888888" }}>
            <span style={{ fontSize: "13px", letterSpacing: "0.08em", textTransform: "uppercase" }}>Redirecting to login...</span>
          </div>
        )}
      </body>
    </html>
  );
}
