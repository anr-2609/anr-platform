"use client";

import { useEffect, useState, useRef } from "react";

interface StatConfig {
  icon: string;
  target: number;
  suffix: string;
  decimals: number;
  label: string;
  delay: string;
}

const STATS_DATA: StatConfig[] = [
  { icon: "#", target: 5, suffix: "M+", decimals: 0, label: "Global Downloads", delay: "0.5s" },
  { icon: "*", target: 4.8, suffix: " /5", decimals: 1, label: "Play Store Rating", delay: "0.58s" },
  { icon: "+", target: 12, suffix: "+", decimals: 0, label: "Published Apps", delay: "0.66s" },
  { icon: "%", target: 99.9, suffix: "%", decimals: 1, label: "Crash-Free Rate", delay: "0.74s" },
];

export default function HomePage() {
  const [activeModal, setActiveModal] = useState<"apps" | "platform" | "about" | "contact" | null>(null);
  const [copiedEmail, setCopiedEmail] = useState(false);
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);
  const [displayValues, setDisplayValues] = useState<string[]>(STATS_DATA.map(s => (0).toFixed(s.decimals)));
  const footerRef = useRef<HTMLElement>(null);
  const hasAnimatedRef = useRef(false);

  const openModal = (modal: "apps" | "platform" | "about" | "contact") => {
    setActiveModal(modal);
    setMobileMenuOpen(false);
  };

  const handleCopyEmail = () => {
    navigator.clipboard.writeText("contact@anr-studio.com");
    setCopiedEmail(true);
    setTimeout(() => setCopiedEmail(false), 2000);
  };

  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === "Escape") {
        setMobileMenuOpen(false);
        setActiveModal(null);
      }
    };

    const handleResize = () => {
      if (window.innerWidth > 720) {
        setMobileMenuOpen(false);
      }
    };

    window.addEventListener("keydown", handleKeyDown);
    window.addEventListener("resize", handleResize);

    return () => {
      window.removeEventListener("keydown", handleKeyDown);
      window.removeEventListener("resize", handleResize);
    };
  }, []);

  useEffect(() => {
    const footerElem = footerRef.current;
    if (!footerElem) return;

    const observer = new IntersectionObserver(
      (entries) => {
        const [entry] = entries;
        if (entry.isIntersecting && !hasAnimatedRef.current) {
          hasAnimatedRef.current = true;

          STATS_DATA.forEach((stat, i) => {
            const startDelay = 480 + i * 90;
            const duration = 1500 + i * 80;

            setTimeout(() => {
              const startTime = performance.now();

              const frame = (now: number) => {
                const elapsed = now - startTime;
                const progress = Math.min(elapsed / duration, 1);
                const eased = 1 - Math.pow(1 - progress, 3);
                const currentVal = stat.target * eased;

                setDisplayValues((prev) => {
                  const copy = [...prev];
                  copy[i] = currentVal.toFixed(stat.decimals);
                  return copy;
                });

                if (progress < 1) {
                  requestAnimationFrame(frame);
                } else {
                  setDisplayValues((prev) => {
                    const copy = [...prev];
                    copy[i] = stat.target.toFixed(stat.decimals);
                    return copy;
                  });
                }
              };

              requestAnimationFrame(frame);
            }, startDelay);
          });
        }
      },
      { threshold: 0.25 }
    );

    observer.observe(footerElem);

    return () => {
      observer.disconnect();
    };
  }, []);

  return (
    <>
      <div className="bg">
        <video className="bg-video" autoPlay muted loop playsInline>
          <source
            src="https://d8j0ntlcm91z4.cloudfront.net/user_38xzZboKViGWJOttwIXH07lWA1P/hf_20260809_012548_ef22562c-c0ae-4816-ad9d-f8922af4e6a7.mp4"
            type="video/mp4"
          />
        </video>
        <div className="bg-overlay" />
      </div>

      {mobileMenuOpen && (
        <>
          <div
            className="mobile-overlay"
            onClick={() => setMobileMenuOpen(false)}
          />
          <div className="mobile-sheet">
            <nav className="mobile-nav">
              <a
                href="#apps"
                className={`mobile-nav-link ${activeModal === "apps" ? "active" : ""}`}
                onClick={(e) => {
                  e.preventDefault();
                  openModal("apps");
                }}
              >
                Apps
              </a>
              <a
                href="#platform"
                className={`mobile-nav-link ${activeModal === "platform" ? "active" : ""}`}
                onClick={(e) => {
                  e.preventDefault();
                  openModal("platform");
                }}
              >
                Platform
              </a>
              <a
                href="#about"
                className={`mobile-nav-link ${activeModal === "about" ? "active" : ""}`}
                onClick={(e) => {
                  e.preventDefault();
                  openModal("about");
                }}
              >
                About
              </a>
              <a
                href="#contact"
                className={`mobile-nav-link ${activeModal === "contact" ? "active" : ""}`}
                onClick={(e) => {
                  e.preventDefault();
                  openModal("contact");
                }}
              >
                Contact
              </a>
            </nav>
          </div>
        </>
      )}

      <div className="page">
        <header className="header">
          <a
            href="#"
            className="logo-btn"
            aria-label="Home"
            onClick={(e) => {
              e.preventDefault();
              setActiveModal(null);
            }}
          >
            <img
              src="/assets/logo.webp"
              alt="ANR Studio"
              width={52}
              height={52}
              className="logo-img"
            />
          </a>

          <nav className="nav-pill desktop-nav">
            <a
              href="#apps"
              className={`nav-link ${activeModal === "apps" ? "active" : ""}`}
              onClick={(e) => {
                e.preventDefault();
                openModal("apps");
              }}
            >
              Apps
            </a>
            <a
              href="#platform"
              className={`nav-link ${activeModal === "platform" ? "active" : ""}`}
              onClick={(e) => {
                e.preventDefault();
                openModal("platform");
              }}
            >
              Platform
            </a>
            <a
              href="#about"
              className={`nav-link ${activeModal === "about" ? "active" : ""}`}
              onClick={(e) => {
                e.preventDefault();
                openModal("about");
              }}
            >
              About
            </a>
            <a
              href="#contact"
              className={`nav-link ${activeModal === "contact" ? "active" : ""}`}
              onClick={(e) => {
                e.preventDefault();
                openModal("contact");
              }}
            >
              Contact
            </a>
          </nav>

          <button
            className="burger-btn"
            aria-label="Toggle menu"
            aria-expanded={mobileMenuOpen}
            onClick={() => setMobileMenuOpen(!mobileMenuOpen)}
          >
            <span className="burger-bar" />
            <span className="burger-bar" />
            <span className="burger-bar" />
          </button>
        </header>

        <main className="hero">
          <div className="trust-row anim" style={{ ["--d" as string]: "0.05s" }}>
            <div className="trust-avatars">
              <div className="trust-ring avatar-1">
                <div className="trust-inner">
                  <i className="fa-brands fa-android" />
                </div>
              </div>
              <div className="trust-ring avatar-2">
                <div className="trust-inner">
                  <i className="fa-brands fa-google-play" />
                </div>
              </div>
              <div className="trust-ring avatar-3">
                <div className="trust-inner">
                  <i className="fa-solid fa-bolt" />
                </div>
              </div>
            </div>
            <div className="trust-pill">
              <span>Crafted for Android • Google Play</span>
            </div>
          </div>

          <h1 className="headline">
            <span className="headline-line line-1">ANR STUDIO</span>
            <span className="headline-line line-2">DESIGNED TO EVOLVE</span>
          </h1>

          <p className="subhead anim" style={{ ["--d" as string]: "0.28s" }}>
            Crafting high-quality, privacy-first mobile applications and vibrant
            <br />
            digital experiences powered by the unified ANR Platform architecture.
          </p>

          <div className="cta-wrapper">
            <a
              href="#apps"
              className="cta-btn"
              style={{ ["--d" as string]: "0.4s" }}
              onClick={(e) => {
                e.preventDefault();
                openModal("apps");
              }}
            >
              Explore Ecosystem
            </a>
          </div>
        </main>

        <footer className="stats-footer" ref={footerRef}>
          <div className="stats-grid">
            {STATS_DATA.map((stat, i) => (
              <div
                key={stat.label}
                className="stat-item anim"
                style={{ ["--d" as string]: stat.delay }}
              >
                <span className="stat-icon">{stat.icon}</span>
                <div className="stat-value-wrap">
                  <span className="stat-num">{displayValues[i]}</span>
                  <span className="stat-suffix">{stat.suffix}</span>
                </div>
                <span className="stat-label">{stat.label}</span>
              </div>
            ))}
          </div>
        </footer>
      </div>

      {activeModal && (
        <div
          className="modal-backdrop"
          onClick={() => setActiveModal(null)}
        >
          <div
            className="modal-card"
            onClick={(e) => e.stopPropagation()}
          >
            <div className="modal-header">
              <div>
                <h2 className="modal-title">
                  {activeModal === "apps" && "Ecosystem Applications"}
                  {activeModal === "platform" && "ANR Platform Core"}
                  {activeModal === "about" && "About ANR Studio"}
                  {activeModal === "contact" && "Get In Touch"}
                </h2>
                <p className="modal-subtitle">
                  {activeModal === "apps" && "Explore Android apps engineered by ANR Studio"}
                  {activeModal === "platform" && "Unified high-availability engine powering our ecosystem"}
                  {activeModal === "about" && "Crafting high-performance, aesthetic mobile software"}
                  {activeModal === "contact" && "Connect with our team for inquiries and feedback"}
                </p>
              </div>
              <button
                type="button"
                className="modal-close-btn"
                aria-label="Close"
                onClick={() => setActiveModal(null)}
              >
                ✕
              </button>
            </div>

            <div className="modal-body">
              {activeModal === "apps" && (
                <>
                  <div className="modal-app-card">
                    <div className="modal-app-top">
                      <span className="modal-app-title">ANR Wallpaper</span>
                      <span className="modal-app-tag">Flagship</span>
                    </div>
                    <p className="modal-app-desc">
                      Curated 4K and Ultra-HD wallpapers engineered for vibrant AMOLED displays with fluid category exploration and minimal battery consumption.
                    </p>
                    <div className="modal-badges-row">
                      <span className="modal-badge">4K Ultra HD</span>
                      <span className="modal-badge">Lightweight &lt; 12MB</span>
                      <span className="modal-badge">4.8 ★ Google Play</span>
                    </div>
                    <div className="modal-app-action">
                      <a
                        href="https://play.google.com"
                        target="_blank"
                        rel="noreferrer"
                        className="modal-play-btn"
                      >
                        <i className="fa-brands fa-google-play" /> View on Play Store
                      </a>
                    </div>
                  </div>

                  <div className="modal-app-card">
                    <div className="modal-app-top">
                      <span className="modal-app-title">ANR Focus &amp; Ambient</span>
                      <span className="modal-status-badge">In Development</span>
                    </div>
                    <p className="modal-app-desc">
                      Minimalist digital wellbeing companion designed with offline-first soundscapes, zero intrusive ads, and respectful device telemetry.
                    </p>
                    <div className="modal-badges-row">
                      <span className="modal-badge">Privacy First</span>
                      <span className="modal-badge">Offline Capable</span>
                      <span className="modal-badge">Coming Q4 2026</span>
                    </div>
                  </div>
                </>
              )}

              {activeModal === "platform" && (
                <div className="modal-feature-grid">
                  <div className="modal-feature-item">
                    <h3 className="modal-feature-title">
                      <i className="fa-solid fa-server" /> Go Backend Monolith
                    </h3>
                    <p className="modal-feature-desc">
                      Clean architecture backend delivering sub-50ms API responses with minimal resource overhead.
                    </p>
                  </div>
                  <div className="modal-feature-item">
                    <h3 className="modal-feature-title">
                      <i className="fa-solid fa-bolt" /> Redis Ephemeral Cache
                    </h3>
                    <p className="modal-feature-desc">
                      Distributed in-memory caching and real-time device session synchronization.
                    </p>
                  </div>
                  <div className="modal-feature-item">
                    <h3 className="modal-feature-title">
                      <i className="fa-solid fa-shield-halved" /> First-Party Auth
                    </h3>
                    <p className="modal-feature-desc">
                      Proprietary device auth tokens without dependence on third-party tracking services.
                    </p>
                  </div>
                  <div className="modal-feature-item">
                    <h3 className="modal-feature-title">
                      <i className="fa-solid fa-chart-line" /> 99.9% Uptime Stack
                    </h3>
                    <p className="modal-feature-desc">
                      Containerized Docker deployments monitored 24/7 with Prometheus and Grafana.
                    </p>
                  </div>
                </div>
              )}

              {activeModal === "about" && (
                <>
                  <p className="modal-prose">
                    ANR Studio is a specialized mobile app studio focused on creating fast, refined, and aesthetic Android experiences.
                  </p>
                  <p className="modal-prose">
                    We believe everyday mobile tools should be lightning fast, beautifully designed, and respect user privacy above all else. Every app in our ecosystem is powered by a unified platform built for speed and durability.
                  </p>
                  <div className="modal-badges-row">
                    <span className="modal-badge">Android Specialist</span>
                    <span className="modal-badge">Privacy-Centric</span>
                    <span className="modal-badge">Crafted with Care</span>
                  </div>
                </>
              )}

              {activeModal === "contact" && (
                <>
                  <p className="modal-prose">
                    Have feedback, partnership ideas, or business inquiries? Reach out directly to our team.
                  </p>
                  <div className="modal-contact-row">
                    <span className="modal-email-val">contact@anr-studio.com</span>
                    <button
                      type="button"
                      className="modal-copy-btn"
                      onClick={handleCopyEmail}
                    >
                      {copiedEmail ? "Copied!" : "Copy Email"}
                    </button>
                  </div>
                  <div className="modal-badges-row" style={{ marginTop: "12px" }}>
                    <a
                      href="mailto:contact@anr-studio.com"
                      className="modal-play-btn"
                      style={{ padding: "8px 16px" }}
                    >
                      <i className="fa-regular fa-envelope" /> Send Email
                    </a>
                  </div>
                </>
              )}
            </div>
          </div>
        </div>
      )}
    </>
  );
}
