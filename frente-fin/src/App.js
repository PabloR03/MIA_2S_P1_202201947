import React, { useEffect, useState } from 'react';
import './App.css';
import { Entrada } from './components/Entrada';
import '@fortawesome/fontawesome-free/css/all.min.css';
import { BarraNavegacion } from './components/BarraNavegacion';
import { Salida } from './components/Salida';

function App() {
  const [message, setMessage] = useState('');

  useEffect(() => {
    // Realizar una solicitud GET al backend cuando el componente se monta
    fetch('http://localhost:8080/api/message')
      .then(response => response.json())
      .then(data => setMessage(data.text))
      .catch(error => console.error('Error:', error));
  }, []);

  const sendMessage = () => {
    // Realizar una solicitud POST al backend
    fetch('http://localhost:8080/api/message', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ text: 'Hello from React!' }),
    })
      .then(response => response.json())
      .then(data => console.log('Response:', data))
      .catch(error => console.error('Error:', error));
  };

  return (
    <div className="App">
      <header className="App-header">
        <h1>MIA Project 202201947</h1>
        <Entrada />
        <BarraNavegacion />
        <Salida />
        <p>{message}</p> {/* Mostrar el mensaje recibido del backend */}
        <button onClick={sendMessage}>Send Message</button>
      </header>
    </div>
  );
}

export default App;
