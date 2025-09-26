import axios from 'axios';
import type { Product } from '../types/Product';

const apiClient = axios.create({
    baseURL: 'http://localhost:8080',
    headers: {
        'Content-Type': 'application/json',
    },
});

// Parameters for pagination.
export const getProducts = (page: number, limit: number) =>
    apiClient.get(`/products?page=${page}&limit=${limit}`);

// Corrigido: Usando o tipo Product para mais segurança
export const createProduct = (productData: Product) =>
    apiClient.post('/products', productData);

export const updateProduct = (id: number, productData: Product) =>
    apiClient.put(`/products/${id}`, productData);

export const deleteProduct = (id: number) =>
    apiClient.delete(`/products/${id}`);