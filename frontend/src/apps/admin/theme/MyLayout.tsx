import {
  AppBar,
  AppBarProps,
  TitlePortal,
  Layout,
  LayoutProps,
  Menu,
  MenuProps,
} from "react-admin";
import { Typography, Box } from "@mui/material";
import { Logo } from "@core/logo/Logo";
import { useAppConfig } from "@core/config/useConfig";

export const MyAppBar = (props: AppBarProps) => (
  <AppBar {...props} color="secondary">
    <Box sx={{ flex: "1", display: "flex", alignItems: "center" }}>
      <Logo
        style={{
          height: "32px",
          width: "32px",
          marginRight: "12px",
        }}
      />

      <Typography
        variant="h6"
        color="inherit"
        sx={{
          fontWeight: 100,
          letterSpacing: ".1rem",
          textTransform: "uppercase",
          marginRight: "10em",
        }}
      >
        Tidsapparat
      </Typography>
      <TitlePortal />
    </Box>
  </AppBar>
);

export const MyMenu = (props: MenuProps) => {
  const { version } = useAppConfig();

  return (
    <Box sx={{ display: "flex", flexDirection: "column", height: "100%" }}>
      <Box sx={{ flex: 1 }}>
        <Menu {...props} />
      </Box>

      <Box sx={{ p: 2, textAlign: "center" }}>
        <Typography variant="caption" color="text.secondary">
          Version: {version}
        </Typography>
      </Box>
    </Box>
  );
};

export const MyLayout = (props: LayoutProps) => (
  <Layout {...props} appBar={MyAppBar} menu={MyMenu} />
);
