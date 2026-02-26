import type { Metadata } from "next";
import { Fredoka, Nunito } from "next/font/google";
import type { ReactNode } from "react";
import "./globals.css";

const bodyFont = Nunito({
  subsets: ["latin"],
  variable: "--font-body-next"
});

const displayFont = Fredoka({
  subsets: ["latin"],
  variable: "--font-display-next"
});

export const metadata: Metadata = {
  title: "MikasMissions",
  description: "Age-aware kids streaming with strict safety defaults",
  applicationName: "MikasMissions",
  manifest: "/manifest.webmanifest",
  appleWebApp: {
    capable: true,
    statusBarStyle: "default",
    title: "MikasMissions"
  },
  icons: {
    apple: "/apple-icon.svg",
    icon: "/icon.svg"
  }
};

type RootLayoutProps = {
  children: ReactNode;
};

export default function RootLayout({ children }: RootLayoutProps) {
  return (
    <html className={`${bodyFont.variable} ${displayFont.variable}`} lang="de">
      <body>{children}</body>
    </html>
  );
}
