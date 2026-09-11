import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "ANR Admin - Platform Dashboard",
  description: "Management dashboard for ANR Studio applications and platform services.",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <body>
        <div className="dashboard-layout">
          <aside className="sidebar">
            <div className="sidebar-brand">ANR PLATFORM</div>
            <nav className="nav-menu">
              <a href="/" className="nav-item active">Overview</a>
              <a href="#apps" className="nav-item">Applications</a>
              <a href="#devices" className="nav-item">Devices</a>
              <a href="#users" className="nav-item">Users</a>
              <a href="#settings" className="nav-item">Settings</a>
            </nav>
          </aside>
          <main className="main-content">{children}</main>
        </div>
      </body>
    </html>
  );
}
