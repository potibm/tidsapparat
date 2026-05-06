import {
  List,
  DatagridConfigurable,
  TextField,
  FunctionField,
  TopToolbar,
  CreateButton,
  ExportButton,
  SelectColumnsButton,
  SearchInput,
} from "react-admin";
import { LocationRecord } from "./location.types";

const locationFilters = [
  <SearchInput key="q" source="q" alwaysOn placeholder="Search by name..." />,
];

const ListActions = () => (
  <TopToolbar>
    <SelectColumnsButton />
    <CreateButton />
    <ExportButton />
  </TopToolbar>
);

export const LocationsList = () => (
  <List
    sort={{ field: "id", order: "DESC" }}
    actions={<ListActions />}
    filters={locationFilters}
  >
    <DatagridConfigurable rowClick="edit" bulkActionButtons={false}>
      <TextField source="id" />
      <TextField source="name" />
      <FunctionField
        source="address"
        label="Address"
        render={(record: LocationRecord) => record.address ?? "—"}
      />
    </DatagridConfigurable>
  </List>
);

export default LocationsList;
