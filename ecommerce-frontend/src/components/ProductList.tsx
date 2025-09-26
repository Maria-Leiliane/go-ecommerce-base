import React from 'react';
import type { Product } from '../types/Product';

interface ProductListProps {
    products: Product[];
    onEdit: (product: Product) => void;
    onDelete: (id: number) => void;
    onRefresh: () => void;
}

const ProductList: React.FC<ProductListProps> = ({
                                                     products,
                                                     onEdit,
                                                     onDelete,
                                                     onRefresh,
                                                 }) => {
    if (products.length === 0) {
        return (
            <div className="card product-list-empty">
                <p>Nenhum produto encontrado.</p>
                <button onClick={onRefresh}>Atualizar Lista</button>
            </div>
        );
    }

    return (
        <div className="card product-list-section">
            <h2>Lista de Produtos</h2>
            <button onClick={onRefresh} className="refresh-btn">Atualizar</button>
            <table>
                <thead>
                <tr>
                    <th>Nome</th>
                    <th>Preço</th>
                    <th>Qtd.</th>
                    <th>Ações</th>
                </tr>
                </thead>
                <tbody>
                {products.map((product) => (
                    <tr key={product.id}>
                        <td>{product.name}</td>
                        <td>
                            {product.price.toLocaleString('pt-BR', {
                                style: 'currency',
                                currency: 'BRL',
                            })}
                        </td>
                        <td>{product.amount}</td>
                        <td className="actions-cell">
                            <button
                                className="edit-btn"
                                onClick={() => onEdit(product)}
                            >
                                Editar
                            </button>
                            <button
                                className="delete-btn"
                                onClick={() => product.id && onDelete(product.id)}
                                disabled={!product.id}
                            >
                                Excluir
                            </button>
                        </td>
                    </tr>
                ))}
                </tbody>
            </table>
        </div>
    );
};

export default ProductList;