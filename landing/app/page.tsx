export default function HomePage() {
  return (
    <main>
      <div className="container">
        <header className="navbar">
          <div className="logo">ANR STUDIO</div>
          <nav className="nav-links">
            <a href="#apps">Apps</a>
            <a href="#platform">Platform</a>
            <a href="mailto:contact@anr-studio.com">Contact</a>
          </nav>
        </header>

        <section className="hero">
          <span className="badge">Next Generation Mobile Apps</span>
          <h1>
            Crafting Exceptional <span>Digital Experiences</span>
          </h1>
          <p>
            ANR Studio creates high quality, privacy-first mobile applications designed with speed, aesthetics, and reliability at their core.
          </p>
          <div className="cta-group">
            <a href="#apps" className="btn-primary">Explore Products</a>
            <a href="mailto:contact@anr-studio.com" className="btn-secondary">Contact Us</a>
          </div>
        </section>

        <section id="apps" className="section">
          <h2 className="section-title">Ecosystem Applications</h2>
          <p className="section-desc">Mobile applications powered by the unified ANR Platform architecture.</p>
          
          <div className="grid">
            <div className="card">
              <h3>ANR Wallpaper</h3>
              <p>
                Curated 4K and Ultra-HD wallpapers engineered for vibrant screens with lightweight performance and smooth categorization.
              </p>
            </div>
            <div className="card">
              <h3>Unified Platform</h3>
              <p>
                High availability backend infrastructure with first-party device synchronization, real-time analytics, and low latency caching.
              </p>
            </div>
            <div className="card">
              <h3>Privacy Centric</h3>
              <p>
                Strict data privacy standards with transparent session management, zero intrusive tracking, and secure guest access.
              </p>
            </div>
          </div>
        </section>

        <footer className="footer">
          <p>&copy; {new Date().getFullYear()} ANR Studio. All rights reserved.</p>
        </footer>
      </div>
    </main>
  );
}
