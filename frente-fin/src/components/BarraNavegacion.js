import React from 'react';
import './BarraNavegacion.css';

export const BarraNavegacion = ({ onExecute }) => {
    return (
        <div className="navbar">
            <button className="nav-button">
                <i className="fas fa-file-upload"></i> Cargar Archivo
            </button>
            <button className="nav-button" onClick={onExecute}>
                <i className="fas fa-play"></i> Ejecutar
            </button>
        </div>
    );
};