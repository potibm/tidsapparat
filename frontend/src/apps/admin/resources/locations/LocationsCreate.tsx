import { Create, SimpleForm, TextInput } from "react-admin";

export const LocationsCreate = () => (
  <Create title="Add Location">
    <SimpleForm>
      <TextInput source="name" />
      <TextInput source="address" />
    </SimpleForm>
  </Create>
);

export default LocationsCreate;
