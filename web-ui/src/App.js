import React, { useState } from 'react';

const App = () => {
  const [gameID, setGameID] = useState('');
  const [data, setData] = useState(null);

  const fetchScoreboard = () => {
    fetch(`http://localhost:8080/scoreboard/${gameID}`)
      .then(response => response.json())
      .then(data => {
        setData(data);
      })
      .catch(error => console.error(error));
  };

  return (
    <>
      <div>
        <label htmlFor="gameID">Enter Game ID:</label>
        <input type="text" id="gameID" value={gameID} onChange={e => setGameID(e.target.value)} />
        <button onClick={fetchScoreboard}>Get Game Score</button>
      </div>
      {data && (
        <div>
          <h1>Final Game Score</h1>
          <h2>{JSON.parse(data).team1Name} vs {JSON.parse(data).team2Name}</h2>
          <p>{JSON.parse(data).team1Name} score: {JSON.parse(data).team1Score}</p>
          <p>{JSON.parse(data).team2Name} score: {JSON.parse(data).team2Score}</p>
        </div>
      )}
    </>
  );
};

export default App;
