import axios from 'axios';
import type { Product } from '../types/Product'; // Boa prática: importar o tipo

const apiClient = axios.create({
    baseURL: 'http://localhost:8080', // O endereço da sua API Go
    headers: {
        'Content-Type': 'application/json',
    },
});

// Corrigido: Usando os parâmetros para a paginação.
// Renomeei 'p0' para 'limit' para clareza.
export const getProducts = (page: number, limit: number) =>
    apiClient.get(`/products?page=${page}&limit=${limit}`);

// Corrigido: Usando o tipo Product para mais segurança
export const createProduct = (productData: Product) =>
    apiClient.post('/products', productData);

export const updateProduct = (id: number, productData: Product) =>
    apiClient.put(`/products/${id}`, productData);

export const deleteProduct = (id: number) =>
    apiClient.delete(`/products/${id}`);