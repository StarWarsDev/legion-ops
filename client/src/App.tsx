import React from "react";
import clsx from "clsx";
import { ThemeProvider, makeStyles } from "@material-ui/core/styles";
import theme from "./theme";
import { CssBaseline, Container, Grid, Paper } from "@material-ui/core";
import Header from "./Header";

const useStyles = makeStyles(theme => ({
  root: {
    display: "flex"
  },
  appBarSpacer: {
    ...theme.mixins.toolbar
  },
  content: {
    flexGrow: 1,
    height: "100vh",
    overflow: "auto"
  },
  container: {
    paddingTop: theme.spacing(4),
    paddingBottom: theme.spacing(4)
  },
  paper: {
    padding: theme.spacing(2),
    display: "flex",
    overflow: "auto",
    flexDirection: "column"
  },
  fixedHeight: {
    height: 240
  }
}));

const App: React.FC = () => {
  const classes = useStyles();
  const fixedHeightPaper = clsx(classes.paper, classes.fixedHeight);
  return (
    <ThemeProvider theme={theme}>
      <div className={classes.root}>
        <CssBaseline />
        <Header />
        <main className={classes.content}>
          <div className={classes.appBarSpacer} />
          <Container maxWidth="lg" className={classes.container}>
            <Grid container spacing={3}>
              {/* Chart */}
              <Grid item xs={12} md={8} lg={9}>
                <Paper className={fixedHeightPaper}>Paper</Paper>
              </Grid>
              {/* Recent Deposits */}
              <Grid item xs={12} md={4} lg={3}>
                <Paper className={fixedHeightPaper}>Paper</Paper>
              </Grid>
              {/* Recent Orders */}
              <Grid item xs={12}>
                <Paper className={classes.paper}>Paper</Paper>
              </Grid>
            </Grid>
          </Container>
        </main>
        <div>Footer</div>
      </div>
    </ThemeProvider>
  );
};

export default App;
