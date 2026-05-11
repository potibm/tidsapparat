import { Create, SimpleForm, TextInput } from "react-admin";
import { useFormContext } from "react-hook-form";
import { InputAdornment, IconButton } from "@mui/material";
import AutorenewIcon from "@mui/icons-material/Autorenew";
import randomColor from "randomcolor";

const getRandomColor = () => {
  return randomColor({
    luminosity: "bright",
    format: "hex",
  });
};

const RandomColorButton = () => {
  const { setValue } = useFormContext();

  const handleRandomize = () => {
    setValue("color", getRandomColor(), { shouldDirty: true });
  };

  return (
    <InputAdornment position="end">
      <IconButton onClick={handleRandomize} title="Random Color">
        <AutorenewIcon />
      </IconButton>
    </InputAdornment>
  );
};

export const CategoriesCreate = () => (
  <Create title="Add Category">
    <SimpleForm defaultValues={{ color: getRandomColor() }}>
      <TextInput source="name" />
      <TextInput
        source="color"
        type="color"
        helperText="Pick a color or click the randomize button"
        slotProps={{
          input: {
            endAdornment: <RandomColorButton />,
          },
        }}
      />
    </SimpleForm>
  </Create>
);

export default CategoriesCreate;
