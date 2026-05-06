import { describe, it, expect } from "vitest";
import { transformScheduleToAPI } from "./time";

describe("transformScheduleToAPI", () => {
  it("computes start_time and end_time ISO strings", () => {
    const result = transformScheduleToAPI(
      {
        party_day: "2024-06-15",
        start_time_only: "14:30",
        duration_mins: 90,
      },
      "Europe/London",
    );

    expect(result.start_time).toMatch(/^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}/);
    expect(result.end_time).toMatch(/^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}/);
  });

  it("offsets end_time by duration_mins from start_time", () => {
    const result = transformScheduleToAPI(
      {
        party_day: "2024-06-15",
        start_time_only: "10:00",
        duration_mins: 45,
      },
      "Europe/London",
    );

    const start = new Date(result.start_time);
    const end = new Date(result.end_time);
    const diffMinutes = (end.getTime() - start.getTime()) / 60000;

    expect(diffMinutes).toBe(45);
  });

  it("clears internal fields", () => {
    const result = transformScheduleToAPI(
      {
        party_day: "2024-06-15",
        start_time_only: "08:00",
        duration_mins: 30,
      },
      "Europe/London",
    );

    expect(result.party_day).toBeUndefined();
    expect(result.start_time_only).toBeUndefined();
    expect(result.duration_mins).toBeUndefined();
  });

  it("preserves extra properties from the input", () => {
    const result = transformScheduleToAPI(
      {
        party_day: "2024-06-15",
        start_time_only: "20:00",
        duration_mins: 60,
        title: "Main Event",
        location: "Hall A",
      } as {
        party_day: string;
        start_time_only: string;
        duration_mins: number;
        title: string;
        location: string;
      },
      "Europe/London",
    );

    expect((result as Record<string, unknown>).title).toBe("Main Event");
    expect((result as Record<string, unknown>).location).toBe("Hall A");
  });
});
