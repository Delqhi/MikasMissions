import { cookies, headers } from "next/headers";
import { deMessages } from "./messages/de";
import { enMessages } from "./messages/en";
import { esMessages } from "./messages/es";
import type { LocalizedMessages } from "./messages/types";

export const supportedLocales = ["de", "en", "es"] as const;
export type Locale = (typeof supportedLocales)[number];

export const defaultLocale: Locale = "de";
export const localeCookieName = "mm_locale";
const localeHeaderName = "x-mm-locale";

const messagesByLocale: Record<Locale, LocalizedMessages> = {
  de: deMessages,
  en: enMessages,
  es: esMessages
};

export function isLocale(value: string | null | undefined): value is Locale {
  return value === "de" || value === "en" || value === "es";
}

export function localeFromPathname(pathname: string): Locale | null {
  const [, segment] = pathname.split("/");
  return isLocale(segment) ? segment : null;
}

export function stripLocalePrefix(pathname: string): string {
  const locale = localeFromPathname(pathname);
  if (!locale) {
    return pathname;
  }
  const stripped = pathname.slice(locale.length + 1);
  return stripped === "" ? "/" : stripped;
}

function normalizePath(pathname: string): string {
  if (pathname === "") {
    return "/";
  }
  return pathname.startsWith("/") ? pathname : `/${pathname}`;
}

export function withLocalePath(locale: Locale, pathname: string): string {
  const normalized = normalizePath(pathname);
  const prefixLocale = localeFromPathname(normalized);

  if (prefixLocale) {
    const withoutPrefix = stripLocalePrefix(normalized);
    return withoutPrefix === "/" ? `/${locale}` : `/${locale}${withoutPrefix}`;
  }

  return normalized === "/" ? `/${locale}` : `/${locale}${normalized}`;
}

export function preferredLocaleFromAcceptLanguage(acceptLanguage: string | null | undefined): Locale {
  if (!acceptLanguage) {
    return defaultLocale;
  }

  const rawCandidates = acceptLanguage.split(",").map((value) => value.trim().toLowerCase());

  for (const candidate of rawCandidates) {
    const language = candidate.split(";")[0]?.trim() ?? "";
    if (isLocale(language)) {
      return language;
    }
    const primary = language.split("-")[0];
    if (isLocale(primary)) {
      return primary;
    }
  }

  return defaultLocale;
}

export function resolvePreferredLocale(cookieLocale: string | null | undefined, acceptLanguage: string | null | undefined): Locale {
  if (isLocale(cookieLocale)) {
    return cookieLocale;
  }

  return preferredLocaleFromAcceptLanguage(acceptLanguage);
}

export async function getLocaleFromRequest(): Promise<Locale> {
  const [cookieStore, headerStore] = await Promise.all([cookies(), headers()]);
  const headerLocale = headerStore.get(localeHeaderName);

  if (isLocale(headerLocale)) {
    return headerLocale;
  }

  return resolvePreferredLocale(cookieStore.get(localeCookieName)?.value, headerStore.get("accept-language"));
}

export async function getMessages(locale?: Locale): Promise<LocalizedMessages> {
  const resolvedLocale = locale ?? (await getLocaleFromRequest());
  return messagesByLocale[resolvedLocale];
}

export async function getLocaleAndMessages(): Promise<{ locale: Locale; messages: LocalizedMessages }> {
  const locale = await getLocaleFromRequest();
  return { locale, messages: messagesByLocale[locale] };
}

export const i18nRequestHeader = localeHeaderName;
