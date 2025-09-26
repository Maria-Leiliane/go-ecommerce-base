import React from "react";
import logoImage from "../assets/icom-removebg.png";

const Header: React.FC = () => (
    <header className="app-header">
        <img src={logoImage} alt="Logo E-Commerce" className="app-logo" />
        <h1>Gerenciador de Produtos</h1>
    </header>
);

export default Header;
