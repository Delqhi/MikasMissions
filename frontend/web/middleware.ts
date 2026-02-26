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

    return response;
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

  return response;
}

export const config = {
  matcher: [
    "/((?!api|_next/static|_next/image|_next/data|favicon.ico|robots.txt|sitemap.xml|manifest.webmanifest|icon.svg|apple-icon.svg|.*\\..*).*)"
  ]
};
