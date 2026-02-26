import type { NextRequest } from "next/server";
import { NextResponse } from "next/server";
import {
  i18nRequestHeader,
  isLocale,
  localeCookieName,
  resolvePreferredLocale,
  stripLocalePrefix,
  withLocalePath
} from "./lib/i18n";

const localeCookieMaxAgeSeconds = 60 * 60 * 24 * 365;
const contentSecurityPolicy = [
  "default-src 'self'",
  "base-uri 'self'",
  "font-src 'self' https://fonts.gstatic.com data:",
  "img-src 'self' data: https:",
  "object-src 'none'",
  "script-src 'self' 'unsafe-inline'",
  "style-src 'self' 'unsafe-inline' https://fonts.googleapis.com",
  "connect-src 'self' https:",
  "frame-ancestors 'none'",
  "form-action 'self'"
].join("; ");

const securityHeaders = [
  ["Content-Security-Policy", contentSecurityPolicy],
  ["Cross-Origin-Opener-Policy", "same-origin"],
  ["Cross-Origin-Resource-Policy", "same-origin"],
  ["Origin-Agent-Cluster", "?1"],
  ["Permissions-Policy", "camera=(), microphone=(), geolocation=(), payment=()"],
  ["Referrer-Policy", "strict-origin-when-cross-origin"],
  ["Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload"],
  ["X-Content-Type-Options", "nosniff"],
  ["X-DNS-Prefetch-Control", "off"],
  ["X-Frame-Options", "DENY"],
  ["X-Permitted-Cross-Domain-Policies", "none"]
] as const;

function withSecurityHeaders(response: NextResponse): NextResponse {
  for (const [key, value] of securityHeaders) {
    response.headers.set(key, value);
  }
  return response;
}

export function middleware(request: NextRequest) {
  const { pathname } = request.nextUrl;
  const localeFromPath = pathname.split("/")[1] ?? "";

  if (isLocale(localeFromPath)) {
    const rewrittenURL = request.nextUrl.clone();
    rewrittenURL.pathname = stripLocalePrefix(pathname);

    const requestHeaders = new Headers(request.headers);
    requestHeaders.set(i18nRequestHeader, localeFromPath);

    const response = NextResponse.rewrite(rewrittenURL, {
      request: {
        headers: requestHeaders
      }
    });

    response.cookies.set(localeCookieName, localeFromPath, {
      maxAge: localeCookieMaxAgeSeconds,
      path: "/",
      sameSite: "lax"
    });

    return withSecurityHeaders(response);
  }

  const localeFromRewrite = request.headers.get(i18nRequestHeader);
  if (isLocale(localeFromRewrite)) {
    const response = NextResponse.next();
    response.cookies.set(localeCookieName, localeFromRewrite, {
      maxAge: localeCookieMaxAgeSeconds,
      path: "/",
      sameSite: "lax"
    });
    return withSecurityHeaders(response);
  }

  const preferredLocale = resolvePreferredLocale(
    request.cookies.get(localeCookieName)?.value,
    request.headers.get("accept-language")
  );

  const redirectURL = request.nextUrl.clone();
  redirectURL.pathname = withLocalePath(preferredLocale, pathname);

  const response = NextResponse.redirect(redirectURL);
  response.cookies.set(localeCookieName, preferredLocale, {
    maxAge: localeCookieMaxAgeSeconds,
    path: "/",
    sameSite: "lax"
  });

  return withSecurityHeaders(response);
}

export const config = {
  matcher: [
    "/((?!api|v1|_next/static|_next/image|_next/data|favicon.ico|robots.txt|sitemap.xml|manifest.webmanifest|icon.svg|apple-icon.svg|.*\\..*).*)"
  ]
};
