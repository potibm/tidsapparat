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
import { Box } from "@mui/material";
import { CategoryRecord } from "./category.types";

const categoryFilters = [
  <SearchInput key="q" source="q" alwaysOn placeholder="Search by name..." />,
];

const ListActions = () => (
  <TopToolbar>
    <SelectColumnsButton />
    <CreateButton />
    <ExportButton />
  </TopToolbar>
);

export const CategoriesList = () => (
  <List
    sort={{ field: "id", order: "DESC" }}
    actions={<ListActions />}
    filters={categoryFilters}
  >
    <DatagridConfigurable rowClick="edit" bulkActionButtons={false}>
      <TextField source="id" />
      <TextField source="name" />
      <FunctionField
        source="color"
        label="Color"
        render={(record: CategoryRecord) => (
          <Box sx={{ display: "flex", alignItems: "center", gap: 1 }}>
            <Box
              sx={{
                width: 16,
                height: 16,
                borderRadius: "50%",
                backgroundColor: record.color,
                border: "1px solid rgba(0,0,0,0.12)",
              }}
            />
            <span>{record.color}</span>
          </Box>
        )}
      />
    </DatagridConfigurable>
  </List>
);

export default CategoriesList;
