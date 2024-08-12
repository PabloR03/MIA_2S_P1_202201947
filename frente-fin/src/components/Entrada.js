import React, { useState } from 'react';
import './Entrada.css'; // Asegúrate de tener este archivo CSS

export const Entrada = () => {
    const [inputText, setInputText] = useState('');

    const handleInputChange = (event) => {
        setInputText(event.target.value);
    };

    return (
        <div className="entrada-container">
            <h2 className="entrada-title">Entrada</h2>
            <textarea 
                className="entrada-textarea"
                value={inputText} 
                onChange={handleInputChange} 
                placeholder="Escribe tu texto aquí..."
                rows="10"
                cols="400"
            />
        </div>
    );
};
