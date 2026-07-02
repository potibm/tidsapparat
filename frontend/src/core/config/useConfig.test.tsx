import { describe, it, expect } from "vitest";
import { renderHook } from "@testing-library/react";
import { useAppConfig } from "./useConfig";
import { ConfigContext } from "./ConfigContext";
import type { AppConfig } from "./config.schemas";

const mockConfig: AppConfig = {
  version: "1.0.0",
  environment: "test",
  sentry: {
    dsn: "",
    environment: "test",
    version: "1.0.0",
    replay_session_sample_rate: 0,
    replay_error_sample_rate: 1,
  },
  date_locale: "en-US",
  date_options: { weekday: "long", hour: "2-digit", minute: "2-digit" },
  timezone: "Europe/Berlin",
  party_days: [{ id: "2024-06-15", name: "Saturday" }],
  event_durations: [15, 30, 60],
};

function wrapper({ children }: { children: React.ReactNode }) {
  return <ConfigContext value={mockConfig}>{children}</ConfigContext>;
}

describe("useAppConfig", () => {
  it("returns config when inside provider", () => {
    const { result } = renderHook(() => useAppConfig(), { wrapper });

    expect(result.current.version).toBe("1.0.0");
    expect(result.current.timezone).toBe("Europe/Berlin");
    expect(result.current.party_days).toHaveLength(1);
  });

  it("throws when used outside provider", () => {
    expect(() => renderHook(() => useAppConfig())).toThrow(
      "useAppConfig must be used within a ConfigContext.Provider",
    );
  });
});
