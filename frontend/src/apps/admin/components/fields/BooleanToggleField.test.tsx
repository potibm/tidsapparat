import { describe, it, expect, vi } from "vitest";
import { render, screen, fireEvent } from "@testing-library/react";
import { BooleanToggleField } from "./BooleanToggleField";
import { RecordContextProvider, ResourceContextProvider } from "react-admin";

const mockUpdate = vi.fn();
const mockNotify = vi.fn();

vi.mock("react-admin", async () => {
  const actual =
    await vi.importActual<typeof import("react-admin")>("react-admin");
  return {
    ...actual,
    useUpdate: () => [mockUpdate, { isLoading: false }],
    useNotify: () => mockNotify,
  };
});

describe("BooleanToggleField", () => {
  it("renders switch based on record value", () => {
    render(
      <ResourceContextProvider value="schedule-entries">
        <RecordContextProvider value={{ id: 1, hidden: true }}>
          <BooleanToggleField source="hidden" />
        </RecordContextProvider>
      </ResourceContextProvider>,
    );

    const switchEl = screen.getByRole("switch");
    expect(switchEl).toBeChecked();
  });

  it("returns null when no record", () => {
    const { container } = render(
      <ResourceContextProvider value="schedule-entries">
        <RecordContextProvider value={undefined}>
          <BooleanToggleField source="hidden" />
        </RecordContextProvider>
      </ResourceContextProvider>,
    );

    expect(container.firstChild).toBeNull();
  });

  it("calls update on change", () => {
    render(
      <ResourceContextProvider value="schedule-entries">
        <RecordContextProvider value={{ id: 1, hidden: false }}>
          <BooleanToggleField source="hidden" />
        </RecordContextProvider>
      </ResourceContextProvider>,
    );

    const switchEl = screen.getByRole("switch");
    fireEvent.click(switchEl);

    expect(mockUpdate).toHaveBeenCalled();
  });
});
