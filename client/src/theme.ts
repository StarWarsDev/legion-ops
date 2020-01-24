import { createMuiTheme } from "@material-ui/core";

const primary = {
  main: "#424242",
  light: "#6d6d6d",
  dark: "#1b1b1b"
};

const secondary = {
  main: "#9c27b0",
  light: "#d05ce3",
  dark: "#6a0080"
};

const theme = createMuiTheme({
  palette: {
    type: "dark",
    primary,
    secondary
  },
  overrides: {
    MuiAvatar: {
      colorDefault: {
        backgroundColor: secondary.main,
        color: "white"
      }
    }
  }
});

export default theme;
