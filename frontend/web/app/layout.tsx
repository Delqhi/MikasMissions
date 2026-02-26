import type { Metadata } from "next";
import { Baloo_2, Nunito } from "next/font/google";
import type { ReactNode } from "react";
import { getLocaleFromRequest } from "../lib/i18n";
import "./globals.css";

const bodyFont = Nunito({
  subsets: ["latin"],
  variable: "--font-body-next"
});

const displayFont = Baloo_2({
  subsets: ["latin"],
  variable: "--font-display-next",
  weight: ["500", "700", "800"]
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

export default async function RootLayout({ children }: RootLayoutProps) {
  const locale = await getLocaleFromRequest();

  return (
    <html className={`${bodyFont.variable} ${displayFont.variable}`} lang={locale}>
      <body>{children}</body>
    </html>
  );
}
