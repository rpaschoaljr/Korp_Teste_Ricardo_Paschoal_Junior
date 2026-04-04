CREATE TABLE IF NOT EXISTS clientes (
    id SERIAL PRIMARY KEY,
    nome VARCHAR(255) NOT NULL,
    telefone VARCHAR(20),
    endereco TEXT,
    cpf VARCHAR(14),
    cnpj VARCHAR(18)
);

CREATE SEQUENCE IF NOT EXISTS itens_codigo_seq START 1;

CREATE TABLE IF NOT EXISTS itens (
    id SERIAL PRIMARY KEY,
    codigo VARCHAR(50) UNIQUE DEFAULT 'PROD-' || LPAD(nextval('itens_codigo_seq')::text, 4, '0'),
    descricao VARCHAR(255) NOT NULL,
    saldo INT NOT NULL DEFAULT 0,
    preco_base DECIMAL(10, 2) NOT NULL
);

CREATE TABLE IF NOT EXISTS faturas (
    id SERIAL PRIMARY KEY,
    cliente_id INT REFERENCES clientes(id),
    status VARCHAR(20) DEFAULT 'ABERTA',
    valor_total DECIMAL(10, 2) NOT NULL DEFAULT 0.00,
    data_criacao TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS itens_fatura (
    id SERIAL PRIMARY KEY,
    fatura_id INT REFERENCES faturas(id) ON DELETE CASCADE,
    item_id INT REFERENCES itens(id),
    quantidade INT NOT NULL,
    preco_unitario DECIMAL(10, 2) NOT NULL,
    subtotal DECIMAL(10, 2) NOT NULL
);
