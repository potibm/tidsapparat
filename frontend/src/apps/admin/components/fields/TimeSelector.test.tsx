import { describe, it, expect, vi, beforeEach, type Mock } from "vitest";
import { render, screen, fireEvent } from "@testing-library/react";
import "@testing-library/jest-dom/vitest";
import { TimeSelector } from "./TimeSelector";
import { ConfigContext } from "@core/config/ConfigContext";
import { AppConfig } from "@core/config/config.schemas";
import { FormProvider, useForm } from "react-hook-form";
import dayjs from "dayjs";
import utc from "dayjs/plugin/utc";
import timezone from "dayjs/plugin/timezone";

dayjs.extend(utc);
dayjs.extend(timezone);

const mockSetValue = vi.fn();

vi.mock("react-admin", async () => {
  const actual =
    await vi.importActual<typeof import("react-admin")>("react-admin");
  return {
    ...actual,
    DateInput: vi.fn(({ source, label, ...props }) => (
      <input
        data-testid={`date-input-${source}`}
        aria-label={label}
        {...props}
      />
    )),
    TimeInput: vi.fn(({ source, label, ...props }) => (
      <input
        data-testid={`time-input-${source}`}
        aria-label={label}
        {...props}
      />
    )),
    NumberInput: vi.fn(({ source, label, ...props }) => (
      <input
        data-testid={`number-input-${source}`}
        aria-label={label}
        {...props}
      />
    )),
    useRecordContext: vi.fn(),
  };
});

vi.mock("react-hook-form", async () => {
  const actual =
    await vi.importActual<typeof import("react-hook-form")>("react-hook-form");
  return {
    ...actual,
    useWatch: vi.fn(),
    useFormContext: vi.fn(),
  };
});

import { useWatch, useFormContext } from "react-hook-form";
import { useRecordContext } from "react-admin";

const defaultConfig: AppConfig = {
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
  date_options: {
    weekday: "long",
    hour: "2-digit",
    minute: "2-digit",
  },
  timezone: "Europe/Berlin",
  party_days: [],
  event_durations: [15, 30, 60, 90],
};

function TestWrapper({
  children,
  config = defaultConfig,
}: {
  children: React.ReactNode;
  config?: AppConfig;
}) {
  const methods = useForm();
  return (
    <ConfigContext value={config}>
      <FormProvider {...methods}>{children}</FormProvider>
    </ConfigContext>
  );
}

const mockPartyDays = [
  { id: "2024-06-15", name: "Saturday" },
  { id: "2024-06-16", name: "Sunday" },
];

const mockedUseWatch = useWatch as Mock;
const mockedUseRecordContext = useRecordContext as Mock;

describe("TimeSelector", () => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(useFormContext).mockReturnValue({
      setValue: mockSetValue,
    } as unknown as ReturnType<typeof useFormContext>);
  });

  it("renders date, time and duration inputs", () => {
    mockedUseWatch.mockReturnValue(undefined);
    mockedUseRecordContext.mockReturnValue(undefined);

    render(
      <TestWrapper>
        <TimeSelector partyDays={mockPartyDays} />
      </TestWrapper>,
    );

    expect(screen.getByTestId("date-input-party_day")).toBeInTheDocument();
    expect(
      screen.getByTestId("time-input-start_time_only"),
    ).toBeInTheDocument();
    expect(
      screen.getByTestId("number-input-duration_mins"),
    ).toBeInTheDocument();
  });

  it("renders party day chips", () => {
    mockedUseWatch.mockImplementation(({ name }: { name: string }) => {
      if (name === "party_day") return undefined;
      if (name === "start_time_only") return undefined;
      if (name === "duration_mins") return undefined;
      return undefined;
    });
    mockedUseRecordContext.mockReturnValue(undefined);

    render(
      <TestWrapper>
        <TimeSelector partyDays={mockPartyDays} />
      </TestWrapper>,
    );

    expect(screen.getByText("Saturday")).toBeInTheDocument();
    expect(screen.getByText("Sunday")).toBeInTheDocument();
  });

  it("highlights the active party day chip", () => {
    mockedUseWatch.mockImplementation(({ name }: { name: string }) => {
      if (name === "party_day") return "2024-06-15";
      if (name === "start_time_only") return undefined;
      if (name === "duration_mins") return undefined;
      return undefined;
    });
    mockedUseRecordContext.mockReturnValue(undefined);

    render(
      <TestWrapper>
        <TimeSelector partyDays={mockPartyDays} />
      </TestWrapper>,
    );

    const saturdayChip = screen.getByText("Saturday").closest(".MuiChip-root");
    expect(saturdayChip).toHaveClass("MuiChip-filled");
    expect(saturdayChip).toHaveClass("MuiChip-colorPrimary");

    const sundayChip = screen.getByText("Sunday").closest(".MuiChip-root");
    expect(sundayChip).toHaveClass("MuiChip-outlined");
  });

  it("clicking a party day chip calls setValue with the day id", () => {
    mockedUseWatch.mockImplementation(({ name }: { name: string }) => {
      if (name === "party_day") return undefined;
      if (name === "start_time_only") return undefined;
      if (name === "duration_mins") return undefined;
      return undefined;
    });
    mockedUseRecordContext.mockReturnValue(undefined);

    render(
      <TestWrapper>
        <TimeSelector partyDays={mockPartyDays} />
      </TestWrapper>,
    );

    fireEvent.click(screen.getByText("Sunday"));

    expect(mockSetValue).toHaveBeenCalledWith(
      "party_day",
      "2024-06-16",
      expect.objectContaining({
        shouldDirty: true,
        shouldValidate: true,
        shouldTouch: true,
      }),
    );
  });

  it("renders duration chips with default durations", () => {
    mockedUseWatch.mockImplementation(({ name }: { name: string }) => {
      if (name === "party_day") return undefined;
      if (name === "start_time_only") return undefined;
      if (name === "duration_mins") return undefined;
      return undefined;
    });
    mockedUseRecordContext.mockReturnValue(undefined);

    render(
      <TestWrapper>
        <TimeSelector partyDays={mockPartyDays} />
      </TestWrapper>,
    );

    expect(screen.getByText("15m")).toBeInTheDocument();
    expect(screen.getByText("30m")).toBeInTheDocument();
    expect(screen.getByText("60m")).toBeInTheDocument();
    expect(screen.getByText("90m")).toBeInTheDocument();
  });

  it("clicking a duration chip calls setValue with the minutes", () => {
    mockedUseWatch.mockImplementation(({ name }: { name: string }) => {
      if (name === "party_day") return undefined;
      if (name === "start_time_only") return undefined;
      if (name === "duration_mins") return undefined;
      return undefined;
    });
    mockedUseRecordContext.mockReturnValue(undefined);

    render(
      <TestWrapper>
        <TimeSelector partyDays={mockPartyDays} />
      </TestWrapper>,
    );

    fireEvent.click(screen.getByText("30m"));

    expect(mockSetValue).toHaveBeenCalledWith(
      "duration_mins",
      30,
      expect.objectContaining({
        shouldDirty: true,
        shouldValidate: true,
        shouldTouch: true,
      }),
    );
  });

  it("shows preview when all fields are filled", () => {
    mockedUseWatch.mockImplementation(({ name }: { name: string }) => {
      if (name === "party_day") return "2024-06-15";
      if (name === "start_time_only") return "14:30";
      if (name === "duration_mins") return 60;
      return undefined;
    });
    mockedUseRecordContext.mockReturnValue(undefined);

    render(
      <TestWrapper>
        <TimeSelector partyDays={mockPartyDays} />
      </TestWrapper>,
    );

    expect(screen.getByText(/Event ends on:/)).toBeInTheDocument();
  });

  it("shows placeholder text when fields are incomplete", () => {
    mockedUseWatch.mockImplementation(({ name }: { name: string }) => {
      if (name === "party_day") return undefined;
      if (name === "start_time_only") return "14:30";
      if (name === "duration_mins") return 60;
      return undefined;
    });
    mockedUseRecordContext.mockReturnValue(undefined);

    render(
      <TestWrapper>
        <TimeSelector partyDays={mockPartyDays} />
      </TestWrapper>,
    );

    expect(
      screen.getByText("Please fill out the three fields."),
    ).toBeInTheDocument();
  });

  it("renders custom preset durations when provided", () => {
    mockedUseWatch.mockImplementation(({ name }: { name: string }) => {
      if (name === "party_day") return undefined;
      if (name === "start_time_only") return undefined;
      if (name === "duration_mins") return undefined;
      return undefined;
    });
    mockedUseRecordContext.mockReturnValue(undefined);

    render(
      <TestWrapper>
        <TimeSelector partyDays={mockPartyDays} presetDurations={[5, 10, 15]} />
      </TestWrapper>,
    );

    expect(screen.getByText("5m")).toBeInTheDocument();
    expect(screen.getByText("10m")).toBeInTheDocument();
    expect(screen.getByText("15m")).toBeInTheDocument();
    expect(screen.queryByText("30m")).not.toBeInTheDocument();
  });
});
