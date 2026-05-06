import { defaultTheme } from "react-admin";
import type { ThemeOptions } from "@mui/material";

export const MyTheme: ThemeOptions = {
  ...defaultTheme,
  palette: {
    ...defaultTheme.palette,
    mode: "light",
    primary: {
      main: "#3B82F6",
    },
    secondary: {
      main: "#BF00FF",
      contrastText: "#FFFFFF",
    },
    background: {
      default: "#F8FAFC",
      paper: "#FFFFFF",
    },
    error: { main: "#FF3366" },
    warning: { main: "#FF9800" },
    info: { main: "#00BCD4" },
    success: { main: "#10B981" },
  },
  components: {
    ...defaultTheme.components,
    MuiButton: {
      styleOverrides: {
        root: {
          borderRadius: "6px",
          fontWeight: 600,
          textTransform: "none",
        },
      },
    },
  },
};

export const MyDarkTheme: ThemeOptions = {
  ...defaultTheme,
  palette: {
    ...defaultTheme.palette,
    mode: "dark",
    primary: {
      main: "#BF00FF",
      light: "#D946EF",
      dark: "#9300C4",
      contrastText: "#FFFFFF",
    },
    secondary: {
      main: "#00E5FF",
      contrastText: "#000000",
    },
    background: {
      default: "#0F172A",
      paper: "#1E293B",
    },
    error: { main: "#FF3366" },
    warning: { main: "#FFB300" },
    info: { main: "#00E5FF" },
    success: { main: "#00E676" },
  },
  components: {
    ...defaultTheme.components,
    MuiButton: {
      styleOverrides: {
        root: {
          borderRadius: "6px",
          fontWeight: 600,
          textTransform: "none",
        },
      },
    },
    MuiTableCell: {
      styleOverrides: {
        head: {
          backgroundColor: "#0F172A",
        },
      },
    },
    MuiChip: {
      styleOverrides: {
        filledPrimary: {
          backgroundColor: "#BF00FF",
          color: "#FFFFFF",
        },
        root: {
          backgroundColor: "rgba(15, 23, 42, 0.6)",
          border: "1px solid rgba(191, 0, 255, 0.5)",
        },
        label: {
          color: "#E066FF",
        },
      },
    },
    MuiTableRow: {
      styleOverrides: {
        root: {
          "&:hover": {
            backgroundColor: "rgba(255, 255, 255, 0.04) !important",
          },
        },
      },
    },
  },
};
