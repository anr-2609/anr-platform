import Link from "next/link";

export default function NotFound() {
  return (
    <div style={{
      minHeight: "100vh",
      display: "flex",
      flexDirection: "column",
      alignItems: "center",
      justifyContent: "center",
      background: "#000000",
      backgroundImage: "radial-gradient(circle at 50% 30%, rgba(255, 255, 255, 0.05) 0%, rgba(0, 0, 0, 0.98) 70%)",
      color: "#ffffff",
      padding: "24px",
      textAlign: "center",
      fontFamily: "'Inter', system-ui, sans-serif"
    }}>
      <div style={{
        width: "56px",
        height: "56px",
        borderRadius: "50%",
        background: "#000000",
        border: "1px solid rgba(255, 255, 255, 0.18)",
        display: "grid",
        placeItems: "center",
        marginBottom: "28px",
        boxShadow: "0 0 24px rgba(255, 255, 255, 0.12)"
      }}>
        <img
          src="/assets/logo.webp?v=2"
          alt="ANR Studio"
          width={40}
          height={40}
          style={{ width: "76%", height: "76%", objectFit: "contain" }}
        />
      </div>

      <div style={{
        fontFamily: "'BubbledotICG-FinePos', monospace",
        fontSize: "clamp(64px, 12vw, 96px)",
        letterSpacing: "-0.04em",
        lineHeight: 1,
        marginBottom: "16px",
        color: "#ffffff"
      }}>
        404
      </div>

      <h1 style={{
        fontSize: "20px",
        fontWeight: 600,
        letterSpacing: "-0.02em",
        marginBottom: "10px",
        color: "#ffffff"
      }}>
        Page Not Found
      </h1>

      <p style={{
        fontSize: "14px",
        color: "#888888",
        maxWidth: "420px",
        lineHeight: 1.6,
        marginBottom: "32px"
      }}>
        The page you are looking for does not exist or has been moved within the ANR Studio platform.
      </p>

      <Link
        href="/"
        style={{
          display: "inline-flex",
          alignItems: "center",
          gap: "8px",
          background: "#ffffff",
          color: "#000000",
          fontSize: "14px",
          fontWeight: 600,
          padding: "12px 28px",
          borderRadius: "999px",
          textDecoration: "none",
          boxShadow: "0 0 24px rgba(255, 255, 255, 0.25)",
          transition: "transform 0.2s ease"
        }}
      >
        Return to Homepage
      </Link>
    </div>
  );
}
