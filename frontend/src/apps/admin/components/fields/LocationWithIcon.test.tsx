import { describe, it, expect } from "vitest";
import { render, screen } from "@testing-library/react";
import { LocationWithIcon } from "./LocationWithIcon";
import { RecordContextProvider } from "react-admin";

describe("LocationWithIcon", () => {
  it("renders location name with icon", () => {
    render(
      <RecordContextProvider value={{ id: 1, name: "Main Hall" }}>
        <LocationWithIcon />
      </RecordContextProvider>,
    );

    expect(screen.getByText("Main Hall")).toBeInTheDocument();
  });

  it("returns null when no record", () => {
    const { container } = render(
      <RecordContextProvider value={undefined}>
        <LocationWithIcon />
      </RecordContextProvider>,
    );

    expect(container.firstChild).toBeNull();
  });
});
