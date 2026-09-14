import Link from "next/link";

export default function NotFound() {
  return (
    <div style={{
      minHeight: "70vh",
      display: "flex",
      flexDirection: "column",
      alignItems: "center",
      justifyContent: "center",
      textAlign: "center",
      padding: "48px 24px",
      fontFamily: "'Inter', system-ui, sans-serif"
    }}>
      <div style={{
        width: "52px",
        height: "52px",
        borderRadius: "50%",
        background: "#000000",
        border: "1px solid rgba(255, 255, 255, 0.18)",
        display: "grid",
        placeItems: "center",
        marginBottom: "24px",
        boxShadow: "0 0 24px rgba(255, 255, 255, 0.12)"
      }}>
        <img
          src="/assets/logo.webp?v=2"
          alt="ANR Studio"
          width={36}
          height={36}
          style={{ width: "76%", height: "76%", objectFit: "contain" }}
        />
      </div>

      <div style={{
        fontFamily: "'BubbledotICG-FinePos', monospace",
        fontSize: "clamp(56px, 10vw, 84px)",
        letterSpacing: "-0.04em",
        lineHeight: 1,
        marginBottom: "14px",
        color: "#ffffff"
      }}>
        404
      </div>

      <h1 style={{
        fontSize: "19px",
        fontWeight: 600,
        letterSpacing: "-0.02em",
        marginBottom: "8px",
        color: "#ffffff"
      }}>
        Console Route Not Found
      </h1>

      <p style={{
        fontSize: "13.5px",
        color: "#888888",
        maxWidth: "380px",
        lineHeight: 1.6,
        marginBottom: "28px"
      }}>
        The administrative panel or resource endpoint you requested is not available.
      </p>

      <Link
        href="/"
        style={{
          display: "inline-flex",
          alignItems: "center",
          gap: "8px",
          background: "#ffffff",
          color: "#000000",
          fontSize: "13.5px",
          fontWeight: 600,
          padding: "11px 24px",
          borderRadius: "999px",
          textDecoration: "none",
          boxShadow: "0 0 20px rgba(255, 255, 255, 0.2)",
          transition: "transform 0.2s ease"
        }}
      >
        Return to Dashboard
      </Link>
    </div>
  );
}
