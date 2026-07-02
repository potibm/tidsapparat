import { describe, it, expect, vi, beforeEach } from "vitest";
import { httpClient } from "./dataProvider";

const mockFetchJson = vi.fn();

vi.mock("react-admin", async (importOriginal) => {
  const actual = await importOriginal<typeof import("react-admin")>();
  return {
    ...actual,
    fetchUtils: {
      ...actual.fetchUtils,
      fetchJson: (...args: unknown[]) => mockFetchJson(...args),
    },
  };
});

vi.mock("./authProvider", () => ({
  getAccessToken: vi.fn(),
}));

import { getAccessToken } from "./authProvider";

describe("httpClient", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("adds Bearer token when available", async () => {
    vi.mocked(getAccessToken).mockResolvedValue("test-token");
    mockFetchJson.mockResolvedValue({ json: {}, status: 200 });

    await httpClient("/api/admin/items", { method: "GET" });

    expect(mockFetchJson).toHaveBeenCalledOnce();
    const callArgs = mockFetchJson.mock.calls[0];
    const headers = callArgs[1].headers as Headers;
    expect(headers.get("Authorization")).toBe("Bearer test-token");
  });

  it("does not add Authorization header when no token", async () => {
    vi.mocked(getAccessToken).mockResolvedValue(null);
    mockFetchJson.mockResolvedValue({ json: {}, status: 200 });

    await httpClient("/api/admin/items");

    const callArgs = mockFetchJson.mock.calls[0];
    const headers = callArgs[1].headers as Headers;
    expect(headers.get("Authorization")).toBeNull();
  });

  it("creates default headers if none provided", async () => {
    vi.mocked(getAccessToken).mockResolvedValue(null);
    mockFetchJson.mockResolvedValue({ json: {}, status: 200 });

    await httpClient("/api/admin/items");

    const callArgs = mockFetchJson.mock.calls[0];
    const headers = callArgs[1].headers as Headers;
    expect(headers.get("Accept")).toBe("application/json");
  });

  it("preserves existing headers", async () => {
    vi.mocked(getAccessToken).mockResolvedValue(null);
    mockFetchJson.mockResolvedValue({ json: {}, status: 200 });

    await httpClient("/api/admin/items", {
      headers: new Headers({ "X-Custom": "value" }),
    });

    const callArgs = mockFetchJson.mock.calls[0];
    const headers = callArgs[1].headers as Headers;
    expect(headers.get("X-Custom")).toBe("value");
    // Accept is only added when no headers are provided
    expect(headers.get("Accept")).toBeNull();
  });

  it("converts plain object headers to Headers instance", async () => {
    vi.mocked(getAccessToken).mockResolvedValue(null);
    mockFetchJson.mockResolvedValue({ json: {}, status: 200 });

    await httpClient("/api/admin/items", {
      headers: { "X-Custom": "value" } as unknown as Headers,
    });

    const callArgs = mockFetchJson.mock.calls[0];
    const headers = callArgs[1].headers as Headers;
    expect(headers.get("X-Custom")).toBe("value");
  });
});
