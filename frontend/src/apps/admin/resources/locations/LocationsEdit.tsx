import { Edit, SimpleForm, TextInput } from "react-admin";

export const LocationsEdit = () => (
  <Edit title="Edit Location">
    <SimpleForm>
      <TextInput source="id" disabled />
      <TextInput source="name" />
      <TextInput source="address" />
    </SimpleForm>
  </Edit>
);

export default LocationsEdit;
