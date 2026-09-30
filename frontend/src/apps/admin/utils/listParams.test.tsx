import { describe, it, expect } from "vitest";
import { render, screen } from "@testing-library/react";
import { parse, stringify } from "query-string";
import {
  StoreContextProvider,
  TestMemoryRouter,
  memoryStore,
  useListParams,
} from "react-admin";

// react-admin reads and writes the whole list state (page, perPage, sort,
// order, filter) through query-string, which is CommonJS and does
// `require('decode-uri-component')`. Pinning decode-uri-component to an
// ESM-only release (>=0.5.0) turns that require into a namespace object, so
// `parse` throws `TypeError: decodeComponent is not a function` on every list
// render and pagination breaks. Do not add a decode-uri-component override to
// package.json to silence GHSA-vcc3-ghjq-m6fr; these tests cover that.
describe("query-string integration", () => {
  it("exposes callable parse and stringify", () => {
    expect(typeof parse).toBe("function");
    expect(typeof stringify).toBe("function");
  });

  it("parses the react-admin list state from a query string", () => {
    const search =
      "?page=2&perPage=25&sort=name&order=ASC&filter=" +
      encodeURIComponent(JSON.stringify({ q: "Köln" }));

    expect(parse(search)).toEqual({
      page: "2",
      perPage: "25",
      sort: "name",
      order: "ASC",
      filter: JSON.stringify({ q: "Köln" }),
    });
  });

  it("serialises the react-admin list state", () => {
    const encoded = stringify({
      page: 2,
      perPage: 25,
      sort: "name",
      order: "ASC",
      filter: JSON.stringify({ q: "Köln" }),
    });

    expect(encoded).toBe(
      "filter=%7B%22q%22%3A%22K%C3%B6ln%22%7D&order=ASC&page=2&perPage=25&sort=name",
    );
  });

  it("survives malformed percent-encoding", () => {
    expect(() => parse("?q=100%25%20%zz")).not.toThrow();
    expect(parse("?q=100%25%20%zz")).toEqual({ q: "100% %zz" });
  });
});

const ListParamsProbe = () => {
  const [listParams] = useListParams({ resource: "categories" });

  return (
    <div data-testid="list-params">
      {JSON.stringify({
        page: listParams.page,
        perPage: listParams.perPage,
        sort: listParams.sort,
        order: listParams.order,
      })}
    </div>
  );
};

const readProbe = () =>
  JSON.parse(screen.getByTestId("list-params").textContent ?? "{}") as {
    page: number;
    perPage: number;
    sort: string;
    order: string;
  };

const renderProbe = (search: string) =>
  render(
    <StoreContextProvider value={memoryStore()}>
      <TestMemoryRouter initialEntries={[`/categories${search}`]}>
        <ListParamsProbe />
      </TestMemoryRouter>
    </StoreContextProvider>,
  );

describe("useListParams", () => {
  it("reads pagination and sorting from the location", () => {
    renderProbe("?page=2&perPage=25&sort=name&order=ASC");

    expect(readProbe()).toEqual({
      page: 2,
      perPage: 25,
      sort: "name",
      order: "ASC",
    });
  });

  it("reads the filter from the location", () => {
    renderProbe(
      "?page=2&filter=" + encodeURIComponent(JSON.stringify({ q: "Köln" })),
    );

    expect(readProbe().page).toBe(2);
  });

  it("falls back to defaults without a query string", () => {
    renderProbe("");

    const { page } = readProbe();
    expect(page).toBe(1);
  });
});
