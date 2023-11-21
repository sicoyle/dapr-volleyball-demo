import { Box, Button, Input, Typography } from "@mui/material";
import { useState } from "react";

function App() {
  const [gameID, setGameID] = useState("");
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  const [data, setData] = useState(null);

  const fetchScoreboard = (e: React.FormEvent) => {
    e.preventDefault();
    e.stopPropagation();
    const url = import.meta.env.VITE_REACT_APP_GAME_SERVICE_URL + "/scoreboard/";
    console.log("Calling: " + url + gameID);
    fetch(url + gameID)
      .then((response) => response.json())
      .then((data) => {
        setData(data);
      })
      .catch((error) => console.error(error));
  };
  return (
    <Box
      sx={{
        p: 2,
      }}
    >
      <Typography variant="h5">Welcome to our Dapr Volleyball Demo!</Typography>

      <Box className="game-id-container" component="form" noValidate autoComplete="off" onSubmit={fetchScoreboard} sx={{ my: 2 }}>
        <Input
          type="text"
          placeholder="Enter Game ID"
          value={gameID}
          onChange={(event: React.ChangeEvent<HTMLInputElement>) => {
            setGameID(event.target.value);
          }}
        />
        <Box my={2}>
          <Button type="submit" onClick={fetchScoreboard}>
            Get Game Score
          </Button>
        </Box>
      </Box>
      {data && (
        <div>
          <h1>Final Game Score</h1>
          <h2>
            {JSON.parse(data).team1Name} vs {JSON.parse(data).team2Name}
          </h2>
          <p>
            {JSON.parse(data).team1Name} score: {JSON.parse(data).team1Score}
          </p>
          <p>
            {JSON.parse(data).team2Name} score: {JSON.parse(data).team2Score}
          </p>
        </div>
      )}
    </Box>
  );
}

export default App;
