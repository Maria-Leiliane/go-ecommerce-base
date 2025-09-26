
import { useState, useEffect } from 'react';
import type { ChangeEvent, FormEvent } from 'react';
import type { Product } from './types/Product';
import * as api from './services/api';

import Header from './components/Header';
import Loader from './components/Loader';
import ProductList from './components/ProductList';
import ProductForm from './components/ProductForm';
import './App.css';

const INITIAL_PRODUCT_STATE: Product = {
    name: '',
    price: 0,
    amount: 0,
    description: '',
};

function App() {
    const [products, setProducts] = useState<Product[]>([]);
    const [currentProduct, setCurrentProduct] = useState<Product>(INITIAL_PRODUCT_STATE);
    const [isEditing, setIsEditing] = useState<boolean>(false);
    const [isLoading, setIsLoading] = useState<boolean>(false);
    const [error, setError] = useState<string | null>(null);
    const [currentPage, setCurrentPage] = useState<number>(1);
    const [totalPages, setTotalPages] = useState<number>(1);

    const loadProducts = async (page: number) => {
        setIsLoading(true);
        setError(null);
        try {
            const response = await api.getProducts(page, 10); // Carregando 10 por página
            setProducts(response.data.data || []);
            setTotalPages(response.data.total_pages || 1);
            setCurrentPage(response.data.current_page || 1);
        } catch (err) {
            setError('Falha ao carregar produtos. Tente novamente.');
            console.error('Erro ao buscar produtos:', err);
        } finally {
            setIsLoading(false);
        }
    };

    useEffect(() => {
        void loadProducts(currentPage);
    }, [currentPage]);

    const refreshProducts = () => {
        if (currentPage !== 1) {
            setCurrentPage(1);
        } else {
            void loadProducts(1);
        }
    };

    const handleDelete = async (id: number) => {
        if (window.confirm('Tem certeza que deseja deletar este produto?')) {
            try {
                await api.deleteProduct(id);
                refreshProducts();
            } catch (err) {
                setError('Falha ao deletar produto.');
                console.error('Erro ao deletar produto:', err);
            }
        }
    };

    const handleEdit = (product: Product) => {
        setIsEditing(true);
        setCurrentProduct(product);
        window.scrollTo(0, 0);
    };

    const handleCancel = () => {
        setIsEditing(false);
        setCurrentProduct(INITIAL_PRODUCT_STATE);
    };

    const handleInputChange = (e: ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
        const { name, value } = e.target;
        setCurrentProduct({
            ...currentProduct,
            [name]: name === 'amount' ? parseInt(value, 10) || 0 : value,
        });
    };

    const handlePriceChange = (value: string | undefined, name?: string) => {
        const numericValue = value ? parseFloat(value.replace(',', '.')) : 0;
        if (name) {
            setCurrentProduct((prev) => ({ ...prev, [name]: numericValue }));
        }
    };

    const handleSubmit = async (e: FormEvent) => {
        e.preventDefault();
        setIsLoading(true);
        setError(null);

        try {
            if (isEditing && currentProduct.id) {
                await api.updateProduct(currentProduct.id, currentProduct);
            } else {
                await api.createProduct(currentProduct);
            }
            handleCancel();
            refreshProducts();
        } catch (err) {
            setError('Falha ao salvar produto.');
            console.error('Erro ao salvar produto:', err);
        } finally {
            setIsLoading(false);
        }
    };

    const handlePageChange = (newPage: number) => {
        if (newPage >= 1 && newPage <= totalPages && !isLoading) {
            setCurrentPage(newPage);
        }
    };

    return (
        <div className="container">
            <Header />
            <main>
                <ProductForm
                    currentProduct={currentProduct}
                    isEditing={isEditing}
                    onInputChange={handleInputChange}
                    onValueChange={handlePriceChange}
                    onSubmit={handleSubmit}
                    onCancel={handleCancel}
                />
                {error && <p className="error-message">{error}</p>}
                {isLoading ? (
                    <Loader />
                ) : (
                    <ProductList
                        products={products}
                        onEdit={handleEdit}
                        onDelete={handleDelete}
                        onRefresh={refreshProducts}
                    />
                )}
                {totalPages > 1 && !isLoading && (
                    <div className="pagination-controls">
                        <button
                            onClick={() => handlePageChange(currentPage - 1)}
                            disabled={currentPage <= 1}
                        >
                            Anterior
                        </button>
                        <span>
              Página {currentPage} de {totalPages}
            </span>
                        <button
                            onClick={() => handlePageChange(currentPage + 1)}
                            disabled={currentPage >= totalPages}
                        >
                            Próxima
                        </button>
                    </div>
                )}
            </main>
        </div>
    );
}

export default App;