import { Create, SimpleForm, TextInput } from "react-admin";

export const CategoriesCreate = () => (
  <Create title="Add Category">
    <SimpleForm defaultValues={{ color: "#BF00FF" }}>
      <TextInput source="name" />
      <TextInput source="color" type="color" />
    </SimpleForm>
  </Create>
);

export default CategoriesCreate;
