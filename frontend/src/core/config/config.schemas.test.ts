import { describe, it, expect } from "vitest";
import { AppConfigSchema } from "./config.schemas";

describe("AppConfigSchema", () => {
  const validConfig = {
    version: "1.0.0",
    environment: "test",
    sentry: {
      dsn: "https://test@sentry.io/123",
      environment: "test",
      version: "1.0.0",
      replay_session_sample_rate: 0,
      replay_error_sample_rate: 1,
    },
    date_locale: "en-US",
    date_options: {
      weekday: "long",
      hour: "2-digit",
      minute: "2-digit",
    },
    timezone: "Europe/Berlin",
    party_days: [{ id: "2024-06-15", name: "Saturday" }],
    event_durations: [15, 30, 60],
  };

  it("parses a valid full config", () => {
    const result = AppConfigSchema.parse(validConfig);
    expect(result.version).toBe("1.0.0");
    expect(result.party_days).toHaveLength(1);
    expect(result.event_durations).toEqual([15, 30, 60]);
  });

  it("parses a valid minimal config without optional fields", () => {
    const minimal = {
      version: "1.0.0",
      environment: "test",
      sentry: {
        dsn: "",
        environment: "test",
        version: "1.0.0",
      },
      date_locale: "da-DK",
      timezone: "Europe/Copenhagen",
      party_days: [],
    };

    const result = AppConfigSchema.parse(minimal);
    expect(result.date_options).toEqual({
      weekday: "long",
      hour: "2-digit",
      minute: "2-digit",
    });
    expect(result.event_durations).toEqual([0, 5, 10, 15, 30, 45, 60, 90, 120]);
  });

  it("fails when required fields are missing", () => {
    expect(() =>
      AppConfigSchema.parse({
        version: "1.0.0",
      }),
    ).toThrow();
  });

  it("fails when date_locale is too short", () => {
    expect(() =>
      AppConfigSchema.parse({
        ...validConfig,
        date_locale: "e",
      }),
    ).toThrow();
  });

  it("fails when sentry sample rates are out of range", () => {
    expect(() =>
      AppConfigSchema.parse({
        ...validConfig,
        sentry: {
          ...validConfig.sentry,
          replay_session_sample_rate: 2,
        },
      }),
    ).toThrow();
  });

  it("fails when auth authority is not a valid URL", () => {
    expect(() =>
      AppConfigSchema.parse({
        ...validConfig,
        auth: {
          type: "oidc",
          name: "dex",
          authority: "not-a-url",
          client_id: "client",
        },
      }),
    ).toThrow();
  });

  it("parses config with valid auth", () => {
    const result = AppConfigSchema.parse({
      ...validConfig,
      auth: {
        type: "oidc",
        name: "dex",
        authority: "https://dex.example.com",
        client_id: "react-admin-client",
      },
    });

    expect(result.auth).toBeDefined();
    expect(result.auth?.type).toBe("oidc");
  });

  it("accepts empty party_days array", () => {
    const result = AppConfigSchema.parse({
      ...validConfig,
      party_days: [],
    });

    expect(result.party_days).toEqual([]);
  });

  it("accepts custom event durations", () => {
    const result = AppConfigSchema.parse({
      ...validConfig,
      event_durations: [5, 10],
    });

    expect(result.event_durations).toEqual([5, 10]);
  });
});
