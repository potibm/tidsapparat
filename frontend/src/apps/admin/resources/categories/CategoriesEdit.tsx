import { Edit, SimpleForm, TextInput } from "react-admin";

export const CategoriesEdit = () => (
  <Edit title="Edit Category">
    <SimpleForm>
      <TextInput source="id" disabled />
      <TextInput source="name" />
      <TextInput source="color" type="color" />
    </SimpleForm>
  </Edit>
);

export default CategoriesEdit;
