import React, { useState } from 'react';
import './App.css';
import { Entrada } from './components/Entrada';
import '@fortawesome/fontawesome-free/css/all.min.css';
import { BarraNavegacion } from './components/BarraNavegacion';
import { Salida } from './components/Salida';

function App() {
    const [inputText, setInputText] = useState('');  // Almacena el texto de entrada
    const [outputText, setOutputText] = useState('');  // Almacena el texto de salida

    const handleInputChange = (text) => {
        setInputText(text);  // Actualiza el texto de entrada
    };

    const handleExecute = () => {
        // Realizar una solicitud POST al backend con el texto de entrada
        fetch('http://localhost:8080/api/message', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ text: inputText }),
        })
        .then(response => response.json())
        .then(data => setOutputText(data.text))  // Actualiza el texto de salida con la respuesta del backend
        .catch(error => console.error('Error:', error));
    };

    return (
        <div className="App">
            <header className="App-header">
                <h1>MIA Project 202201947</h1>
                <Entrada onInputChange={handleInputChange} />  {/* Pasamos la función para actualizar el texto de entrada */}
                <BarraNavegacion onExecute={handleExecute} />  {/* Pasamos la función para ejecutar la solicitud */}
                <Salida outputText={outputText} />  {/* Pasamos el texto de salida al componente Salida */}
            </header>
        </div>
    );
}

export default App;
