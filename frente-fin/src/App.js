import React from 'react';
import './App.css';
import { Entrada } from './components/Entrada';
import '@fortawesome/fontawesome-free/css/all.min.css';
import { BarraNavegacion } from './components/BarraNavegacion';
import { Salida } from './components/Salida';

function App() {
  return (
    <div className="App">
      <header className="App-header">
        <h1>MIA Project 202201947</h1>
        <Entrada />
        <BarraNavegacion />
        <Salida/>
      </header>
    </div>
  );
}

export default App;
