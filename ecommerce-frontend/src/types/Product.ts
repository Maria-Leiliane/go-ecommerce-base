export interface Product {
    id?: number;       // opcional, pois produto novo ainda não tem ID
    name: string;
    price: number;     // sempre número, inicialize como 0
    amount: number;
    description?: string;
}
