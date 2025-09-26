import React from "react";
import CurrencyInput from "react-currency-input-field";
import type {Product} from "../types/Product";

interface ProductFormProps {
    currentProduct: Product;
    isEditing: boolean;
    onInputChange: (
        e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>
    ) => void;
    onValueChange: (value: string | undefined, name?: string) => void;
    onSubmit: (e: React.FormEvent<HTMLFormElement>) => void;
    onCancel: () => void;
}

const ProductForm: React.FC<ProductFormProps> = ({
                                                     currentProduct,
                                                     isEditing,
                                                     onInputChange,
                                                     onValueChange,
                                                     onSubmit,
                                                     onCancel,
                                                 }) => {
    return (
        <div className="form-section card">
            <h2>{isEditing ? "Editar Produto" : "Adicionar Novo Produto"}</h2>
            <form onSubmit={onSubmit}>
                <div className="form-group">
                    <label htmlFor="name">Nome do Produto</label>
                    <input
                        id="name"
                        name="name"
                        value={currentProduct.name}
                        onChange={onInputChange}
                        placeholder="Ex: Mouse Gamer"
                        required
                    />
                </div>

                <div className="form-group">
                    <label htmlFor="price">Preço (R$)</label>
                    <CurrencyInput
                        id="price"
                        name="price"
                        placeholder="Ex: R$ 15,50"
                        value={currentProduct.price}
                        decimalsLimit={2}
                        decimalSeparator=","
                        groupSeparator="."
                        prefix="R$ "
                        onValueChange={onValueChange}
                        className="currency-input"
                        required
                    />
                </div>

                <div className="form-group">
                    <label htmlFor="amount">Quantidade</label>
                    <input
                        id="amount"
                        name="amount"
                        type="number"
                        min="0"
                        value={currentProduct.amount}
                        onChange={onInputChange}
                        placeholder="Ex: 50"
                        required
                    />
                </div>

                <div className="form-group full-width">
                    <label htmlFor="description">Descrição</label>
                    <textarea
                        id="description"
                        name="description"
                        value={currentProduct.description || ""}
                        onChange={onInputChange}
                        placeholder="Ex: Mouse gamer com iluminação RGB..."
                    />
                </div>

                <div className="form-buttons full-width">
                    <button type="submit">
                        {isEditing ? "Salvar Alterações" : "Criar Produto"}
                    </button>
                    {isEditing && (
                        <button type="button" className="cancel-btn" onClick={onCancel}>
                            Cancelar
                        </button>
                    )}
                </div>
            </form>
        </div>
    );
};

export default ProductForm;
