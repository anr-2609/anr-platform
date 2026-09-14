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
    }
  }, [pathname, router]);

  function handleLogout() {
    localStorage.removeItem("anr_admin_token");
    localStorage.removeItem("anr_admin_user");
    router.replace("/login");
  }

  if (pathname === "/login") {
    return (
      <html lang="en">
        <body>{children}</body>
      </html>
    );
  }

  return (
    <html lang="en">
      <body>
        {authorized ? (
          <div className="dashboard-layout">
            <aside className="sidebar">
              <div className="sidebar-brand">ANR PLATFORM</div>
              <nav className="nav-menu">
                <a href="/" className="nav-item active">Overview</a>
                <a href="#apps" className="nav-item">Applications</a>
                <a href="#devices" className="nav-item">Devices</a>
                <button onClick={handleLogout} className="nav-item btn-logout">
                  Sign Out
                </button>
              </nav>
            </aside>
            <main className="main-content">{children}</main>
          </div>
        ) : null}
      </body>
    </html>
  );
}
