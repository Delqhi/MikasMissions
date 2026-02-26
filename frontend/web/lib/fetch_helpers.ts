export function apiBaseURL(): string {
  return process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080";
}

export function useFallbackData(): boolean {
  if (process.env.NODE_ENV === "production") {
    return false;
  }
  return process.env.NEXT_PUBLIC_USE_API_FALLBACKS === "true";
}

type JSONFetchOptions = {
  body?: unknown;
  headers?: Record<string, string>;
  method?: "GET" | "POST" | "PUT";
  token?: string;
};

export async function safeJSONFetch<T>(input: string, options: JSONFetchOptions = {}): Promise<T> {
  const headers: Record<string, string> = {
    Accept: "application/json",
    ...options.headers
  };
  if (options.body !== undefined) {
    headers["Content-Type"] = "application/json";
  }
  if (options.token) {
    headers.Authorization = `Bearer ${options.token}`;
  }
  try {
    const response = await fetch(input, {
      method: options.method ?? "GET",
      cache: "no-store",
      headers,
      body: options.body === undefined ? undefined : JSON.stringify(options.body)
    });

    if (!response.ok) {
      throw new Error(`request failed with status ${response.status}`);
    }

    return (await response.json()) as T;
  } catch (error) {
    throw new Error(`fetch failed for ${input}: ${String(error)}`);
  }
}
