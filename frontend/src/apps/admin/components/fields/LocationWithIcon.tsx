import { useRecordContext } from "react-admin";
import PlaceIcon from "@mui/icons-material/Place";
import { Stack, Typography } from "@mui/material";

export const LocationWithIcon = () => {
  const record = useRecordContext();
  if (!record) return null;

  return (
    <Stack direction="row" alignItems="center" spacing={0.5}>
      <PlaceIcon fontSize="small" sx={{ color: "#BF00FF", opacity: 0.8 }} />
      <Typography
        variant="body2"
        sx={{
          color: "text.primary",
          fontWeight: 500,
          letterSpacing: "0.02em",
        }}
      >
        {record.name}
      </Typography>
    </Stack>
  );
};
