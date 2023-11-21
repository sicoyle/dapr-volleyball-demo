import React from "react";
import ReactDOM from "react-dom/client";
import App from "./App.tsx";
import { AppBar, Box, Button, Container, CssBaseline, Link, ThemeProvider, Toolbar, Typography } from "@mui/material";
import theme from "./theme/index.ts";
import { SportsVolleyballSharp } from "@mui/icons-material";

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <AppBar position="static">
        <Container maxWidth="xl">
          <Toolbar disableGutters>
            <SportsVolleyballSharp fontSize="large" />
            <Typography
              variant="h6"
              noWrap
              component="a"
              href="#app-bar-with-responsive-menu"
              sx={{
                ml: 1,
                mr: 2,
                display: { xs: "none", md: "flex" },
                fontFamily: "monospace",
                fontWeight: 700,
                letterSpacing: ".3rem",
                color: "inherit",
                textDecoration: "none",
              }}
            >
              VolleyTracker
            </Typography>
            <Box sx={{ flexGrow: 1, display: { xs: "none", md: "flex" } }}>
              <Button
                component={Link}
                // sx={{ my: 2, color: "white", display: "block" }}
                href="https://dapr.io/"
                target="_blank"
                rel="noopener noreferrer"
              >
                Dapr
              </Button>
              <Button
                component={Link}
                // sx={{ my: 2, color: "white", display: "block" }}
                href={import.meta.env.VITE_REACT_APP_ZIPKIN_URL}
                target="_blank"
                rel="noopener noreferrer"
              >
                Zipkin
              </Button>
            </Box>
          </Toolbar>
        </Container>
      </AppBar>
      <App />
    </ThemeProvider>
  </React.StrictMode>
);
