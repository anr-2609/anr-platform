import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "ANR Studio - Modern Mobile Applications",
  description: "ANR Studio builds high performance, user-centric mobile applications.",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  );
}
