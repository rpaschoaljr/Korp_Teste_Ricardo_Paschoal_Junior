-- Popula clientes (Padronizado: Maiúsculo, Sem Acentos, Apenas Números no CPF/CNPJ/Telefone)
INSERT INTO clientes (nome, telefone, endereco, cpf)
SELECT 'JOAO SILVA', '11987654321', 'RUA DAS FLORES, 123', '12345678900'
WHERE NOT EXISTS (SELECT 1 FROM clientes WHERE cpf = '12345678900');

INSERT INTO clientes (nome, telefone, endereco, cnpj)
SELECT 'EMPRESA XYZ', '1133334444', 'AV PAULISTA, 1000', '12345678000190'
WHERE NOT EXISTS (SELECT 1 FROM clientes WHERE cnpj = '12345678000190');

-- Popula itens (Padronizado: Maiúsculo, Sem Acentos)
-- Obs: O campo codigo agora é gerado automaticamente pelo banco, 
-- mas inserimos manualmente no seed para manter os IDs de teste.
INSERT INTO itens (codigo, descricao, saldo, preco_base)
SELECT 'PROD-0001', 'PARAFUSO SEXTAVADO', 1000, 0.50
WHERE NOT EXISTS (SELECT 1 FROM itens WHERE codigo = 'PROD-0001');

INSERT INTO itens (codigo, descricao, saldo, preco_base)
SELECT 'PROD-0002', 'PORCA DE ACO', 500, 0.20
WHERE NOT EXISTS (SELECT 1 FROM itens WHERE codigo = 'PROD-0002');
